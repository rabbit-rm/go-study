package config

import (
	"fmt"

	"github.com/rabbit-rm/rabbit/viperKit"
)

const path = "resources/config.yaml"

var (
	c Config
)

func init() {
	err := viperKit.UnmarshalFromFile(path, nil, &c)
	if err != nil {
		panic(err)
	}
}

func MySQLConf() MySQL {
	return c.MySQL
}

func Host() string {
	return fmt.Sprintf("%s:%d", c.IP, c.Port)
}
