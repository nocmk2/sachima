package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/nocmk2/mathinterval"
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

		r.srcsql = "SELECT " + r.pk + "," + strings.Join(r.featureList, ",") + " FROM " + r.table + " limit 300"

	})
}

// GetScore return an int value get from feature.colname.bin
func (r *Rule) GetScore(colname string, cell string) int64 {
	var res int64
	res = r.featureRaw.Get(colname + ".default").Int()
	bin := r.featureRaw.Get(colname + ".bin") //gjson.GetBytes(r.featureRaw, colname+".bin")
	bintype := r.featureRaw.Get(colname + ".bintype").String()

	// if number
	if bintype == "math" {
		bin.ForEach(func(k, v gjson.Result) bool {
			n, err := strconv.ParseFloat(cell, 64)
			if err != nil {
				n = 0
			}
			if mathinterval.Get(k.String()).Hit(n) {
				res = v.Int()
			}
			return true
		})
	}

	if bintype == "text" {
		bin.ForEach(func(k, v gjson.Result) bool {
			if k.String() == cell {
				res = v.Int()
			}
			return true
		})
	}

	return res
}

func (r *Rule) cal(d dur.Data) {
	r.lazyInit()
	var scores []int64
	// r.featureList

	log.Println(d)

	for i := 0; i < d.Rows(); i++ {
		var score int64
		pk := strings.Split(r.pk, ",")
		log.Println(pk[0], ":", d.Row(i).Col(pk[0]))
		log.Println(pk[1], ":", d.Row(i).Col(pk[1]))
		for _, colname := range r.featureList {
			cell := d.Row(i).Col(colname)
			log.Println(colname, ":", cell)
			score += r.GetScore(colname, cell)
		}
		log.Println("------------")
		scores = append(scores, score)
	}

	log.Println(scores)

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
