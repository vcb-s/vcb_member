package conf

import (
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

type conf struct {
	Database struct {
		Host   string
		Port   int
		User   string
		Pass   string
		Dbname string
	}
	Jwt struct {
		Mac        string
		Encryption string
	}
	Server struct {
		Port int
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
