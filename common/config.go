package common

import (
	"github.com/Terry-Mao/goconf"
	"fmt"
)


type Config struct {
	Port int `goconf:"http:port"`

	MongDB_Addr string `goconf:"db:mongodb-addr"`
	MongDB_DB string `goconf:"db:mongodb-db"`

	LOG_Debug bool `goconf:"log:debug"`
}

var myConfig *Config


func InitConfig(configFile string) *Config{
	conf := goconf.New()

	if err := conf.Parse(configFile);err !=nil{
		panic(err)
	}

	myConfig := &Config{}

	if err := conf.Unmarshal(myConfig);err!=nil{
		panic(err)
	}

	fmt.Printf("config is: %#v\n", myConfig)
	return myConfig
}