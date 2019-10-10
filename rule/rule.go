package main

import (
	"io/ioutil"
	"log"

	"github.com/nocmk2/sachima/dur"
	"github.com/tidwall/gjson"
)

// read and parse rule.json

const jsonPath string = "../data/rule.json"

func parse() {
	f, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		log.Fatal(err)
	}

	r := gjson.GetBytes(f, "feature.shop_has_negative")
	log.Println(string(f))
	log.Println(r)

	t := gjson.GetBytes(f, "datasrc.type")
	log.Println(t)
	d := dur.ReadSQL("select shop_id,shop_has_negative,1PD7_pct from hawaiidb.risk_shop_grade limit 5", "hawaii")
	log.Println("rows = ", d.Rows())
	d.ToSQL("dx_xx", "localmysql8")
	// res := rule.cal(d.GetRow(0))
}

func main() {
	parse()
}
