package main

import (
	"io/ioutil"
	"log"
	"strings"
	"sync"

	"github.com/nocmk2/sachima/dur"
	"github.com/tidwall/gjson"
)

// read and parse rule.json

const jsonPath string = "../data/rule.json"

//Rule read from rule.json by default
type Rule struct {
	rulePath    string
	featureRaw  gjson.Result
	featureList []string
	srcsql      string
	colName     string
	table       string
	pk          string
	doOnce      sync.Once
}

func (r *Rule) lazyInit() {
	r.doOnce.Do(func() {

		f, err := ioutil.ReadFile(r.rulePath)
		if err != nil {
			log.Fatal(err)
		}

		r.featureRaw = gjson.GetBytes(f, "feature")
		r.colName = gjson.GetBytes(f, "colname").String()
		r.table = gjson.GetBytes(f, "datasrc.name").String()
		r.pk = gjson.GetBytes(f, "datasrc.pk").String()

		r.featureRaw.ForEach(func(k, v gjson.Result) bool {
			r.featureList = append(r.featureList, k.String())
			log.Println(k)
			return true
		})

		r.srcsql = "SELECT " + r.pk + "," + strings.Join(r.featureList, ",") + " FROM " + r.table + " limit 3"

	})
}

func (r *Rule) cal(d dur.Data) {
	r.lazyInit()
	// scores := make([]int, d.Rows())
	// r.featureList

	log.Println(d)

	for i := 0; i < d.Rows(); i++ {
		for _, colname := range r.featureList {
			cell := d.Row(i).Col(colname)
			log.Println(cell)
		}
	}

	// d.ForEach(func(row dur.Row()) bool {
	// 	res = 1334888
	// 	scores = append(scores,res)
	// 	return true
	// })
	//d.Add(scores, "dxscore")
	// d.ToSQL("dx_shop_scores", "localmysql8")
}

func parse() {
	// d.ToSQL("dx_xx", "localmysql8")
	rule1 := Rule{rulePath: jsonPath}
	rule1.lazyInit()
	log.Println(rule1.srcsql)
	d := dur.ReadSQL(rule1.srcsql, "hawaii")
	log.Println(d.Rows())
	// d.Add(scores, "dxscore")
	rule1.cal(d)
}

func main() {
	parse()
}
