package check

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"

	consulapi "github.com/hashicorp/consul/api"
)

func IsPingType(user, password string) {
	host, port, checktime_string, timeout, defaultDb, ping_type := GetConfig()
	servicename := beego.AppConfig.String("servicename")
	if ping_type == "select,replication" || ping_type == "select" {
		TrySelectCheckTime(user, password, host, port, defaultDb, checktime_string, ping_type, timeout, servicename)
	} else if ping_type == "update,replication" || ping_type == "update" {
		TryUpdateCheckTime(user, password, host, port, defaultDb, checktime_string, ping_type, timeout, servicename)
	} else {
		fmt.Print("Configuration error")
		os.Exit(2)
	}

}

func TrySelectCheckTime(user, password, host, port, defaultDb, checktime_string, ping_type, timeout, servicename string) {
	checktime, _ := strconv.Atoi(checktime_string)
	for {
		if checktime == 0 {
			break
		} else {
			checktime--
			MYSQL_OK := SelectCheckMysqlHealth(user, password, host, port, defaultDb, timeout, checktime)
			if MYSQL_OK == 0 {
				if strings.Contains(ping_type, "replication") {
					isyes, err := ShowSlave(user, password, host, port, defaultDb, timeout)
					if err != nil {
						fmt.Print(err)
						os.Exit(2)
					}
					if isyes == "Yes" {
						fmt.Print("check ok")
						UpdateSessionTTL(servicename, host)
						os.Exit(0)
					} else if isyes == "noreplication" {
						fmt.Print("replication is not configured")
						UpdateSessionTTL(servicename, host)
						os.Exit(1)
					} else {
						fmt.Print("check replication io_thread fail:", isyes)
						UpdateSessionTTL(servicename, host)
						os.Exit(1)
					}
				} else {
					fmt.Print("check ok")
					UpdateSessionTTL(servicename, host)
					os.Exit(0)
				}
			}
			if MYSQL_OK == 1 && checktime == 0 {
				os.Exit(2)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func TryUpdateCheckTime(user, password, host, port, defaultDb, checktime_string, ping_type, timeout, servicename string) {
	checktime, _ := strconv.Atoi(checktime_string)
	for {
		if checktime == 0 {
			break
		} else {
			checktime--
			MYSQL_OK := CheckMysqlHealth(user, password, host, port, defaultDb, timeout, checktime)
			if MYSQL_OK == 0 {
				if strings.Contains(ping_type, "replication") {
					isyes, err := ShowSlave(user, password, host, port, defaultDb, timeout)
					if err != nil {
						fmt.Print(err)
						os.Exit(2)
					}
					if isyes == "Yes" {
						fmt.Print("check ok")
						UpdateSessionTTL(servicename, host)
						os.Exit(0)
					} else if isyes == "noreplication" {
						fmt.Print("replication is not configured")
						UpdateSessionTTL(servicename, host)
						os.Exit(1)
					} else {
						fmt.Print("check replication io_thread fail:", isyes)
						UpdateSessionTTL(servicename, host)
						os.Exit(1)
					}
				} else {
					fmt.Print("check ok")
					UpdateSessionTTL(servicename, host)
					os.Exit(0)
				}
			}
			if MYSQL_OK == 1 && checktime == 0 {
				os.Exit(2)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func GetConsulConfig(address string) *consulapi.Config {
	config := &consulapi.Config{
		Datacenter: beego.AppConfig.String("datacenter"),
		Token:      beego.AppConfig.String("token"),
		Address:    address,
	}
	return config
}

func GetClient(config *consulapi.Config) (*consulapi.Client, error) {
	client, err := consulapi.NewClient(config)
	if err != nil {
		fmt.Println("Create consul-api client failed!", err)
		return nil, err
	}
	return client, nil
}

func UpdateSessionTTL(servicename, ip string) {
	address := ReadCaAddress()
	config := GetConsulConfig(address)
	client, err := GetClient(config)
	if err != nil {
		return
	}
	kv := client.KV()
	//Get is used to lookup a single key
	kvPair, _, err := kv.Get("cmha/service/"+servicename+"/db/leader", nil)
	if err != nil {
		fmt.Println("Get leader failed!", err)
		return
	}
	if kvPair == nil {
		fmt.Println("Not service leader,Please create!")
		return
	}
	if string(kvPair.Value) == ip {
		if kvPair.Session != "" {
			session := client.Session()
			sessionentry, _, err := session.Renew(kvPair.Session, nil)
			if err != nil {
				fmt.Println("session renew failed:", err)
				return
			}
			fmt.Println("update session success,sessionentry:", sessionentry)
		}
	}
}

func ReadCaAddress() string {
	consul_agent_ip := beego.AppConfig.String("consul_agent_ip")
	consul_agent_port := beego.AppConfig.String("consul_agent_port")
	address := consul_agent_ip + ":" + consul_agent_port
	return address
}
