package main

import (
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/getsentry/raven-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/srikrsna/security-headers"
	"github.com/therecluse26/fortisure-api/src/auth"
	"github.com/therecluse26/fortisure-api/src/config/mainconf"
	"github.com/therecluse26/fortisure-api/src/controllers"
	"github.com/therecluse26/fortisure-api/src/database"
	"github.com/therecluse26/fortisure-api/src/routes"
	"github.com/therecluse26/fortisure-api/src/util"
	"net/http"
	"strconv"
)

// Initializes variables in global scope
var conf = mainconf.BuildConfig()
var Router *mux.Router
var ProtectedRouter *mux.Router

// CSP Middleware for blocking favicon requests
var csp = &secure.CSP{
	Value: `img-src 'none'`,
}

func init() {

	// Initializes sentry connection
	err := raven.SetDSN(conf.SentryHost)
	if err != nil {
		util.ErrorHandler(err)
	}

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
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Access-Control-Allow-Origin, Origin, User-Agent, Referer, Cache-Control, X-header, Content-Security-Policy: img-src none")
		w.Header().Set("Access-Control-Allow-Origin", conf.AllowedOrigins[1])
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusNoContent)
		return
	})


	// Handle all preflight requests for protected router
	ProtectedRouter.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Access-Control-Allow-Origin, Origin, User-Agent, Referer, Cache-Control, X-header, Content-Security-Policy: img-src none")
		w.Header().Set("Access-Control-Allow-Origin", conf.AllowedOrigins[1])
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.WriteHeader(http.StatusNoContent)
		return
	})


	// Set routes from individual routes files
	routes.SetGeneralRoutes(Router, ProtectedRouter, controllers.JwtMiddleware)
	routes.SetCourseRoutes(Router, ProtectedRouter, controllers.JwtMiddleware)
	routes.SetUserRoutes(Router, ProtectedRouter, controllers.JwtMiddleware)
	routes.SetStudentRoutes(Router, ProtectedRouter, controllers.JwtMiddleware)
	routes.SetEmployeeRoutes(Router, ProtectedRouter, controllers.JwtMiddleware)
	routes.SetProductRoutes(Router, ProtectedRouter, controllers.JwtMiddleware)
	routes.SetConsultingRoutes(Router, ProtectedRouter, controllers.JwtMiddleware)
	/*****************************/

	corsObj := handlers.AllowedOrigins(conf.AllowedOrigins)

	// Initialize http server
	fmt.Println("Listening on port " + strconv.Itoa(conf.ApiPort))
	err := http.ListenAndServe(":"+strconv.Itoa(conf.ApiPort), csp.Middleware()(handlers.CORS(corsObj)(Router)))
	if err != nil {
		panic(err)
	}

}

