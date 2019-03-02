package routes

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type RouteConf struct {
	ListEndpoints string
	Group []struct{
		Name string
		SubGroup []struct{
			Name string
			Routes []struct{
				Path string
				HttpMethod string `yaml:"httpmethod"`
				Function string
				Access []string
				Query string
				Params []string
				Description string
			}
		}
	}
}

func LoadRoutesFromFile(RouteFile string) RouteConf {

	r := RouteConf{}

	rtConfRaw, err := ioutil.ReadFile(RouteFile); if err!=nil {log.Fatalf("error: %v", err)}

	err = yaml.Unmarshal(rtConfRaw, &r); if err!=nil {log.Fatalf("error: %v", err)}

	return r
}
