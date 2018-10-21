package routes

import (
	"github.com/gorilla/mux"
	"github.com/auth0/go-jwt-middleware"
	"fmt"
)

type RouteMapping struct {
	path string
	object string
}

type RouteObjects map[int]*RouteMapping

var ObjMap RouteObjects

/**
 * Adds mapping to ObjMap object between url path and database object (for permissions checks)
 */
func (om *RouteObjects) addRouteMapping (p string, o string){
	k := len(*om)
	if k == 0 {
		*om = RouteObjects{0: &RouteMapping{path: p, object: o}}
	}
	(*om)[k] = &RouteMapping{path: p, object: o}
}

func GetObjectFromPath (path string) string {
	var obj = ""

	// Add logic for inserting path parameters into template path

	for _, v := range ObjMap {
		if v.path == path {
			obj = v.object
		}
	}
	return obj
}

// GENERAL
func SetUserRoutes(router *mux.Router, middleware *jwtmiddleware.JWTMiddleware){
	// User Paths
	router.HandleFunc("/users", GetUsers).Methods("GET")
	router.HandleFunc("/users", NewUser).Methods("PUT")
		ObjMap.addRouteMapping("/users", "People.Users")

	router.HandleFunc("/users/{id}", GetUser).Methods("GET")
		ObjMap.addRouteMapping("/users/{id}", "People.Users")

	router.HandleFunc("/users/{id}/info", GetUserInfo).Methods("GET")
		ObjMap.addRouteMapping("/users/{id}/info", "People.Users")
		ObjMap.addRouteMapping("/users/{id}/info", "People.User_Info")

	router.HandleFunc("/users/{id}/roles", GetUserRoles).Methods("GET")
		ObjMap.addRouteMapping("/users/{id}/roles", "People.Users")
		ObjMap.addRouteMapping("/users/{id}/roles", "Security.Roles")

	// Role Paths
	router.HandleFunc("/roles", GetRoles).Methods("GET")
		ObjMap.addRouteMapping("/roles", "Security.Roles")

	router.HandleFunc("/roles/{id}", GetRole).Methods("GET")
		ObjMap.addRouteMapping("/roles/{id}", "Security.Roles")

	router.HandleFunc("/roles/{id}/users", GetRoleUsers).Methods("GET")
		ObjMap.addRouteMapping("/roles/{id}/users", "Security.Roles")

	for _, item := range ObjMap {
		fmt.Println(item.path + ": " + item.object)
	}

}

// SCHOOL
func SetCourseRoutes(router *mux.Router, middleware *jwtmiddleware.JWTMiddleware){
	router.HandleFunc("/courses/all", GetAllCourseData).Methods("GET")
	//router.Handle("/courses", middleware.Handler(GetCourses)).Methods("GET")

	router.HandleFunc("/courses", GetCourses).Methods("GET")

	router.HandleFunc("/courses", NewCourse).Methods("PUT")
	router.HandleFunc("/courses/{course_id}", GetCourse).Methods("GET")
	//router.HandleFunc("/courses/{course_id}/grades", GetCourseGrades).Methods("GET")



	router.HandleFunc("/courses/{course_id}/assignments", GetCourseAssignments).Methods("GET")

	//router.HandleFunc("/courses/{code}", UpdateCourse).Methods("POST")
	//router.HandleFunc("/courses/{code}", DeleteCourse).Methods("DELETE")
	router.HandleFunc("/courses/{course_id}/registrants", GetCourseRegistrants).Methods("GET")

	// Course Session Routes
	router.HandleFunc("/courses/{course_id}/sessions", GetCourseSessions).Methods("GET")
	router.HandleFunc("/courses/{course_id}/sessions", NewCourseSession).Methods("PUT")
	router.HandleFunc("/courses/{course_id}/sessions/{session_id}", GetCourseSession).Methods("GET")
	router.HandleFunc("/courses/{course_id}/sessions/{session_id}/assignments", GetSessionAssignments).Methods("GET")
	//router.HandleFunc("/courses/{code}/sessions/{session_number}", UpdateCourseSession).Methods("POST")
	//router.HandleFunc("/courses/{code}/sessions/{session_number}", DeleteCourseSession).Methods("DELETE")
}

func SetStudentRoutes(router *mux.Router, middleware *jwtmiddleware.JWTMiddleware){

	// Student Paths
	router.HandleFunc("/students", GetStudents).Methods("GET")
	router.HandleFunc("/users/{id}/student_info", GetStudent).Methods("GET")
	router.HandleFunc("/users/{id}/student", NewStudent).Methods("PUT")
}

// COMPANY
func SetEmployeeRoutes(router *mux.Router, middleware *jwtmiddleware.JWTMiddleware){

	// Student Paths
	router.HandleFunc("/employees", GetEmployees).Methods("GET")
	router.HandleFunc("/users/{id}/employee_info", GetEmployee).Methods("GET")
	router.HandleFunc("/users/{id}/employee", NewEmployee).Methods("PUT")
}

// SALES
func SetProductRoutes(router *mux.Router, middleware *jwtmiddleware.JWTMiddleware){

	// Product Paths
	router.HandleFunc("/products", GetProducts).Methods("GET")
	router.HandleFunc("/products/{id}", GetProduct).Methods("GET")
	router.HandleFunc("/products/{id}", NewProduct).Methods("PUT")

	// Order Paths
	router.HandleFunc("/orders", GetOrders).Methods("GET")
	router.HandleFunc("/orders/{id}", GetOrder).Methods("GET")
	router.HandleFunc("/orders/{id}/details", GetOrderDetails).Methods("GET")

}
