package main

import (
	"database/sql"
	"io/ioutil"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tidwall/gjson"
)

//Data struct
type Data struct {
	dt map[string][]interface{}
}

//ReadSQL from db
func ReadSQL(sql string, con string) Data {
	return Data{}
}

//ToSQL write Data to db table
func (d *Data) ToSQL(table string, con string) {

}

func readDBConfig(con string) *sql.DB {
	dat, err := ioutil.ReadFile("../data/db.json")
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
	var text string
	db, err := sql.Open(dbtype, ""+dbuser+":"+dbpass+"@tcp("+dbip+":"+dbport+")/"+dbdb+"")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("select * from aaa")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		err := rows.Scan(&text)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(text)
	}

	return db
}

func main() {
	//db, err := sql.Open("mysql", "user:password@/dbname")
	dt := map[string][]interface{}{
		"name": []interface{}{"User1", "User2", "User3"},
		"age":  []interface{}{22, 34, 87},
	}
	readDBConfig("localmysql8")

	d2 := Data{dt}

	d1 := ReadSQL("select id,name from table", "mysql1")
	d1.ToSQL("dw_test", "mysql1")
	d2.ToSQL("dw_test2", "mysql1")

}
