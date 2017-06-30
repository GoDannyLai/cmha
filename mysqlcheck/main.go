package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/astaxie/beego"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/upmio/cmha/mysqlcheck/check"
)

var (
	UserFlag     = flag.String("user", "", "usage: the mysql user")
	PasswordFlag = flag.String("password", "", "usage: the mysql password")
)
var (
	client *consulapi.Client
	err    error
)

func main() {
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			fmt.Println("version 1.2.0")
			return
		}
	}
	flag.Parse()
	if *UserFlag == "" || *PasswordFlag == "" {
		os.Exit(2)
	}
	config := GetConsulConfig()
	if client == nil {
		client, err = GetClient(config)
		if err != nil {
			fmt.Print(err)
			os.Exit(2)
		}
	}
	user := *UserFlag
	password := *PasswordFlag
	check.IsPingType(user, password, client)
}

func GetConsulConfig() *consulapi.Config {
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
		Address:    "127.0.0.1:8500",
		WaitTime:   1 * time.Second,
	}
	return config
}

func GetClient(config *consulapi.Config) (*consulapi.Client, error) {
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	return client, nil
}
