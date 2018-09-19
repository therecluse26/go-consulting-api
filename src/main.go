package main

import (
	_ "github.com/denisenkom/go-mssqldb"
	"./config/mainconf"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"./database"
	"./routes"
)

// Initializes variables in global scope
var conf mainconf.Configuration
var AuthConf mainconf.AuthConfig


func init() {

	// Pulls config variables
	conf = mainconf.BuildConfig()

	/*AuthConf = mainconf.GetAuthConfig()

	os.Setenv("AuthHost", AuthConf.AuthHost)
	os.Setenv("AuthSecret", AuthConf.AuthSecret)*/


	// Creates database connection
	database.DbConnection(conf)



}

func main() {


	routes.SpawnRouter(conf)

	/* Initialize Routing Paths */

	//router.HandleFunc("/", ViewRoutes).Methods("GET")




}
