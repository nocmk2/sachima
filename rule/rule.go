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
	CatalogRaw      gjson.Result
	rulePath        string
	featureRaw      gjson.Result
	featureList     []string
	srcsql          string
	colName         string
	table           string
	pk              string
	DataTargetTable string
	DataTargetType  string
	RulersRaw       gjson.Result
	Where           string
	doOnce          sync.Once
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

func (r *Rule) getFeatures() []string {
	var res []string
	r.featureRaw.ForEach(func(k, v gjson.Result) bool {
		res = append(res, k.String())
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
		r.CatalogRaw = gjson.GetBytes(f, "catalog")
		r.DataTargetTable = gjson.GetBytes(f, "datatarget.table").String()
		r.DataTargetType = gjson.GetBytes(f, "datatarget.type").String()
		r.RulersRaw = gjson.GetBytes(f, "rulers")
		r.Where = gjson.GetBytes(f, "datasrc.where").String()

		r.featureRaw.ForEach(func(k, v gjson.Result) bool {
			r.featureList = append(r.featureList, k.String())
			log.Println(k)
			return true
		})

		r.srcsql = "SELECT " + r.pk + "," + strings.Join(r.featureList, ",") + " FROM " + r.table + " " + r.Where + ""

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

// TODO: 支持前置条件
func (r *Rule) cal(d dur.Data) {
	var resData dur.Data
	r.lazyInit()
	// r.featureList
	var catalogList []string

	for _, colname := range r.getFeatures() {
		resData.InsertCol(colname, make(dur.Col, d.Rows()))
	}

	r.CatalogRaw.ForEach(func(k, v gjson.Result) bool {
		catalogName := k.String()
		var scores dur.Col
		log.Println(k, v)
		catalogList = append(catalogList, catalogName)
		weight := v.Get("weight").Float()
		initScore := v.Get("init_score").Int()

		for i := 0; i < d.Rows(); i++ {
			score := initScore
			for _, colname := range r.getFeaturesByCatalog(catalogName) {
				cell := d.Row(i).Col(colname)
				// log.Println(colname, ":", cell)
				binscore := r.GetScore(colname, cell)
				score += binscore
				resData.Col(colname)[i] = binscore
				// resData.Col(colname).Row(i) = score
				// resData.Row(i).Col(colname) = score
			}
			scores = append(scores, score)
			log.Println(score)
		}

		// log.Println(scores)
		normScores := scores.Normalize().Mul(weight)
		resData.InsertCol(catalogName+"_ORI", scores)
		resData.InsertCol(catalogName, normScores)

		// d.InsertCol(k.String(), normScores)
		log.Println("normScores", normScores)

		return true
	})

	log.Println(d)

	col := resData.Col(catalogList[0])
	log.Println("catalogList", catalogList)
	for i := 1; i < len(catalogList); i++ {
		// log.Println(catalogList[i])
		col = col.Add(resData.Col(catalogList[i]))
		// log.Println("-----------")
		// log.Println(catalogList[i])
		log.Println(col)
	}
	resData.InsertCol(r.colName, col)
	pk := strings.Split(r.pk, ",")
	// log.Println(pk[0], ":", d.Row(i).Col(pk[0]))
	// log.Println(pk[1], ":", d.Row(i).Col(pk[1]))
	resData.InsertCol(pk[0], d.Col(pk[0]))
	resData.InsertCol(pk[1], d.Col(pk[1]))
	// log.Println(d.Rows())
	// resData.ToSQL("dx_score_res", "localmysql8")
	resData.InsertCol("GRADE", r.getRulerGrade(resData.Col(r.colName).Percentile()))
	resData.Col("xxxxll")
	log.Println(resData)
	resData.ToSQL(r.DataTargetTable, r.DataTargetType)
	// log.Println(d.Row(1).Col("GRADEX"))

}

func (r *Rule) getRulerGrade(col dur.Col) dur.Col {
	percents := col.Percentile()
	var res dur.Col
	for i := 0; i < col.Len(); i++ {
		r.RulersRaw.ForEach(func(k, v gjson.Result) bool {
			if mathinterval.Get(k.String()).Hit((percents[i]).(float64)) {
				res = append(res, v.String())
				return false
			}
			return true
		})
	}

	return res
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
