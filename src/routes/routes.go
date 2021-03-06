package routes

import (
	"github.com/gorilla/mux"
	"github.com/auth0/go-jwt-middleware"
	"github.com/urfave/negroni"
	"github.com/therecluse26/fortisure-api/src/auth"
	"net/http"
	"fmt"
	"github.com/therecluse26/fortisure-api/src/routes/school"
	"github.com/therecluse26/fortisure-api/src/routes/sales"
	"github.com/therecluse26/fortisure-api/src/routes/general"
	"github.com/therecluse26/fortisure-api/src/routes/company"
	"github.com/therecluse26/fortisure-api/src/routes/consulting"
	"github.com/therecluse26/fortisure-api/src/util"
)

// GENERAL
func SetGeneralRoutes(router *mux.Router, protectedRouter *mux.Router, middleware *jwtmiddleware.JWTMiddleware){
	router.HandleFunc("/", general.GetStats).Methods("GET")

	router.HandleFunc("/login", general.Login).Methods("GET")

}

// User-related routes
func SetUserRoutes(router *mux.Router, protectedRouter *mux.Router, middleware *jwtmiddleware.JWTMiddleware){

	// User Paths
	router.HandleFunc("/users", general.GetUsers).Methods("GET")
	router.HandleFunc("/users", general.NewUser).Methods("PUT")

	protectedRouter.HandleFunc("/users/{id:[0-9]+}", general.GetUser).Methods("GET")
	protectedRouter.HandleFunc("/users/{id:[0-9]+}/info", general.GetUserInfo).Methods("GET")
	protectedRouter.HandleFunc("/users/{id:[0-9]+}/roles", general.GetUserRoles).Methods("GET")

	router.PathPrefix("/users/{id:[0-9]+}").Handler(negroni.New(
		negroni.HandlerFunc(auth.ProtectedEndpoint),
		negroni.Wrap(protectedRouter),
	))

	// Role Paths
	protectedRouter.HandleFunc("/roles", general.GetRoles).Methods("GET")
	protectedRouter.HandleFunc("/roles/{id:[0-9]+}", general.GetRole).Methods("GET")
	protectedRouter.HandleFunc("/roles/{id:[0-9]+}/users", general.GetRoleUsers).Methods("GET")

	router.PathPrefix("/roles").Handler(negroni.New(
		negroni.HandlerFunc(auth.ProtectedEndpoint),
		negroni.Wrap(protectedRouter),
	))
}

// SCHOOL
func SetCourseRoutes(router *mux.Router, protectedRouter *mux.Router, middleware *jwtmiddleware.JWTMiddleware){

	router.HandleFunc("/courses/all", school.GetAllCourseData).Methods("GET")
	router.HandleFunc("/courses", school.GetCourses).Methods("GET")
	protectedRouter.HandleFunc("/courses", school.NewCourse).Methods("PUT")
		router.Path("/courses").Handler(negroni.New(
			negroni.HandlerFunc(auth.ProtectedEndpoint),
			negroni.Wrap(protectedRouter),
		))
	router.HandleFunc("/courses/{course_id:[0-9]+}", school.GetCourse).Methods("GET")

	protectedRouter.HandleFunc("/courses/{course_id:[0-9]+}/grades", school.GetCourseGrades).Methods("GET")
		router.Path("/courses/{course_id:[0-9]+}/grades").Handler(negroni.New(
			negroni.HandlerFunc(auth.ProtectedEndpoint),
			negroni.Wrap(protectedRouter),
		))

	protectedRouter.HandleFunc("/courses/{course_id:[0-9]+}/assignments", school.GetCourseAssignments).Methods("GET")

	protectedRouter.HandleFunc("/courses/{course_id:[0-9]+}/registrants", school.GetCourseRegistrants).Methods("GET")

	// Course Session Routes
	router.HandleFunc("/courses/{course_id:[0-9]+}/sessions", school.GetCourseSessions).Methods("GET")
	protectedRouter.HandleFunc("/courses/{course_id:[0-9]+}/sessions", school.NewCourseSession).Methods("PUT")
	//protectedRouter.HandleFunc("/courses/{course_id:[0-9]+}/sessions", school.UpdateCourseSession).Methods("POST")
	router.HandleFunc("/courses/{course_id:[0-9]+}/sessions/{session_id:[0-9]+}", school.GetCourseSession).Methods("GET")


	protectedRouter.HandleFunc("/courses/{course_id:[0-9]+}/sessions/{session_id:[0-9]+}/assignments", school.GetSessionAssignments).Methods("GET")

}

func SetStudentRoutes(router *mux.Router, protectedRouter *mux.Router, middleware *jwtmiddleware.JWTMiddleware){

	// Student Paths
	router.HandleFunc("/students", school.GetStudents).Methods("GET")
	router.HandleFunc("/users/{id:[0-9]+}/student_info", school.GetStudent).Methods("GET")
	router.HandleFunc("/users/{id:[0-9]+}/student", school.NewStudent).Methods("PUT")
}

// COMPANY
func SetEmployeeRoutes(router *mux.Router, protectedRouter *mux.Router, middleware *jwtmiddleware.JWTMiddleware){

	// Student Paths
	router.HandleFunc("/employees", company.GetEmployees).Methods("GET")
	router.HandleFunc("/users/{id:[0-9]+}/employee_info", company.GetEmployee).Methods("GET")
	router.HandleFunc("/users/{id:[0-9]+}/employee", company.NewEmployee).Methods("PUT")
}

// SALES
func SetProductRoutes(router *mux.Router, protectedRouter *mux.Router, middleware *jwtmiddleware.JWTMiddleware){

	// Product Paths
	router.HandleFunc("/products", sales.GetProducts).Methods("GET")
	router.HandleFunc("/products/{id:[0-9]+}", sales.GetProduct).Methods("GET")
	router.HandleFunc("/products/{id:[0-9]+}", sales.NewProduct).Methods("PUT")

	// Order Paths
	router.HandleFunc("/orders", sales.GetOrders).Methods("GET")
	router.HandleFunc("/orders/{id:[0-9]+}", sales.GetOrder).Methods("GET")
	router.HandleFunc("/orders/{id:[0-9]+}/details", sales.GetOrderDetails).Methods("GET")
}

// CONSULTING
func SetConsultingRoutes(router *mux.Router, protectedRouter *mux.Router, middleware *jwtmiddleware.JWTMiddleware){

	// Product Paths
	router.HandleFunc("/projects", consulting.GetAllProjects).Methods("GET")
	router.HandleFunc("/projects/{id:[0-9]+}", consulting.GetProject).Methods("GET")

}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	_, err := fmt.Fprintf(w, "404 - Not Found")
	if err != nil {
		util.ErrorHandler(err)
	}
}

// https://login.microsoftonline.com/common/adminconsent?client_id=c7ba4700-1b55-4563-b066-9d103d59efcc&state=12345&redirect_uri=http://localhost:9988