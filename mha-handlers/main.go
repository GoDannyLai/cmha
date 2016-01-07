package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"os"
	"time"
)

func init() {
	beego.SetLogger("file", `{"filename":"logs/mha-handlers.log"}`)
	beego.SetLogFuncCall(true)
}

func main() {
	defer beego.BeeLogger.Close()
	defer time.Sleep(100 * time.Millisecond)
	args := os.Args[1:]
	for _, arg := range args {
		if arg == "-v" || arg == "--version" {
			fmt.Println("version 1.0.0")
			return
		} else {
			return
		}
	}
	SessionAndChecks()
}
