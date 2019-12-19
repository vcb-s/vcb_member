package models

import (
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

type conf struct {
	Database struct {
		Host  string
		Ports int
		User  string
		Pass  string
	}
	Jwt struct {
		Mac        string
		Encryption string
	}
}

// Conf 配置内容
var Conf conf

func init() {
	tomlFile, err := ioutil.ReadFile("./config.toml")

	if err != nil {
		panic(err)
	}

	err = toml.Unmarshal(tomlFile, &Conf)

	if err != nil {
		panic(err)
	}
}
