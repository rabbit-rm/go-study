package yaml

import (
	"os"
	"testing"

	yamlV3 "gopkg.in/yaml.v3"
)

func TestYaml(t *testing.T) {
	type name1 struct {
		Name1 string `yaml:"name1"`
		Name2 string `yaml:"name2"`
	}
	type name struct {
		Debug   bool   `yaml:"debug"`
		FileDir string `yaml:"fileDir"`
		LogDir  string `yaml:"logDir"`
		Name    name1  `yaml:"name"`
	}

	marshal, err := yamlV3.Marshal(&name{
		Debug:   false,
		FileDir: "D:\\\\web office\\\\officeFiles",
		LogDir:  "D:\\web office\\logs",
		Name: name1{
			Name1: "xxxxx",
			Name2: "xxxxxxx",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(marshal)
	os.WriteFile("W:\\GoProject\\private\\study\\yaml\\test.yaml", marshal, 0644)
}
