package main

import (
	"time"
	//	"bufio"
	"encoding/json"
	"io/ioutil"
	//	"bytes"
	"fmt"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	//"gopkg.in/yaml.v2"
	//	"io"
	//	"os"
	"strings"
)

type Transactions struct {
	Trans []Transaction `json:"transactions"`
}

type Transaction struct {
	Sql []string `json:"sqls"`
}

func Crud() error {
	if !islong {
		db = Conn()
		defer func() {
			fmt.Println("db closing")
			db.Close()
		}()
	}

	affairs := beego.AppConfig.String("affairs")
	//f, err := os.Open(affairs)
	//defer f.Close()
	var trans Transactions
	content, err := ioutil.ReadFile(affairs)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, &trans)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, tran := range trans.Trans {
		begin := time.Now()
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		for _, sql := range tran.Sql {
			if strings.HasPrefix(sql, "select") {
				/**/
				_, err := db.Query(sql)
				if err != nil {
					beego.Error("select data failure", err)
					//tx.Rollback()
					return err
				}
				beego.Info("Executing:", sql)

				/*			for raw.Next() {
								var id int
								var names string
								var created_at string
								err = raw.Scan(&id, &names, &created_at)
								if err != nil {
									beego.Error("scan() traversal data failure!", err)
								}
								//      		beego.Info("scan() traversal data success!")
								//fmt.Println("out", id, names, created_at)
							}
				*/
			} else {
				_, err = tx.Exec(sql)
				if err != nil {
					beego.Error(err, "[sql:]", sql)
					tx.Rollback()
					return err
				}
				beego.Info("Executing:", sql)
			}
		}

		if err := tx.Commit(); err != nil {
			beego.Error("Submit data failure", err)
			tx.Rollback()
			return err
		}
		beego.Alert("transaction use:", time.Now().Sub(begin).String())

	}
	return nil
}
