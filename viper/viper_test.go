package viper

import (
	"fmt"
	"log"
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func TestDefaultValues(t *testing.T) {
	viper.SetDefault("username", "rabbit.rm")
	viper.SetDefault("password", "rabbit.rm")
	viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})
	printAllSetting(viper.GetViper())
}

func TestReadAndWriteFromConfigFile(t *testing.T) {
	viper.SetConfigFile("./config.yaml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("read err:%v\n", err)
	}
	fmt.Printf("password:%s\n", viper.GetString("password"))
	printAllSetting(viper.GetViper())
	viper.Set("password", "123456789")
	err = viper.WriteConfig()
	if err != nil {
		log.Fatalf("write err:%v\n", err)
	}
}

func TestWatchConfig(t *testing.T) {
	viper.SetConfigFile("./config.yaml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("read err:%v\n", err)
	}
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.Name, in.Op)
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatalf("read err:%v\n", err)
		}
		printAllSetting(viper.GetViper())
	})
	printAllSetting(viper.GetViper())
	viper.WatchConfig()
	select {}
}
