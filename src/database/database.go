package database

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"encoding/json"
	"fmt"
	"strings"
	"net/http"
	"github.com/therecluse26/fortisure-api/src/config/mainconf"
	"github.com/therecluse26/fortisure-api/src/util"
)

var Dbconn *sql.DB
var Dberr error
var Conf = mainconf.BuildConfig()

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
		util.ErrorHandler(Dberr)
	}
}

/**
* Performs SELECT query
*/
func DbSelect(dbconn *sql.DB, stmt Statement) (map[int]map[string]interface{}, error) {

	// Runs Database query
	rows, err := dbconn.Query(stmt.MergedStmt())
	if err != nil {
		util.ErrorHandler(err)
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
			util.ErrorHandler(err)
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
	var sqlMerged = s.Sql
	for k, v := range s.Params {
		sqlMerged = strings.Replace(sqlMerged, "{{" + k + "}}", v, -1)
	}
	return sqlMerged
}

/**
* Selects and formats select result sets as JSON
*/
func SelectAndWriteJsonResponse(sql Statement, w http.ResponseWriter) {

	// Gets result set from DbSelect method
	result, err := DbSelect(Dbconn, sql)
	if err != nil {
		util.ErrorHandler(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Formats result set as JSON
	jsonString, _ := json.Marshal(result)

	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")
	// Writes JSON string to http.ResponseWriter
	_, err = w.Write([]byte(jsonString)); if err != nil{ fmt.Println(err) }
}

/**
* Selects and returns simple result set
*/
func SelectAndReturnResultSet(sql Statement) (map[int]map[string]interface{}, error) {

	// Gets result set from DbSelect method
	result, err := DbSelect(Dbconn, sql)

	return result, err
}

/**
* Selects and returns simple result set
*/
func ExecuteRawQuery(stmt string) error {

	sql := Statement{ Sql: stmt }

	// Executes Database query
	_, err := Dbconn.Query(sql.Sql)
	if err != nil {
		util.ErrorHandler(err)
	}

	return err
}


/**
 * Selects a simple count from the database
 */
func SelectAndReturnCount (sql Statement) (int, error) {

	result, err := SelectSingleCountValue(Dbconn, sql)

	return result, err
}



func SelectSingleCountValue (dbconn *sql.DB, stmt Statement) (int, error) {

	var count int

	// Runs Database query
	rows, err := dbconn.Query(stmt.MergedStmt())
	if err != nil {
		util.ErrorHandler(err)
	}

	// Selects first row
	rows.Next()

	// Inserts result into `count` variable
	err = rows.Scan(&count)
	if err != nil {
		util.ErrorHandler(err)
	}

	return count, err

}