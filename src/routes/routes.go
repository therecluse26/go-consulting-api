package routes

import (
	"github.com/gorilla/mux"
	"github.com/auth0/go-jwt-middleware"
	"net/http"
	"github.com/spf13/viper"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"../database"
	"os"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"github.com/gorilla/handlers"
	"../controllers"
	"../config/mainconf"
)

type JsonRoutes struct {
	ListEndpoints bool `json:"listEndpoints"`
	Groups []Group `json:"group"`
}

type Group struct {
	Name string `json:"name"`
	Subgroup []SubGroup `json:"subgroup"`
}

type SubGroup struct {
	Name string `json:"name"`
	Routes []Route `json:"routes"`
}

type Route struct {
	Path string `json:"path"`
	Method string `json:"method"`
	EndpointType string `json:"endpoint-type"`
	Access []string `json:"access"`
	Query string `json:"query"`
	Function func(http.ResponseWriter, *http.Request) `json:"function"`
	Description string `json:"description"`
	Params []string `json:"params"`
}

type PassRoute struct {
	RouteStruct Route
}

var Router *mux.Router

func SpawnRouter(conf mainconf.Configuration){
	Router = mux.NewRouter().StrictSlash(true)

	// Handle all preflight requests
	Router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Access-Control-Allow-Origin, Origin, User-Agent, Referer, Cache-Control, X-header")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNoContent)
		return
	})

	Router.HandleFunc("/get-token", GetTokenHandler).Methods("GET")

	LoadRouteConfig(Router, controllers.JwtMiddleware)

	fmt.Println("Listening on port " + strconv.Itoa(conf.ApiPort))
	http.ListenAndServe(":"+strconv.Itoa(conf.ApiPort), handlers.LoggingHandler(os.Stdout, Router))
}


func DestroyRoutes(router *mux.Router){

	fmt.Println(router)


}

/**
 * Iterates over RouteData from `routes.json` file and sets routes accordingly
 */
func RegisterRoutes (Router *mux.Router, middleware *jwtmiddleware.JWTMiddleware, RouteData JsonRoutes) {

	for _, group := range RouteData.Groups {

		for _, subgroup := range group.Subgroup {

			for _, rt := range subgroup.Routes {

				RouteObj := &PassRoute{RouteStruct: rt}

				RouteObj.SetRoute(Router)

			}
		}
	}
}

/**
 * Simple wrapper to pass route info into appropriate function based on "endpoint-type"
 */
func (rt *PassRoute) SetRoute (router *mux.Router) {

	if rt.RouteStruct.EndpointType == "query" {

		router.HandleFunc(rt.RouteStruct.Path, rt.QueryRoute).Methods(rt.RouteStruct.Method)

	} else if rt.RouteStruct.EndpointType == "function" {

		FunctionRoute(rt.RouteStruct)
	}

}

/**
 * Builds query, executes result and builds JSON response
 */
func (rt *PassRoute) QueryRoute (w http.ResponseWriter, r *http.Request) {

	sql := database.Statement{ Sql: rt.RouteStruct.Query, Params: mux.Vars(r) }

	database.SelectAndReturnJson(sql, w)

}


func FunctionRoute (rt Route) {

	panic(Router.HandleFunc(rt.Path, rt.Function).Methods(rt.Method))

}

func LoadRouteConfig(router *mux.Router, middleware *jwtmiddleware.JWTMiddleware) {

	RouteData := ParseRoutesFile()
	RegisterRoutes(Router, middleware, RouteData)

	viper.SetConfigName("routes")
	viper.SetConfigType("json")
	viper.AddConfigPath("./")
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {

		fmt.Println("Config file changed:", e.Name)

		DestroyRoutes(router)

		RouteData := ParseRoutesFile()
		RegisterRoutes(router, middleware, RouteData)
	})

}


/**
 * Parses `routes.json` file and returns unmarshalled byte array
 */
func ParseRoutesFile () (JsonRoutes) {

	var RouteData JsonRoutes

	jsonFile, err := os.Open("./routes.json")

	if err != nil {
		panic(fmt.Sprintf("%s", "Could not open `routes.json` file"))
	} else {
		fmt.Println("Successfully opened `routes.json` file")
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &RouteData)

	return RouteData

}






