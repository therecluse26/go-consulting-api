package routes

import (
	"net/http"
	"../database"
)

func GetStats(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT (SELECT count(u.id)
										FROM People.Users u ) as users, (SELECT count(c.id) from School.Courses c) as courses` }

	database.SelectAndReturnJson(sql, w)

}

