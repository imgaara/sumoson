package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"net/url"

	"regexp"

	"github.com/PuerkitoBio/goquery"
)

type PropertyInformation struct {
	ID               string   `json:"id"`
	PostingDate      string   `json:"posting_date"`
	Name             string   `json:"name"`
	Price            int      `json:"price_yen"`
	FloorPlan        string   `json:"floor_plan"`
	LandArea         float64  `json:"land_area_m2"`
	BuildingArea     float64  `json:"building_area_m2"`
	Address          string   `json:"address"`
	Traffic          []string `json:"traffic"`
	ConstructionDate string   `json:"construction_date"`
}

func Parse(u *url.URL) (*PropertyInformation, error) {
	doc, err := goquery.NewDocument(u.String())
	if err != nil {
		return nil, err
	}

	var p PropertyInformation

	id, err := extractIDFromURL(u)
	if err != nil {
		return nil, err
	}
	p.ID = id

	pd, err := extractPostingDateFromDocument(doc)
	if err != nil {
		return nil, err
	}
	p.PostingDate = pd

	table := doc.Find("#mainContents .secTitleInnerR:contains('物件詳細情報')").Parent().Parent().Find("table")
	if table == nil {
		err = fmt.Errorf("物件詳細情報 Table not Found %v", u.RawQuery)
		return nil, err
	}

	name, err := extractNameFromSelection(table)
	if err != nil {
		return nil, err
	}
	p.Name = name

	price, err := extractPriceFromSelection(table)
	if err != nil {
		return nil, err
	}
	p.Price = price

	floorPlan, err := extractFloorPlanFromSelection(table)
	if err != nil {
		return nil, err
	}
	p.FloorPlan = floorPlan

	landArea, err := extractLandAreaFromSelection(table)
	if err != nil {
		return nil, err
	}
	p.LandArea = landArea

	buildingArea, err := extractBuildingAreaFromSelection(table)
	if err != nil {
		return nil, err
	}
	p.BuildingArea = buildingArea

	constructionDate, err := extractConstructionDateFromSelection(table)
	if err != nil {
		return nil, err
	}
	p.ConstructionDate = constructionDate

	address, err := extractAddressFromSelection(table)
	if err != nil {
		return nil, err
	}
	p.Address = address

	traffic, err := extractTrafficFromSelection(table)
	if err != nil {
		return nil, err
	}
	p.Traffic = traffic

	return &p, nil
}

func extractTrafficFromSelection(s *goquery.Selection) ([]string, error) {
	var traffic []string
	s.Find("th:contains('交通')").Next().Find("div").Each(func(_ int, s *goquery.Selection) {
		traffic = append(traffic, s.Text())
	})
	if len(traffic) == 0 {
		err := fmt.Errorf("交通 not Found %v", s)
		return nil, err
	}
	return traffic, nil
}

func extractAddressFromSelection(s *goquery.Selection) (string, error) {
	extractText := s.Find("th:contains('住所')").Next().Find("p:first-child").Text()
	if extractText == "" {
		err := fmt.Errorf("住所 not Found %v", s)
		return "", err
	}
	return strings.TrimSpace(extractText), nil
}

func extractConstructionDateFromSelection(s *goquery.Selection) (string, error) {
	extractText := s.Find(".fl:contains('築年月')").ParentsFiltered("th").Next().Text()
	if extractText == "" {
		err := fmt.Errorf("築年月 not Found %v", s)
		return "", err
	}
	constructionDateString := strings.Trim(strings.TrimSpace(extractText), "予定")
	t, err := time.Parse("2006年1月 MST", constructionDateString+" JST")
	if err != nil {
		return "", err
	}
	return t.String(), nil
}

func extractBuildingAreaFromSelection(s *goquery.Selection) (float64, error) {
	extractText := s.Find(".fl:contains('建物面積')").ParentsFiltered("th").Next().Text()
	if extractText == "" {
		err := fmt.Errorf("建物面積 not Found %v", s)
		return 0, err
	}

	r, err := regexp.Compile(`(.*)m2.*$`)
	if err != nil {
		return 0, err
	}

	buildingAreaString := r.FindStringSubmatch(strings.TrimSpace(extractText))[1]

	buildingAreaNumber, err := strconv.ParseFloat(buildingAreaString, 64)
	if err != nil {
		return 0, err
	}

	return buildingAreaNumber, nil
}

func extractLandAreaFromSelection(s *goquery.Selection) (float64, error) {
	extractText := s.Find(".fl:contains('土地面積')").ParentsFiltered("th").Next().Text()
	if extractText == "" {
		err := fmt.Errorf("土地面積 not Found %v", s)
		return 0, err
	}

	r, err := regexp.Compile(`(.*)m2.*$`)
	if err != nil {
		return 0, err
	}

	landAreaString := r.FindStringSubmatch(strings.TrimSpace(extractText))[1]

	landAreaNumber, err := strconv.ParseFloat(landAreaString, 64)
	if err != nil {
		return 0, err
	}

	return landAreaNumber, nil
}

func extractFloorPlanFromSelection(s *goquery.Selection) (string, error) {
	extractText := s.Find(".fl:contains('間取り')").ParentsFiltered("th").Next().Text()
	if extractText == "" {
		err := fmt.Errorf("間取り not Found %v", s)
		return "", err
	}
	return strings.TrimSpace(extractText), nil
}

func extractPriceFromSelection(s *goquery.Selection) (int, error) {
	extractText := s.Find(".fl:contains('価格')").ParentsFiltered("th").Next().Find("p:contains('円')").Text()
	if extractText == "" {
		err := fmt.Errorf("価格 not Found %v", s)
		return 0, err
	}
	r, err := regexp.Compile(`万円$`)
	if err != nil {
		return 0, err
	}
	priceString := strings.Trim(strings.TrimSpace(extractText), "億")
	if !r.MatchString(priceString) {
		err := fmt.Errorf("priceString is not suffix 万円 %v", priceString)
		return 0, err
	}

	priceNumber, err := strconv.Atoi(strings.Trim(priceString, "万円"))
	if err != nil {
		return 0, err
	}
	return priceNumber * 10000, nil
}

func extractNameFromSelection(s *goquery.Selection) (string, error) {
	extractText := s.Find(".fl:contains('物件名')").ParentsFiltered("tr").Find("td").Text()
	if extractText == "" {
		err := fmt.Errorf("物件名 not Found %v", s)
		return "", err
	}
	return strings.TrimSpace(extractText), nil
}

func extractPostingDateFromDocument(d *goquery.Document) (string, error) {
	extractText := d.Find("p:contains('情報提供日')").Text()
	r, err := regexp.Compile(`情報提供日：(.*)$`)
	if err != nil {
		return "", err
	}
	extractDate := r.FindStringSubmatch(extractText)[1]
	t, err := time.Parse("2006/1/2 MST", extractDate+" JST")
	if err != nil {
		return "", err
	}
	return t.String(), nil
}

func extractIDFromURL(u *url.URL) (string, error) {
	r, err := regexp.Compile(`/nc_([0-9]*)/`)
	if err != nil {
		return "", err
	}
	id := r.FindStringSubmatch(u.String())[1]
	if err != nil {
		err = fmt.Errorf("RawPath is not much %v", u.RawPath)
		return "", err
	}
	return id, nil
}
