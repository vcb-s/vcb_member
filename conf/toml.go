package conf

import (
	"io/ioutil"

	"github.com/pelletier/go-toml"
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
		panic(err)
	}

	err = toml.Unmarshal(tomlFile, &Main)

	if err != nil {
		panic(err)
	}

}
