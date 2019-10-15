package dur

import (
	"database/sql"
	"io/ioutil"
	"log"
	"sort"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/tidwall/gjson"
)

const tablePrefix string = "dx"
const dbConPath string = "../data/db.json"

//Data struct
type Data struct {
	dt map[string][]interface{}
}

//Row reps one row of Data
type Row struct {
	dt map[string]interface{}
}

// Col reps one col of Data
type Col []interface{}

// NewCol ...
// func NewCol(d interface{}) Col {
// 	// return Col{d.([]interface{})}
// 	switch d.(type) {
// 	case []string:
// 		log.Println("sssssssssssssss")
// 		return Col{d.([]interface{})}
// 	case []float64:
// 		log.Println("ffffffffffffff")
// 	case []int64:
// 		log.Println("iiiiiiiiiiiiiiiiiiiiii")
// 	}
// }

// Mul  Multiple c by n
func (c Col) Mul(n float64) Col {
	for i := 0; i < len(c); i++ {
		c[i] = c[i].(float64) * n
	}
	return c
}

// Div Divid c by n
func (c Col) Div(n float64) Col {
	switch c[0].(type) {
	case float64:
		for i := 0; i < len(c); i++ {
			c[i] = c[i].(float64) / n
		}
	case int:
		for i := 0; i < len(c); i++ {
			c[i] = (float64)(c[i].(int)) / n
		}
	}
	return c
}

// Normalize ..
func (c Col) Normalize() Col {
	len := len(c)
	min := c[0].(int64)
	max := c[1].(int64)
	res := make(Col, len)
	for i := 0; i < len; i++ {
		if c[i].(int64) < min {
			min = c[i].(int64)
		}

		if c[i].(int64) > max {
			max = c[i].(int64)
		}
	}

	if min == max {
		return res
	}

	for i := 0; i < len; i++ {
		res[i] = ((float64)(c[i].(int64) - min)) / (float64)(max-min)
	}

	return res
}

// Add :  Add tow Col
func (c Col) Add(added Col) Col {
	newCol := make(Col, len(c))
	for i := 0; i < len(c); i++ {
		// log.Println(i, "------------------+++")
		// log.Println(c[i], added[i])
		newCol[i] = c[i].(float64) + added[i].(float64)
		// log.Println(newCol)
	}
	// log.Println("c", c)
	// log.Println("added", added)
	// log.Println("newCol", newCol)
	return newCol
}

func (c Col) Len() int { return len(c) }
func (c Col) Less(i, j int) bool {
	// log.Println(reflect.TypeOf(c[0]))
	switch c[0].(type) {
	case string:
		log.Panic("string Less cmp")
	case int:
		return c[i].(int) < c[j].(int)
	case int64:
		return c[i].(int64) < c[j].(int64)
	case float64:
		return c[i].(float64) < c[j].(float64)
	}
	return true
}
func (c Col) Swap(i, j int) { c[i], c[j] = c[j], c[i] }

// Sort Col
func (c Col) Sort() Col {
	resCol := make(Col, c.Len())
	copy(resCol, c)
	sort.Sort(resCol)
	return resCol
}

// Rank Col
func (c Col) Rank() Col {
	sortedCol := c.Sort()
	resCol := make(Col, c.Len())

	for i := 0; i < c.Len(); i++ {
		for rank, v := range sortedCol {
			if c[i] == v {
				resCol[i] = rank
				continue
			}
		}
	}
	return resCol
}

// Percentile For example, the 20th percentile is the value (or score) below which 20% of the observations may be found
func (c Col) Percentile() Col {
	log.Println("RANK:", c.Rank())
	return c.Rank().Div((float64)(c.Len()))
}

// InsertCol insert a Col []interface{} to data 在数据中插入一列
func (d *Data) InsertCol(name string, values Col) {
	if d.dt == nil {
		d.dt = map[string][]interface{}{}
	}
	d.dt[name] = values
}

//Col return one col of data type is Data.Col type
func (d *Data) Col(name string) Col {
	return d.dt[name]
}

// Row return the ith row from data
func (d *Data) Row(i int) Row {
	oneRow := Row{make(map[string]interface{}, d.Rows())}
	for colname := range d.dt {
		oneRow.dt[colname] = d.dt[colname][i]
	}
	return oneRow
}

// Col return the name cell of the Row
func (r Row) Col(name string) string {
	return r.dt[name].(string)
	// switch r.dt[name].(type) {
	// default:
	// 	return r.dt[name].(string)
	// case string:
	// 	return r.dt[name].(string)
	// case int:
	// 	return r.dt[name].(string)
	// }
}

//ReadSQL from db
func ReadSQL(sqlstr string, con string) Data {
	// https://github.com/go-sql-driver/mysql/wiki/Examples
	db := readDBConfig(con)
	defer db.Close()
	rows, err := db.Query(sqlstr)
	if err != nil {
		log.Fatal(err)
	}
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	dt := map[string][]interface{}{}

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Fatal(err)
		}
		var value string

		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = ""
			} else {
				value = string(col)
			}
			// log.Println(columns[i], ": ", value)
			dt[columns[i]] = append(dt[columns[i]], value)
		}
	}
	// 	dt := map[string][]interface{}{
	// 		"name":  []interface{}{"User1", "User2", "User3"},
	// 		"age":   []interface{}{222222, 34, 87},
	// 		"coomi": []interface{}{"er3r3r", "xxxxfefef", "feeefefe"},
	// 	}

	return Data{dt}
}

// Rows return number of rows
func (d *Data) Rows() int {
	for _, v := range d.dt {
		return len(v)
	}
	return 0
}

func (d *Data) getCreateStmt(table string) string {
	stmt := "CREATE TABLE " + table + "("
	for col := range d.dt {
		switch d.dt[col][0].(type) {
		default:
			log.Println("++++++++++++++++++++=")
		case float64:
			stmt += col + " float(10,6)"
		case int64:
			stmt += col + " int(11)"
		case string:
			stmt += col + " varchar(50)"
		case int:
			stmt += col + " int"
		}
		stmt += ","
	}
	stmt = stmt[:len(stmt)-1] + ")"
	log.Println(stmt)

	return stmt
}

func (d *Data) getInsertStmt(table string) (string, []string) {
	// concat insert stmt
	var values string
	var cols string
	var columns []string
	for colname := range d.dt {
		cols += colname + ","
		values += "?,"
		columns = append(columns, colname)
	}

	values = values[:len(values)-1]
	cols = cols[:len(cols)-1]
	return "INSERT INTO " + table + "(" + cols + ")" + " VALUES(" + values + ")", columns

}

//ToSQL write Data to db table
func (d *Data) ToSQL(table string, con string) {
	db := readDBConfig(con)
	defer db.Close()

	// check table name
	if table[0:2] != tablePrefix {
		panic("invalid table name! table name should begin with " + tablePrefix)
	}

	// if table not exists create it
	_, err := db.Exec("DROP TABLE IF EXISTS " + table)
	if err == nil {
		log.Println(con + " Table " + table + " dropped!")

		stmt, err := db.Prepare(d.getCreateStmt(table))
		defer stmt.Close()
		if err != nil {
			log.Fatal(err)
		}

		stmt.Exec()
	} else {
		log.Fatal(err)
	}

	insertStmt, columns := d.getInsertStmt(table)

	stmt, err := db.Prepare(insertStmt)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < d.Rows(); i++ {
		var box []interface{}
		for _, colname := range columns {
			box = append(box, d.dt[colname][i])
		}
		log.Println(box, " inserted")
		_, err := stmt.Exec(box...)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func readDBConfig(con string) *sql.DB {
	dat, err := ioutil.ReadFile(dbConPath)
	json := string(dat)
	if err != nil {
		log.Fatal(err)
	}
	dbtype := gjson.Get(json, con+".type").String()
	dbport := gjson.Get(json, con+".port").String()
	dbip := gjson.Get(json, con+".ip").String()
	dbuser := gjson.Get(json, con+".user").String()
	dbpass := gjson.Get(json, con+".pass").String()
	dbdb := gjson.Get(json, con+".db").String()
	// var text string
	db, err := sql.Open(dbtype, ""+dbuser+":"+dbpass+"@tcp("+dbip+":"+dbport+")/"+dbdb+"")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// func main() {
// 	dt := map[string][]interface{}{
// 		"name":  []interface{}{"User1", "User2", "User3"},
// 		"age":   []interface{}{222222, 34, 87},
// 		"coomi": []interface{}{"er3r3r", "xxxxfefef", "feeefefe"},
// 	}

// 	d2 := Data{dt}

// 	d2.ToSQL("dx_test4", "localmysql8")

// }
