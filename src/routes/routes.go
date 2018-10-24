package routes

import (
	"github.com/gorilla/mux"
	"github.com/auth0/go-jwt-middleware"
)

// GENERAL
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
