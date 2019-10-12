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
	catalog     gjson.Result
	rulePath    string
	featureRaw  gjson.Result
	featureList []string
	srcsql      string
	colName     string
	table       string
	pk          string
	doOnce      sync.Once
}

func (r *Rule) getFeaturesByCatalog(catalogName string) []string {
	var res []string
	r.featureRaw.ForEach(func(k, v gjson.Result) bool {
		if v.Get("catalog").String() == catalogName {
			res = append(res, k.String())
		}
		return true
	})

	return res
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
		r.catalog = gjson.GetBytes(f, "catalog")

		r.featureRaw.ForEach(func(k, v gjson.Result) bool {
			r.featureList = append(r.featureList, k.String())
			log.Println(k)
			return true
		})

		r.srcsql = "SELECT " + r.pk + "," + strings.Join(r.featureList, ",") + " FROM " + r.table + " limit 15"

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

// Normalize (x-Min)/(MAX-MIN)
func Normalize(a []int64, percent float64) dur.Col {
	len := len(a)
	min := a[0]
	max := a[1]
	res := make(dur.Col, len)
	for i := 0; i < len; i++ {
		if a[i] < min {
			min = a[i]
		}

		if a[i] > max {
			max = a[i]
		}
	}

	if min == max {
		return res
	}

	for i := 0; i < len; i++ {
		res[i] = percent * ((float64)(a[i]-min) / (float64)(max-min))
	}

	return res
}

// TODO: 开发ruler规则
// TODO: 支持前置条件
func (r *Rule) cal(d dur.Data) {
	r.lazyInit()
	// r.featureList
	var catalogList []string

	r.catalog.ForEach(func(k, v gjson.Result) bool {
		catalogName := k.String()
		var scores []int64
		log.Println(k, v)
		catalogList = append(catalogList, catalogName)
		weight := v.Get("weight").Float()
		initScore := v.Get("init_score").Int()

		for i := 0; i < d.Rows(); i++ {
			score := initScore
			// pk := strings.Split(r.pk, ",")
			// log.Println(pk[0], ":", d.Row(i).Col(pk[0]))
			// log.Println(pk[1], ":", d.Row(i).Col(pk[1]))
			for _, colname := range r.getFeaturesByCatalog(catalogName) {
				cell := d.Row(i).Col(colname)
				// log.Println(colname, ":", cell)
				score += r.GetScore(colname, cell)
			}
			// log.Println("------------")
			scores = append(scores, score)
			log.Println(score)
		}

		// log.Println(scores)
		normScores := Normalize(scores, weight)

		d.InsertCol(catalogName, dur.Col(normScores))

		// d.InsertCol(k.String(), normScores)
		log.Println("normScores", normScores)

		return true
	})

	log.Println(d)

	col := d.Col(catalogList[0])
	log.Println("catalogList", catalogList)
	for i := 1; i < len(catalogList); i++ {
		// log.Println(catalogList[i])
		col = col.Add(d.Col(catalogList[i]))
		// log.Println("-----------")
		// log.Println(catalogList[i])
	}
	d.InsertCol(r.colName, col)
	// log.Println(d.Rows())
	log.Println(d)
	// log.Println(d.Row(1).Col("GRADEX"))

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
