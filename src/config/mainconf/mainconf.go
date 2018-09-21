package mainconf

import (
	"database/sql"
	"os"
	"encoding/json"
	"fmt"
)

var Dbconn *sql.DB

type Configuration struct {
	SqlHost string `json:"SqlHost"`
	SqlPort int `json:"SqlPort"`
	SqlUser string `json:"SqlUser"`
	SqlPass string `json:"SqlPass"`
	SqlDB string `json:"SqlDB"`
	ApiPort int `json:"ApiPort"`
	SentryHost string `json:"SentryHost"`
}

type AuthConfig struct {
	AuthHost string
	AuthSecret string
}

var configJson string
var configMap map[string]string

// Builds primary app configuration object
func BuildConfig() Configuration {

	configJson = os.Getenv("CFGJSON")

	var conf Configuration

	err := json.Unmarshal([]byte(configJson), &conf)
	if err != nil {
		fmt.Println("error:", err)
	}

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

