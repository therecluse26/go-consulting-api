package routes

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"../database"
)

func GetStudents(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT u.id, u.first_name, u.last_name, u.username, ui.preferred_name, ui.gender, ui.date_of_birth,
        								ui.email, ui.phone_primary, ui.phone_secondary, ui.address1, ui.address2, ui.city,
        								ui.state, ui.zip, ui.bio
									FROM Users u
										INNER JOIN User_Info ui on u.id = ui.user_id
										INNER JOIN Student_Info si on si.user_id = u.id
  										INNER JOIN User_Roles ur on u.id = ur.user_id
  										INNER JOIN Roles r on ur.role_id = r.id
										WHERE r.name = 'Student'` }

	result, err := database.DbSelect(database.Dbconn, sql)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonString, _ := json.Marshal(result)

	fmt.Fprintf(w, "%s", jsonString)

}

func GetStudent(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT u.id, u.first_name, u.last_name, u.username, ui.preferred_name, ui.gender, ui.date_of_birth,
        								ui.email, ui.phone_primary, ui.phone_secondary, ui.address1, ui.address2, ui.city,
        								ui.state, ui.zip, ui.bio
									FROM Users u
										INNER JOIN User_Info ui on u.id = ui.user_id
										INNER JOIN Student_Info si on si.user_id = u.id
  										INNER JOIN User_Roles ur on u.id = ur.user_id
  										INNER JOIN Roles r on ur.role_id = r.id
										WHERE r.name = 'Student' 
										AND u.id = {{id}}`, Params: mux.Vars(r) }

	result, err := database.DbSelect(database.Dbconn, sql)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	jsonString, _ := json.Marshal(result)

	fmt.Fprintf(w, "%s", jsonString)

}

func NewStudent(w http.ResponseWriter, r *http.Request){

	r.ParseForm()

	courseData := map[string]string{ "code": r.Form.Get("code"), "name": r.Form.Get("name"), "description": r.Form.Get("description")  }

	sql := database.Statement{ Sql: `INSERT INTO User_Roles (user_id, role_id) VALUES ('{{id}}', (SELECT cast(r.id AS char(36)) FROM Roles r WHERE r.name = 'Student'));
										INSERT INTO Student_Info (user_id) VALUES ('{{id}}');`, Params: courseData }

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

