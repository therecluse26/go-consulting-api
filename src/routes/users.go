package routes

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"../database"
	"github.com/auth0/go-jwt-middleware"
)

func SetUserRoutes(router *mux.Router, middleware *jwtmiddleware.JWTMiddleware){
	// User Paths
	router.HandleFunc("/users", GetUsers).Methods("GET")
	router.HandleFunc("/users", NewUser).Methods("PUT")
	router.HandleFunc("/users/{id}", GetUser).Methods("GET")
	router.HandleFunc("/users/{id}/info", GetUserInfo).Methods("GET")
	router.HandleFunc("/users/{id}/roles", GetUserRoles).Methods("GET")

	// Role Paths
	router.HandleFunc("/roles", GetRoles).Methods("GET")
	router.HandleFunc("/roles/{id}", GetRole).Methods("GET")
	router.HandleFunc("/roles/{id}/users", GetRoleUsers).Methods("GET")
}

func GetUsers(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT u.id, u.first_name, u.last_name, u.username 
									FROM People.Users u` }

	database.SelectAndReturnJson(sql, w)

}

func GetRoles(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT r.id, r.name, r.description 
									FROM Company.Roles r` }

	database.SelectAndReturnJson(sql, w)

}

func GetRole(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT r.id, r.name, r.description 
										FROM Company.Roles r WHERE r.id = {{id}}`, Params: mux.Vars(r) }

	database.SelectAndReturnJson(sql, w)

}

func GetRoleUsers(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT u.id, u.first_name, u.last_name, u.username, r.id as role_id, r.name as role_name
										FROM People.Users u
										INNER JOIN People.User_Roles ur on ur.user_id = u.id
										INNER JOIN Company.Roles r on ur.role_id = r.id WHERE r.id = {{id}}`, Params: mux.Vars(r) }

	database.SelectAndReturnJson(sql, w)

}

func GetUser(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT u.id, u.first_name, u.last_name, u.username 
										FROM People.Users u WHERE u.id = {{id}}`, Params: mux.Vars(r) }

	database.SelectAndReturnJson(sql, w)

}

func GetUserInfo(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT u.id, u.first_name, u.last_name, u.username, ui.preferred_name,
										ui.gender, ui.date_of_birth, ui.email, ui.phone_primary, ui.phone_secondary,
										ui.address1, ui.address2, ui.city, ui.state, ui.zip, ui.bio
									FROM People.Users u
										INNER JOIN People.User_Info ui on ui.user_id = u.id
									WHERE u.id = {{id}}`, Params: mux.Vars(r) }

	database.SelectAndReturnJson(sql, w)

}

func GetUserRoles(w http.ResponseWriter, r *http.Request){

	sql := database.Statement{ Sql: `SELECT r.id, r.name, r.description
										FROM Company.Roles r
  											INNER JOIN People.User_Roles ur on ur.role_id = r.id
											INNER JOIN People.Users u on ur.user_id = u.id 
										WHERE u.id = {{id}}`, Params: mux.Vars(r) }

	database.SelectAndReturnJson(sql, w)

}

func NewUser(w http.ResponseWriter, r *http.Request){

	r.ParseForm()

	courseData := map[string]string{ "code": r.Form.Get("code"), "name": r.Form.Get("name"), "description": r.Form.Get("description")  }

	sql := database.Statement{ Sql: ``, Params: courseData }

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
		w.Write([]byte(returnVal))
	}

}

