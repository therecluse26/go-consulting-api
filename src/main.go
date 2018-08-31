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
	"os"
	"strconv"
)

// Initializes variables in global scope
var conf mainconf.Configuration
var AuthConf mainconf.AuthConfig
var Router *mux.Router
var RouteMap map[int]string
var RouteCount = 1

func init() {

	// Pulls config variables
	conf = mainconf.BuildConfig()

	/*AuthConf = mainconf.GetAuthConfig()

	os.Setenv("AuthHost", AuthConf.AuthHost)
	os.Setenv("AuthSecret", AuthConf.AuthSecret)*/


	// Creates database connection
	database.DbConnection(conf)

	Router = mux.NewRouter().StrictSlash(true)

}

func main() {

	// Handle all preflight requests
	Router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Access-Control-Allow-Origin, Origin, User-Agent, Referer, Cache-Control, X-header")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNoContent)
		return
	})

	/* Initialize Routing Paths */

	//router.HandleFunc("/", ViewRoutes).Methods("GET")

	Router.HandleFunc("/", routes.GetStats).Methods("GET")

	Router.HandleFunc("/get-token", routes.GetTokenHandler).Methods("GET")
	routes.SetCourseRoutes(Router, controllers.JwtMiddleware)
	routes.SetUserRoutes(Router, controllers.JwtMiddleware)
	routes.SetStudentRoutes(Router, controllers.JwtMiddleware)
	routes.SetEmployeeRoutes(Router, controllers.JwtMiddleware)
	routes.SetProductRoutes(Router, controllers.JwtMiddleware)
	/*****************************/

	fmt.Println(conf.SqlDB)
	fmt.Println(conf.SqlHost)
	fmt.Println(conf.SqlPort)
	fmt.Println(conf.SqlUser)

	fmt.Println("Listening on port " + strconv.Itoa(conf.ApiPort))
	http.ListenAndServe(":"+strconv.Itoa(conf.ApiPort), handlers.LoggingHandler(os.Stdout, Router))

}
