package mainconf

import (
	"database/sql"
	"os"
	"encoding/json"
	"../../util"
)

var Dbconn *sql.DB

type Configuration struct {
	SqlHost string `json:"SqlHost"`
	SqlPort int `json:"SqlPort"`
	SqlUser string `json:"SqlUser"`
	SqlPass string `json:"SqlPass"`
	SqlDB string `json:"SqlDB"`
	ApiHost string `json:"ApiHost"`
	ApiPort int `json:"ApiPort"`
	AllowedOrigins []string `json:"AllowedOrigins"`
	CacheMethod string `json:"CacheMethod"`
	SentryHost string `json:"SentryHost"`
}

type AuthConfig struct {
	AuthHost string `json:"AuthHost"`
	AuthClientId string `json:"AuthClientId"`
	AuthSecret string `json:"AuthSecret"`
	AuthAudience string `json:"AuthAudience"`
	AuthTenant string `json:"AuthTenant"`
}

var configJson string

// Builds primary app configuration object
func BuildConfig() Configuration {

	configJson = os.Getenv("CFGJSON")

	var conf Configuration

	err := json.Unmarshal([]byte(configJson), &conf)
	if err != nil {
		util.ErrorHandler(err)
	}

	return conf

}

// Builds Auth configuration object
func GetAuthConfig() AuthConfig {

	configJson = os.Getenv("AUTHCFG")

	var authConf AuthConfig

	err := json.Unmarshal([]byte(configJson), &authConf)
	if err != nil {
		util.ErrorHandler(err)
	}

	return authConf
}

