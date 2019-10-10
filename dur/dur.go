package dur

import (
	"database/sql"
	"io/ioutil"
	"log"

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

//ReadSQL from db
func ReadSQL(sql string, con string) Data {
	return Data{}
}

// Rows return data rows
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
		log.Println("Table " + table + " dropped!")

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
