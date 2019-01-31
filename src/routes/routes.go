package routes

import (
	"github.com/gorilla/mux"
	"github.com/auth0/go-jwt-middleware"
	"github.com/urfave/negroni"
	"../auth"
	"net/http"
	"fmt"
)

// GENERAL
func SetUserRoutes(router *mux.Router, protectedRouter *mux.Router, middleware *jwtmiddleware.JWTMiddleware){

	// User Paths
	router.HandleFunc("/users", GetUsers).Methods("GET")
	router.HandleFunc("/users", NewUser).Methods("PUT")

	protectedRouter.HandleFunc("/users/{id:[0-9]+}", GetUser).Methods("GET")
	protectedRouter.HandleFunc("/users/{id:[0-9]+}/info", GetUserInfo).Methods("GET")
	protectedRouter.HandleFunc("/users/{id:[0-9]+}/roles", GetUserRoles).Methods("GET")

	router.PathPrefix("/users/{id:[0-9]+}").Handler(negroni.New(
		negroni.HandlerFunc(auth.ProtectedEndpoint),
		negroni.Wrap(protectedRouter),
	))

	// Role Paths
	protectedRouter.HandleFunc("/roles", GetRoles).Methods("GET")
	protectedRouter.HandleFunc("/roles/{id:[0-9]+}", GetRole).Methods("GET")
	protectedRouter.HandleFunc("/roles/{id:[0-9]+}/users", GetRoleUsers).Methods("GET")

	router.PathPrefix("/roles").Handler(negroni.New(
		negroni.HandlerFunc(auth.ProtectedEndpoint),
		negroni.Wrap(protectedRouter),
	))
}

// SCHOOL
func SetCourseRoutes(router *mux.Router, protectedRouter *mux.Router, middleware *jwtmiddleware.JWTMiddleware){

	router.HandleFunc("/courses/all", GetAllCourseData).Methods("GET")
	router.HandleFunc("/courses", GetCourses).Methods("GET")
	protectedRouter.HandleFunc("/courses", NewCourse).Methods("PUT")
		router.Path("/courses").Handler(negroni.New(
			negroni.HandlerFunc(auth.ProtectedEndpoint),
			negroni.Wrap(protectedRouter),
		))
	router.HandleFunc("/courses/{course_id:[0-9]+}", GetCourse).Methods("GET")

	protectedRouter.HandleFunc("/courses/{course_id:[0-9]+}/grades", GetCourseGrades).Methods("GET")
		router.Path("/courses/{course_id:[0-9]+}/grades").Handler(negroni.New(
			negroni.HandlerFunc(auth.ProtectedEndpoint),
			negroni.Wrap(protectedRouter),
		))

	protectedRouter.HandleFunc("/courses/{course_id:[0-9]+}/assignments", GetCourseAssignments).Methods("GET")
		router.Path("/courses/{course_id:[0-9]+}/assignments").Handler(negroni.New(
			negroni.HandlerFunc(auth.ProtectedEndpoint),
			negroni.Wrap(protectedRouter),
		))

	protectedRouter.HandleFunc("/courses/{course_id:[0-9]+}/registrants", GetCourseRegistrants).Methods("GET")

	// Course Session Routes
	router.HandleFunc("/courses/{course_id:[0-9]+}/sessions", GetCourseSessions).Methods("GET")
	protectedRouter.HandleFunc("/courses/{course_id:[0-9]+}/sessions", NewCourseSession).Methods("PUT")
	//protectedRouter.HandleFunc("/courses/{course_id:[0-9]+}/sessions", UpdateCourseSession).Methods("POST")
	router.HandleFunc("/courses/{course_id:[0-9]+}/sessions/{session_id:[0-9]+}", GetCourseSession).Methods("GET")


	protectedRouter.HandleFunc("/courses/{course_id:[0-9]+}/sessions/{session_id:[0-9]+}/assignments", GetSessionAssignments).Methods("GET")
		router.Path("/courses/{course_id:[0-9]+}/sessions/{session_id:[0-9]+}/assignments").Handler(negroni.New(
			negroni.HandlerFunc(auth.ProtectedEndpoint),
			negroni.Wrap(protectedRouter),
		))

}

func SetStudentRoutes(router *mux.Router, protectedRouter *mux.Router, middleware *jwtmiddleware.JWTMiddleware){

	// Student Paths
	router.HandleFunc("/students", GetStudents).Methods("GET")
	router.HandleFunc("/users/{id:[0-9]+}/student_info", GetStudent).Methods("GET")
	router.HandleFunc("/users/{id:[0-9]+}/student", NewStudent).Methods("PUT")
}

// COMPANY
func SetEmployeeRoutes(router *mux.Router, protectedRouter *mux.Router, middleware *jwtmiddleware.JWTMiddleware){

	// Student Paths
	router.HandleFunc("/employees", GetEmployees).Methods("GET")
	router.HandleFunc("/users/{id:[0-9]+}/employee_info", GetEmployee).Methods("GET")
	router.HandleFunc("/users/{id:[0-9]+}/employee", NewEmployee).Methods("PUT")
}

// SALES
func SetProductRoutes(router *mux.Router, protectedRouter *mux.Router, middleware *jwtmiddleware.JWTMiddleware){

	// Product Paths
	router.HandleFunc("/products", GetProducts).Methods("GET")
	router.HandleFunc("/products/{id:[0-9]+}", GetProduct).Methods("GET")
	router.HandleFunc("/products/{id:[0-9]+}", NewProduct).Methods("PUT")

	// Order Paths
	router.HandleFunc("/orders", GetOrders).Methods("GET")
	router.HandleFunc("/orders/{id:[0-9]+}", GetOrder).Methods("GET")
	router.HandleFunc("/orders/{id:[0-9]+}/details", GetOrderDetails).Methods("GET")
}


func NotFoundHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 - Not Found")
}

// https://login.microsoftonline.com/common/adminconsent?client_id=c7ba4700-1b55-4563-b066-9d103d59efcc&state=12345&redirect_uri=http://localhost:9988