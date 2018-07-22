package routes

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"../database"
)

func GetUsers(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT u.id, u.first_name, u.last_name, u.username 
									FROM Users u` }

	result, err := database.DbSelect(database.Dbconn, sql)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonString, _ := json.Marshal(result)

	fmt.Fprintf(w, "%s", jsonString)

}

func GetRoles(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT r.id, r.name, r.description 
									FROM Roles r` }

	result, err := database.DbSelect(database.Dbconn, sql)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonString, _ := json.Marshal(result)

	fmt.Fprintf(w, "%s", jsonString)

}

func GetRole(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT r.id, r.name, r.description 
										FROM Roles r WHERE r.id = {{id}}`, Params: mux.Vars(r) }

	result, err := database.DbSelect(database.Dbconn, sql)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonString, _ := json.Marshal(result)

	fmt.Fprintf(w, "%s", jsonString)

}

func GetRoleUsers(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT u.id, u.first_name, u.last_name, u.username, r.id as role_id, r.name as role_name
										FROM Users u
										INNER JOIN User_Roles ur on ur.user_id = u.id
										INNER JOIN Roles r on ur.role_id = r.id WHERE r.id = {{id}}`, Params: mux.Vars(r) }

	result, err := database.DbSelect(database.Dbconn, sql)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonString, _ := json.Marshal(result)

	fmt.Fprintf(w, "%s", jsonString)

}

func GetUser(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT u.id, u.first_name, u.last_name, u.username 
										FROM Users u WHERE u.id = {{id}}`, Params: mux.Vars(r) }

	result, err := database.DbSelect(database.Dbconn, sql)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonString, _ := json.Marshal(result[0])

	fmt.Fprintf(w, "%s", jsonString)

}

func GetUserRoles(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT r.id, r.name, r.description, u.id AS user_id, u.first_name, u.last_name, u.username
										FROM Roles r
  											INNER JOIN User_Roles ur on ur.role_id = r.id
											INNER JOIN Users u on ur.user_id = u.id 
										WHERE u.id = {{id}}`, Params: mux.Vars(r) }

	result, err := database.DbSelect(database.Dbconn, sql)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonString, _ := json.Marshal(result)

	fmt.Fprintf(w, "%s", jsonString)

}

func NewUser(w http.ResponseWriter, r *http.Request){

	r.ParseForm()

	courseData := map[string]string{ "code": r.Form.Get("code"), "name": r.Form.Get("name"), "description": r.Form.Get("description")  }

	sql := database.Statement{ Sql: `INSERT INTO Courses (code, name, description) VALUES ('{{code}}', '{{name}}', '{{description}}')`, Params: courseData }

	_, err := database.DbCreate(database.Dbconn, sql)

	res := map[string]string{}


	if err != nil {
		res["status"] = "error"
		res["data"] = err.Error()
		returnVal, _ := json.Marshal(res)
		fmt.Fprintf(w, "%s", returnVal)

	} else {
		res["status"] = "success"
		returnVal, _ := json.Marshal(res)
		fmt.Fprintf(w, "%s", returnVal)
	}

}
