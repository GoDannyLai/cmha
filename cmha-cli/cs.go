package main

import (
	"fmt"
	"github.com/astaxie/beego"
)

func showlist(args ...string){
	csip :=beego.AppConfig.Strings("cmha-cs-ip")
	_data := make([][]string, len(csip))
	for i,_ := range csip {
		node_info := make([]string, 1)
		node_info[0] = csip[i]
		_data = append(_data, node_info)
	}
	
	_th := []string{
		"cs ip",
	}
	TableRender(_th, _data, ALIGN_CENTRE)
	fmt.Println("")
	
}
