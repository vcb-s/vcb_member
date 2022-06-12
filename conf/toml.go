package conf

import (
	"io/ioutil"
	"log"

	"github.com/pelletier/go-toml/v2"
)

type wpAuth struct {
	ClientID  string
	ClientSec string
}

type conf struct {
	Debug    bool
	Database struct {
		Host   string
		Port   int
		User   string
		Pass   string
		Dbname string
	}
	Redis struct {
		Host string
		Port int
		Pass string
	}
	Jwt struct {
		Mac string
	}
	Server struct {
		Port int
	}
	Third struct {
		Wp wpAuth
	}
}

// Main 配置值
var Main conf

func init() {
	tomlFile, err := ioutil.ReadFile("./config.toml")

	if err != nil {
		log.Print("Failed to open config.toml file")
		log.Panic(err)
	}

	err = toml.Unmarshal(tomlFile, &Main)

	if err != nil {
		log.Print("Failed to open parse.toml file")
		log.Panic(err)
	}
}
