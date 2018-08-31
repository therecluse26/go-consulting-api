package mainconf

import (
	"database/sql"
	"os"
	"strconv"
	"encoding/json"
	"encoding/base64"
)

var Dbconn *sql.DB

type Configuration struct {
	SqlHost string
	SqlPort int
	SqlUser string
	SqlPass string
	SqlDB string
	ApiPort int
}

type AuthConfig struct {
	AuthHost string
	AuthSecret string
}

var configJson []byte
var config64 string
var configMap map[string]string

// Builds primary app configuration object
func BuildConfig() Configuration {

	config64 = os.Getenv("gocfg64")

	configJson, _ = base64.StdEncoding.DecodeString(config64)

	json.Unmarshal(configJson, &configMap)

	apiPort, _ := strconv.Atoi(configMap["apiPort"])
	sqlPort, _ := strconv.Atoi(configMap["sqlPort"])

	conf := Configuration{
		SqlHost: configMap["sqlHost"],
		SqlPort: sqlPort,
		SqlUser: configMap["sqlUser"],
		SqlPass: configMap["sqlPass"],
		SqlDB: configMap["sqlDB"],
		ApiPort: apiPort,
	}

	/*fmt.Println(apiPort)
	fmt.Println(sqlPort)
	fmt.Println(configMap["sqlUser"])
	fmt.Println(configMap["sqlPass"])
	fmt.Println(configMap["sqlDB"])
	fmt.Println(configMap["sqlHost"])*/

	return conf

}

// Builds Auth configuration object
func GetAuthConfig() AuthConfig {

	authConf := AuthConfig{
		AuthHost: os.Getenv("AuthHost"),
		AuthSecret: os.Getenv("AuthSecret"),
	}

	return authConf
}

