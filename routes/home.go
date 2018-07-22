package routes

import (
	"net/http"
	"encoding/json"
	"fmt"
	"../database"
)

func GetStats(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT (SELECT count(u.id)
										FROM Users u ) as users, (SELECT count(c.id) from Courses c) as courses` }

	result, err := database.DbSelect(database.Dbconn, sql)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonString, _ := json.Marshal(result[0])

	fmt.Fprintf(w, "%s", jsonString)

}

