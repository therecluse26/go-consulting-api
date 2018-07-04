package config

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

func BuildConfig() Configuration {

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
