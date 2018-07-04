package main

import (
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"fmt"
	"./routes"
	"./config"
	"./database"
)

// Initializes variables in global scope
var conf config.Configuration

func main() {

	conf = config.BuildConfig()

	database.DbConnection(conf)

	router := mux.NewRouter()

	// Routing paths
	router.HandleFunc("/courses/all", routes.AllCourseData).Methods("GET")
	router.HandleFunc("/courses", routes.Courses).Methods("GET")
	router.HandleFunc("/courses", routes.NewCourse).Methods("PUT")
	router.HandleFunc("/courses/{code}", routes.Course).Methods("GET")
	router.HandleFunc("/courses/{code}", routes.UpdateCourse).Methods("POST")
	router.HandleFunc("/courses/{code}", routes.DeleteCourse).Methods("DELETE")
	router.HandleFunc("/courses/{code}/sessions", routes.CourseSessions).Methods("GET")
	router.HandleFunc("/courses/{code}/sessions", routes.NewCourseSession).Methods("PUT")
	router.HandleFunc("/courses/{code}/sessions/{session_number}", routes.CourseSession).Methods("GET")
	router.HandleFunc("/courses/{code}/sessions/{session_number}", routes.UpdateCourseSession).Methods("POST")
	router.HandleFunc("/courses/{code}/sessions/{session_number}", routes.DeleteCourseSession).Methods("DELETE")

	fmt.Println("Listening on port " + strconv.Itoa(conf.ApiPort))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(conf.ApiPort), router))
}
