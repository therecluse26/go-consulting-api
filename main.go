package main

import (
	_ "github.com/denisenkom/go-mssqldb"
	"./config"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/gorilla/mux"
	"fmt"
	"strconv"
	"log"
	"net/http"
	"./routes"
	"./database"
)

// Initializes variables in global scope
var conf config.Configuration


func init() {

	// Pulls config variables
	conf = config.BuildConfig()

	// Creates database connection
	database.DbConnection(conf)

}

func main() {

	router := mux.NewRouter().StrictSlash(true)

	/**
	* Routing Paths
	*/
		router.HandleFunc("/", routes.GetStats).Methods("GET")

		// Course Paths
		router.HandleFunc("/courses/all", routes.GetAllCourseData).Methods("GET")
		router.HandleFunc("/courses", routes.GetCourses).Methods("GET")
		router.HandleFunc("/courses", routes.NewCourse).Methods("PUT")
		router.HandleFunc("/courses/{course_id}", routes.GetCourse).Methods("GET")
		//router.HandleFunc("/courses/{code}", routes.UpdateCourse).Methods("POST")
		//router.HandleFunc("/courses/{code}", routes.DeleteCourse).Methods("DELETE")
		router.HandleFunc("/courses/{course_id}/registrations", routes.GetCourseRegistrations).Methods("GET")


		// Course Session Paths
		router.HandleFunc("/courses/{course_id}/sessions", routes.GetCourseSessions).Methods("GET")
		router.HandleFunc("/courses/{course_id}/sessions", routes.NewCourseSession).Methods("PUT")
		router.HandleFunc("/courses/{course_id}/sessions/{session_id}", routes.GetCourseSession).Methods("GET")
		//router.HandleFunc("/courses/{code}/sessions/{session_number}", routes.UpdateCourseSession).Methods("POST")
		//router.HandleFunc("/courses/{code}/sessions/{session_number}", routes.DeleteCourseSession).Methods("DELETE")

		// User Paths
		router.HandleFunc("/users", routes.GetUsers).Methods("GET")
		router.HandleFunc("/users/{id}", routes.GetUser).Methods("GET")
		router.HandleFunc("/users/{id}/roles", routes.GetUserRoles).Methods("GET")
		router.HandleFunc("/users", routes.NewUser).Methods("PUT")

		// Role Paths
		router.HandleFunc("/roles", routes.GetRoles).Methods("GET")
		router.HandleFunc("/roles/{id}", routes.GetRole).Methods("GET")
		router.HandleFunc("/roles/{id}/users", routes.GetRoleUsers).Methods("GET")

		// Student Paths
		router.HandleFunc("/students", routes.GetStudents).Methods("GET")
		router.HandleFunc("/students/{id}", routes.GetStudent).Methods("GET")
		router.HandleFunc("/students/{id}", routes.NewStudent).Methods("PUT")

		

		//router.HandleFunc("/roles", routes.NewRole).Methods("PUT")

	fmt.Println("Listening on port " + strconv.Itoa(conf.ApiPort))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(conf.ApiPort), router))

}
