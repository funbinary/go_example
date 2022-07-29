package main

import (
	"fmt"
	"github.com/funbinary/go_example/pkg/bfile"
	"log"

	"github.com/spf13/viper"
)

type Mysql struct {
	IP       string `mapstructure:ip`
	Port     int    `mapstructure:port`
	User     string
	Password string
	Database string
}

type Config struct {
	Mysql Mysql
}

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("ini")
	viper.AddConfigPath(bfile.SelfDir())
	err := viper.ReadInConfig()
	if err != nil {
		log.Panicln(err)
	}
	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(config.Mysql.IP)
	fmt.Println(config.Mysql.Port)
	fmt.Println(config.Mysql.User)
	fmt.Println(config.Mysql.Password)
	fmt.Println(config.Mysql.Database)
}
