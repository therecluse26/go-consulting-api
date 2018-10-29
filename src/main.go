package main

import (
	_ "github.com/denisenkom/go-mssqldb"
	"./config/mainconf"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"fmt"
	"net/http"
	"./routes"
	"./database"
	"./controllers"
	"./auth"
	"os"
	"strconv"
	"github.com/getsentry/raven-go"
	"github.com/urfave/negroni"
)

// Initializes variables in global scope
var conf = mainconf.BuildConfig()
var Router *mux.Router


func init() {

	// Pulls config variables

	raven.SetDSN(conf.SentryHost)

	/*AuthConf = mainconf.GetAuthConfig()

	os.Setenv("AuthHost", AuthConf.AuthHost)
	os.Setenv("AuthSecret", AuthConf.AuthSecret)*/

	// Creates database connection
	database.DbConnection(conf)

	Router = mux.NewRouter().StrictSlash(true)

}

func main() {

	// Initialize Middleware Handlers
	/*ProtectedRead = negroni.New(negroni.HandlerFunc(auth.ProtectedEndpoint))
	ProtectedCreate = negroni.New(negroni.HandlerFunc(auth.ProtectedEndpoint))
	ProtectedUpdate = negroni.New(negroni.HandlerFunc(auth.ProtectedEndpoint))
	ProtectedDelete = negroni.New(negroni.HandlerFunc(auth.ProtectedEndpoint))*/


	// Handle all preflight requests
	Router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Access-Control-Allow-Origin, Origin, User-Agent, Referer, Cache-Control, X-header")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNoContent)
		return
	})

	/* Initialize Routing Paths */

	Router.HandleFunc("/", routes.GetStats).Methods("GET")

	routes.SetCourseRoutes(Router, controllers.JwtMiddleware)
	routes.SetUserRoutes(Router, controllers.JwtMiddleware)
	routes.SetStudentRoutes(Router, controllers.JwtMiddleware)
	routes.SetEmployeeRoutes(Router, controllers.JwtMiddleware)
	routes.SetProductRoutes(Router, controllers.JwtMiddleware)
	/*****************************/

	Router.HandleFunc("/authcallback", auth.AuthCallback).Methods("GET")
	//Router.HandleFunc("/tokenvalidate", auth.ValidateToken).Methods("GET")
	Router.HandleFunc("/login", auth.LoginOrg).Methods("GET")
	Router.HandleFunc("/logout", auth.LogoutOrg).Methods("GET")
	//Router.Handle("/protected", negroni.(auth.ProtectedEndpoint, routes.GetStats) ).Methods("GET")

	//Router.HandleFunc("/protected", routes.GetStats).Methods("GET")
	/*Router.Use(ProtectedRead)

	//Actually set routes
	ProtectedRead.UseHandler(Router)
	ProtectedCreate.UseHandler(Router)
	ProtectedUpdate.UseHandler(Router)
	ProtectedDelete.UseHandler(Router)*/



	protectedRouter := mux.NewRouter().PathPrefix("/protected").Subrouter().StrictSlash(true)
	protectedRouter.HandleFunc("/", routes.GetStats) // "/subpath/"
	protectedRouter.HandleFunc("/{course_id}", routes.GetCourse) // "/subpath/:id"
	// "/subpath" is necessary to ensure the subRouter and main router linkup

	Router.PathPrefix("/protected").Handler(negroni.New(
		negroni.HandlerFunc(auth.ProtectedEndpoint),
		negroni.Wrap(protectedRouter),
	))


	// Updates access policies from database on a loop every x seconds
	auth.LoadAccessPolicyLoopTimer(30)
	auth.CacheAccessKeysTimer(30, "local_env")

	fmt.Println("Listening on port " + strconv.Itoa(conf.ApiPort))

	http.ListenAndServe(":"+strconv.Itoa(conf.ApiPort), handlers.LoggingHandler(os.Stdout, Router))

}
