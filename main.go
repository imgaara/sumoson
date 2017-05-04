package main

import (
	"flag"
	"fmt"
	"net/url"

	"encoding/json"

	"github.com/kanga333/sumoson/parser"
)

type PropertyInformation struct {
	Name             string
	Price            int
	FloorPlan        string
	LandArea         int
	BuildingArea     int
	Address          string
	Traffic          []string
	ConstructionDate string
}

//"https://suumo.jp/ikkodate/tokyo/sc_shinjuku/nc_87706145/"
func main() {
	fmt.Println("Hello, playground")
	flag.Parse()
	u, err := url.Parse(flag.Arg(0))
	if err != nil {
		fmt.Println(err)
	}
	p, err := parser.Parse(u)
	if err != nil {
		fmt.Println(err)
	}
	b, err := json.Marshal(p)
	fmt.Println(string(b))
}
