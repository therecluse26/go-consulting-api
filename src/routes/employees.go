package routes

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"../database"
	"github.com/auth0/go-jwt-middleware"
)

func SetEmployeeRoutes(router *mux.Router, middleware *jwtmiddleware.JWTMiddleware){

	// Student Paths
	router.HandleFunc("/employees", GetEmployees).Methods("GET")
	router.HandleFunc("/users/{id}/employee_info", GetEmployee).Methods("GET")
	router.HandleFunc("/users/{id}/employee", NewEmployee).Methods("PUT")
}

func GetEmployees(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT u.id, u.first_name, u.last_name, u.username, ei.title, ei.department, ei.manager_id, ei.start_date
									FROM People.Users u
  										INNER JOIN People.User_Roles ur on u.id = ur.user_id
  										INNER JOIN Company.Roles r on ur.role_id = r.id
										LEFT JOIN Company.Employee_Info ei on ei.user_id = u.id
									WHERE r.name = 'Employee'` }

	database.SelectAndReturnJson(sql, w)

}

func GetEmployee(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT u.id, u.first_name, u.last_name, u.username, ei.title, ei.department, ei.manager_id, ei.start_date
									FROM People.Users u
  										INNER JOIN People.User_Roles ur on u.id = ur.user_id
  										INNER JOIN Company.Roles r on ur.role_id = r.id
										LEFT JOIN Company.Employee_Info ei on ei.user_id = u.id
									WHERE r.name = 'Employee' 
										AND u.id = {{id}}`, Params: mux.Vars(r) }

	database.SelectAndReturnJson(sql, w)

}

func NewEmployee(w http.ResponseWriter, r *http.Request){

	r.ParseForm()

	employeeData := map[string]string{ "title": r.Form.Get("title"), "department": r.Form.Get("department"), "manager_id": r.Form.Get("manager_id"), "start_date": r.Form.Get("start_date")  }

	sql := database.Statement{ Sql: `INSERT INTO People.User_Roles (user_id, role_id) VALUES ('{{id}}', (SELECT cast(r.id AS char(36)) FROM Company.Roles r WHERE r.name = 'Employee'));
										INSERT INTO Company.Employee_Info (user_id) VALUES ('{{id}}');`, Params: employeeData }

	_, err := database.DbCreate(database.Dbconn, sql)

	res := map[string]string{}


	if err != nil {
		res["status"] = "error"
		res["data"] = err.Error()
		returnVal, _ := json.Marshal(res)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(returnVal))

	} else {
		res["status"] = "success"
		returnVal, _ := json.Marshal(res)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(returnVal))	}

}

