package config

import (
	"fmt"

	"github.com/rabbit-rm/rabbit/viperToolkit"
)

const path = "resources/config.yaml"

var (
	c Config
)

func init() {
	err := viperToolkit.UnmarshalFromFile(path, nil, &c)
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

func UploadConf() Upload {
	return c.Upload
}
