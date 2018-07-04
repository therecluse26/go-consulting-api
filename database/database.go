package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"../config"
)

var Dbconn *sql.DB
var dberr error

type Statement struct {
	Sql string
	Params map[string]string
}

// Establishes database conection
func DbConnection(conf config.Configuration) {

	dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s", conf.SqlHost, conf.SqlUser, conf.SqlPass, conf.SqlPort, conf.SqlDB)

	Dbconn, dberr = sql.Open("mssql", dsn)
	if dberr != nil {
		fmt.Println(dberr.Error())
	}

}

func DbSelect(dbconn *sql.DB, stmt Statement) map[int]map[string]interface{} {

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

	return allResults

}

// Replaces strings following {{pattern}} with corresponding Param value for "pattern"
func (s *Statement) MergedStmt() string {
	var sqlMerged string = s.Sql
	for k, v := range s.Params {
		sqlMerged = strings.Replace(sqlMerged, "{{" + k + "}}", v, -1)
	}
	return sqlMerged
}