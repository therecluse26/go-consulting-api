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
	"strconv"
	"github.com/getsentry/raven-go"
)

// Initializes variables in global scope
var conf = mainconf.BuildConfig()
var Router *mux.Router
var ProtectedRouter *mux.Router

func init() {

	// Pulls config variables
	raven.SetDSN(conf.SentryHost)

	// Creates database connection
	database.DbConnection(conf)

	// Initializes router
	Router = mux.NewRouter().StrictSlash(true)

	ProtectedRouter = mux.NewRouter().StrictSlash(true)

	// Loads access policies from database on a loop every x seconds
	auth.LoadAccessPolicyLoopTimer(600)

	// Caches Azure auth keys on a timer
	auth.CacheAccessKeysTimer(300, conf.CacheMethod)


}

func main() {

	// Handle all preflight requests
	Router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Access-Control-Allow-Origin, Origin, User-Agent, Referer, Cache-Control, X-header")
		w.Header().Set("Access-Control-Allow-Origin", conf.AllowedOrigins[1])
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusNoContent)
		return
	})


	// Handle all preflight requests
	ProtectedRouter.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Access-Control-Allow-Origin, Origin, User-Agent, Referer, Cache-Control, X-header")
		w.Header().Set("Access-Control-Allow-Origin", conf.AllowedOrigins[1])
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusNoContent)
		return
	})

	// Initialize Routing Paths
	Router.HandleFunc("/", routes.GetStats).Methods("GET")

	// Set routes from individual routes files
	routes.SetCourseRoutes(Router, ProtectedRouter, controllers.JwtMiddleware)
	routes.SetUserRoutes(Router, ProtectedRouter, controllers.JwtMiddleware)
	routes.SetStudentRoutes(Router, ProtectedRouter, controllers.JwtMiddleware)
	routes.SetEmployeeRoutes(Router, ProtectedRouter, controllers.JwtMiddleware)
	routes.SetProductRoutes(Router, ProtectedRouter, controllers.JwtMiddleware)
	/*****************************/

	Router.HandleFunc("/authcallback", auth.AuthCallback).Methods("GET")
	Router.HandleFunc("/login", auth.LoginOrg).Methods("GET")
	Router.HandleFunc("/logout", auth.LogoutOrg).Methods("GET")

	corsObj := handlers.AllowedOrigins(conf.AllowedOrigins)

	// Initialize http server
	fmt.Println("Listening on port " + strconv.Itoa(conf.ApiPort))

	http.ListenAndServe(":"+strconv.Itoa(conf.ApiPort), handlers.CORS(corsObj)(Router))

}

