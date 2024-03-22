package viper

import (
	"fmt"

	"github.com/spf13/viper"
)

func printAllSetting(v *viper.Viper) {
	settings := v.AllSettings()
	for k, v := range settings {
		fmt.Printf("key:%v,value:%v\n", k, v)
	}
}
