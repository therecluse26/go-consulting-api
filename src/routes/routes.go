package routes

import (
	"../auth"
	"./company"
	"./consulting"
	"./general"
	"./sales"
	"./school"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"net/http"
	"reflect"

)

type RouteType = func(w http.ResponseWriter, r *http.Request)
type Loader = struct{
	funcRef func(w http.ResponseWriter, r *http.Request)
}

func (routes *RouteConf) RegisterRoutes(router *mux.Router, protectedRouter *mux.Router, middleware *jwtmiddleware.JWTMiddleware){

	for _, group := range routes.Group {

		for _, subGroup := range group.SubGroup {

			for _, route := range subGroup.Routes {

				loader := Loader{}
				//rt := RouteType(func(http.ResponseWriter, *http.Request){})

				/*var inputs = make([]reflect.Value, len(route.Params))

				for i, _ := range inputs {
					inputs[i] = reflect.ValueOf(inputs[i])
				}*/

				fmt.Println(route.Function)

				reflect.ValueOf(&loader.funcRef).MethodByName(route.Function)

				router.HandleFunc(route.Path, loader.funcRef).Methods(route.HttpMethod)

			}
		}
	}

}





// GENERAL
func SetGeneralRoutes(router *mux.Router, protectedRouter *mux.Router, middleware *jwtmiddleware.JWTMiddleware){
	router.HandleFunc("/", general.GetStats).Methods("GET")
}

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
		router.Path("/courses/{course_id:[0-9]+}/assignments").Handler(negroni.New(
			negroni.HandlerFunc(auth.ProtectedEndpoint),
			negroni.Wrap(protectedRouter),
		))

	protectedRouter.HandleFunc("/courses/{course_id:[0-9]+}/registrants", school.GetCourseRegistrants).Methods("GET")

	// Course Session Routes
	router.HandleFunc("/courses/{course_id:[0-9]+}/sessions", school.GetCourseSessions).Methods("GET")
	protectedRouter.HandleFunc("/courses/{course_id:[0-9]+}/sessions", school.NewCourseSession).Methods("PUT")
	//protectedRouter.HandleFunc("/courses/{course_id:[0-9]+}/sessions", school.UpdateCourseSession).Methods("POST")
	router.HandleFunc("/courses/{course_id:[0-9]+}/sessions/{session_id:[0-9]+}", school.GetCourseSession).Methods("GET")


	protectedRouter.HandleFunc("/courses/{course_id:[0-9]+}/sessions/{session_id:[0-9]+}/assignments", school.GetSessionAssignments).Methods("GET")
		router.Path("/courses/{course_id:[0-9]+}/sessions/{session_id:[0-9]+}/assignments").Handler(negroni.New(
			negroni.HandlerFunc(auth.ProtectedEndpoint),
			negroni.Wrap(protectedRouter),
		))

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
	fmt.Fprintf(w, "404 - Not Found")
}

// https://login.microsoftonline.com/common/adminconsent?client_id=c7ba4700-1b55-4563-b066-9d103d59efcc&state=12345&redirect_uri=http://localhost:9988