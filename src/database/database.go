package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"../config/mainconf"
	"net/http"
	"log"
)

var Dbconn *sql.DB
var Dberr error

type Statement struct {
	Sql string
	Params map[string]string
}

/**
* Establishes database connection
*/
func DbConnection(conf mainconf.Configuration) {

	dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", conf.SqlHost, conf.SqlUser, conf.SqlPass, conf.SqlPort, conf.SqlDB)

	Dbconn, Dberr = sql.Open("mssql", dsn)
	if Dberr != nil {
		fmt.Println(Dberr.Error())
	}
}

/**
* Performs SELECT query
*/
func DbSelect(dbconn *sql.DB, stmt Statement) (map[int]map[string]interface{}, error) {

	// Runs Database query
	rows, err := dbconn.Query(stmt.MergedStmt())
	if err != nil {
		fmt.Println(err.Error())
	}

	cols, _ := rows.Columns()

	defer rows.Close()

	allResults := make(map[int]map[string]interface{})
	valuePtrs := make([]interface{}, len(cols))
	var rowResult map[string]interface{}

	counter := 0

	// Iterates over result set and inserts values into allResults
	for rows.Next() {

		rowResult = make(map[string]interface{})

		// Converts column map to interface for passing to rows.Scan() method
		columns := make([]interface{}, len(cols))

		for i, v := range columns {
			valuePtrs[i] = &columns[i]
			columns[i] = v
		}

		// Assigns values to pointers
		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatal(err)
		}

		// Assigns values to rowResult
		for col, val := range cols {
			rowResult[val] = valuePtrs[col]
		}

		// Inserts current rowResult into allResults map
		allResults[counter] = rowResult

		rowResult = nil

		counter++
	}

	rows.Close()

	return allResults, err

}

/**
* Performs INSERT statement
*/
func DbCreate (dbconn *sql.DB, stmt Statement) (sql.Result, error) {

	query, err := dbconn.Prepare(stmt.MergedStmt())

	result, err := query.Exec()

	return result, err

}

/**
* Replaces strings following {{pattern}} with corresponding Param value for "pattern"
*/
func (s *Statement) MergedStmt() string {
	var sqlMerged string = s.Sql
	for k, v := range s.Params {
		sqlMerged = strings.Replace(sqlMerged, "{{" + k + "}}", v, -1)
	}
	return sqlMerged
}

/**
* Selects and formats select result sets as JSON
*/
func SelectAndReturnJson (sql Statement, w http.ResponseWriter) {

	// Gets result set from DbSelect method
	result, err := DbSelect(Dbconn, sql)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Formats result set as JSON
	jsonString, _ := json.Marshal(result)

	// Writes JSON string to http.ResponseWriter
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonString))
}
