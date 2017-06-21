package main

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beego"
	"os"
	"os/signal"
	//"time"
	//"sync"
	"syscall"
)

func init() {
	beego.SetLogger("file", `{"filename":"logs/applong.log"}`)
	beego.SetLogFuncCall(true)

}

var islong bool
var db *sql.DB

func main() {
	defer beego.BeeLogger.Close()
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			beego.Info("version 0.2.0")
			return
		}
	}
	var err error
	islong, err = beego.AppConfig.Bool("islong")
	if err != nil {
		beego.Error(err)
		return
	}
	if islong {
		db = Conn()
		defer func() {
			fmt.Println("db closing")
			db.Close()
		}()
	}
	go func() {
		for {
			//beego.Debug("No. ---> ", i+1)
			err = db.Ping()
			if err != nil {
				beego.Error("ping() database failure!", err)
				db = Conn()
				continue
			}
			beego.Info("ping() database success!")
			err := Crud()
			if err != nil {
				beego.Error("Failure to obtain master information:", err)
				//time.Sleep(10000 * time.Millisecond)
				db = Conn()
			}
		}
	}()
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Kill)
	s := <-ch
	beego.Info(s)
}

func Conn() *sql.DB {
	ip := beego.AppConfig.String("ip")
	port := beego.AppConfig.String("port")
	username := beego.AppConfig.String("username")
	password := beego.AppConfig.String("password")
	database := beego.AppConfig.String("database")
	db, err := sql.Open("mysql", username+":"+password+"@tcp("+ip+":"+port+")/"+database+"?charset=utf8")
	if err != nil {
		beego.Error("open database connection string failure!", err)
		panic(err.Error())
	}
	beego.Info("open database connection string success!")
	//db.SetMaxOpenConns(2000)
	//db.SetMaxIdleConns(1000)
	return db
}
