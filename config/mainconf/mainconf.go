package mainconf

import (
	"github.com/spf13/viper"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"database/sql"
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

// Initializes config file opener
func InitConfReader() {
	viper.SetDefault("Environment", "development")
	viper.SetDefault("ConfigDir", "config")
	viper.SetConfigName("config." + viper.GetString("Environment"))
	viper.AddConfigPath(viper.GetString("ConfigDir"))

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
}

// Builds primary app configuration object
func BuildConfig() Configuration {

	InitConfReader()

	conf := Configuration{
		SqlHost: viper.GetString("sqlHost"),
		SqlPort: viper.GetInt("sqlPort"),
		SqlUser: viper.GetString("sqlUser"),
		SqlPass: viper.GetString("sqlPass"),
		SqlDB: viper.GetString("sqlDB"),
		ApiPort: viper.GetInt("apiPort"),
	}

	return conf

}

// Builds Auth configuration object
func GetAuthConfig() AuthConfig {

	InitConfReader()

	authConf := AuthConfig{
		AuthHost: viper.GetString("AuthHost"),
		AuthSecret: viper.GetString("AuthSecret"),
	}

	return authConf
}

