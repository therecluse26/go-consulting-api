package routes

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"../database"
	"../util"
)

func GetStudents(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT r.name, u.id, u.first_name, u.last_name, u.username, si.major, si.start_date
									FROM People.Users u
										INNER JOIN People.User_Roles ur on u.id = ur.user_id
										INNER JOIN Company.Roles r on ur.role_id = r.id
										LEFT JOIN School.Student_Info si on si.user_id = u.id
									WHERE r.name = 'Student'` }

	database.SelectAndWriteJsonResponse(sql, w)

}

func GetStudent(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT u.id, u.first_name, u.last_name, u.username, si.major, si.start_date
									FROM People.Users u
  										INNER JOIN People.User_Roles ur on u.id = ur.user_id
  										INNER JOIN Company.Roles r on ur.role_id = r.id
										LEFT JOIN School.Student_Info si on si.user_id = u.id
									WHERE r.name = 'Student' 
										AND u.id = {{id}}`, Params: mux.Vars(r) }

	database.SelectAndWriteJsonResponse(sql, w)

}

func NewStudent(w http.ResponseWriter, r *http.Request){

	r.ParseForm()

	courseData := map[string]string{ "code": r.Form.Get("code"), "name": r.Form.Get("name"), "description": r.Form.Get("description")  }

	sql := database.Statement{ Sql: `INSERT INTO People.User_Roles (user_id, role_id) VALUES ('{{id}}', (SELECT cast(r.id AS char(36)) FROM Company.Roles r WHERE r.name = 'Student'));
										INSERT INTO School.Student_Info (user_id) VALUES ('{{id}}');`, Params: courseData }

	_, err := database.DbCreate(database.Dbconn, sql)

	res := map[string]string{}


	if err != nil {
		res["status"] = "error"
		res["data"] = err.Error()
		returnVal, _ := json.Marshal(res)
		util.ErrorHandler(err)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(returnVal))

	} else {
		res["status"] = "success"
		returnVal, _ := json.Marshal(res)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(returnVal))	}

}

