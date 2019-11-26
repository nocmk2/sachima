package rule

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

// ReadRuleFile read rule.json file
func (r *Rule) ReadRuleFile(path string) {
	r.rulePath = path
}

// FeatureList return rule.featurelist
func (r *Rule) FeatureList() []string {
	r.lazyInit()
	return r.featureList
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

//ParsePre run pre condition before bin
func ParsePre(s string, row dur.Row) (int64, bool, error) {
	a := strings.Split(s, " ")
	if len(a) < 4 {
		return 999, false, nil
	}
	log.Println(a)
	col, err := strconv.ParseFloat(row.Col(a[0]), 64)
	if err != nil {
		return 999, false, err
	}
	res, err := strconv.ParseInt(a[3], 0, 64)
	if err != nil {
		return 999, false, err
	}
	var need bool
	need = false

	target, error := strconv.ParseFloat(a[2], 64)
	if error != nil {
		log.Panic(error)
	}
	log.Println(target)

	switch a[1] {
	default:
		log.Panic("wrong pre format" + s)
	case ">":
		if col > target {
			need = true
		}
		log.Println(">")
	case "<":
		if col < target {
			need = true
		}
		log.Println("<")
	case ">=":
		if col >= target {
			need = true
		}
		log.Println(">=")
	case "<=":
		if col <= target {
			need = true
		}
		log.Println("<=")
	case "=":
		if col == target {
			need = true
		}
		log.Println("=")
	case "==":
		if col == target {
			need = true
		}
		log.Println("==")
	}

	// score needscoreornot error
	return res, need, nil
}

// GetScore return an int value get from feature.colname.bin
func (r *Rule) GetScore(colname string, row dur.Row) int64 {
	cell := row.Col(colname)
	var res int64
	res = r.featureRaw.Get(colname + ".default").Int()
	bin := r.featureRaw.Get(colname + ".bin") //gjson.GetBytes(r.featureRaw, colname+".bin")
	bintype := r.featureRaw.Get(colname + ".bintype").String()
	// pre logic
	pre := r.featureRaw.Get(colname + ".pre").String()
	preScore, need, error := ParsePre(pre, row)
	if error != nil {
		log.Println(error)
	}
	if need {
		return preScore
	}

	// if number
	if bintype == "math" {
		bin.ForEach(func(k, v gjson.Result) bool {
			n, err := strconv.ParseFloat(cell, 64)
			if err != nil {
				return false // if cannot conv to flat64 for example null then break and use default value
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
				// pre "m1_apply_cnt <=1 0"

				// log.Println(colname, ":", cell)
				binscore := r.GetScore(colname, d.Row(i))
				score += binscore
				resData.Col(colname)[i] = binscore
				// resData.Col(colname).Row(i) = score
				// resData.Row(i).Col(colname) = score
			}
			scores = append(scores, score)
			// log.Println(score)
		}

		// log.Println(scores)
		normScores := scores.Normalize().Mul(weight)
		resData.InsertCol(catalogName+"_ORI", scores)
		resData.InsertCol(catalogName, normScores)

		// d.InsertCol(k.String(), normScores)
		log.Println("normScores", normScores)

		return true
	})

	// log.Println(d)

	col := resData.Col(catalogList[0])
	log.Println("catalogList", catalogList)
	for i := 1; i < len(catalogList); i++ {
		// log.Println(catalogList[i])
		col = col.Add(resData.Col(catalogList[i]))
		log.Println("-----------")
		// log.Println(catalogList[i])
	}

	pk := strings.Split(r.pk, ",")
	resData.InsertCol(pk[0], d.Col(pk[0]))
	resData.InsertCol(pk[1], d.Col(pk[1]))

	// log.Println(resData.Col(catalogList[0]))
	// log.Println(resData.Col(catalogList[1]))
	// log.Println("col++++++++++++++++++++++++")
	// log.Println(col)

	resData.InsertCol(r.colName, col)
	// rank := resData.Col(r.colName).Rank()
	// resData.InsertCol("rank", rank)

	pctCol := resData.Col(r.colName).Percentile(true)
	resData.InsertCol("PERCENTILE", pctCol)

	resData.InsertCol("GRADE", r.getRulerGrade(pctCol))
	// log.Println(resData)
	resData.ToSQL(r.DataTargetTable, r.DataTargetType, false)
	// log.Println(d.Row(1).Col("GRADEX"))

}

func (r *Rule) getRulerGrade(pctcol dur.Col) dur.Col {
	// percents := col.Percentile(true)
	var res dur.Col
	for i := 0; i < pctcol.Len(); i++ {
		r.RulersRaw.ForEach(func(k, v gjson.Result) bool {
			if mathinterval.Get(k.String()).Hit((pctcol[i]).(float64)) {
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

//Run for batch
func Run() {
	log.Println("rule run.....")
	parse()
}
