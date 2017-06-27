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

func IsPingType(user, password string, cilent *consulapi.Client) {
	host, port, checktime_string, timeout, defaultDb, ping_type := GetConfig()
	servicename := beego.AppConfig.String("servicename")
	if ping_type == "select,replication" || ping_type == "select" {
		TrySelectCheckTime(user, password, host, port, defaultDb, checktime_string, ping_type, timeout, servicename, cilent)
	} else if ping_type == "update,replication" || ping_type == "update" {
		TryUpdateCheckTime(user, password, host, port, defaultDb, checktime_string, ping_type, timeout, servicename, cilent)
	} else {
		fmt.Print("Configuration error")
		os.Exit(2)
	}

}

func TrySelectCheckTime(user, password, host, port, defaultDb, checktime_string, ping_type, timeout, servicename string, client *consulapi.Client) {
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
						UpdateSessionTTL(servicename, host, client)
						os.Exit(0)
					} else if isyes == "noreplication" {
						fmt.Print("replication is not configured")
						UpdateSessionTTL(servicename, host, client)
						os.Exit(1)
					} else {
						fmt.Print("check replication io_thread fail:", isyes)
						UpdateSessionTTL(servicename, host, client)
						os.Exit(1)
					}
				} else {
					fmt.Print("check ok")
					UpdateSessionTTL(servicename, host, client)
					os.Exit(0)
				}
			}
			if MYSQL_OK == 1 && checktime == 0 {
				os.Exit(2)
			}
		}
		time.Sleep(2 * time.Second)
	}
}

func TryUpdateCheckTime(user, password, host, port, defaultDb, checktime_string, ping_type, timeout, servicename string, client *consulapi.Client) {
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
						UpdateSessionTTL(servicename, host, client)
						os.Exit(0)
					} else if isyes == "noreplication" {
						fmt.Print("replication is not configured")
						UpdateSessionTTL(servicename, host, client)
						os.Exit(1)
					} else {
						fmt.Print("check replication io_thread fail:", isyes)
						UpdateSessionTTL(servicename, host, client)
						os.Exit(1)
					}
				} else {
					fmt.Print("check ok")
					UpdateSessionTTL(servicename, host, client)
					os.Exit(0)
				}
			}
			if MYSQL_OK == 1 && checktime == 0 {
				os.Exit(2)
			}
		}
		time.Sleep(2 * time.Second)
	}
}

func UpdateSessionTTL(servicename, ip string, client *consulapi.Client) {
	session := client.Session()
	node, _, err := session.Node(beego.AppConfig.String("hostname"), nil)
	if err != nil {
		fmt.Print(err)
		return
	}
	if node != nil {
		for i := range node {
			for j := 0; j < 3; j++ {
				err := renewSession(node[i].ID, session)
				if err != nil {
					fmt.Print(err)
					continue
				}
				break
			}

			/*_, _, err := session.Renew(node[i].ID, nil)
			if err != nil {
				fmt.Print(err)
				return
			}*/
		}

	}

}

func renewSession(id string, session *consulapi.Session) error {
	_, _, err := session.Renew(id, nil)
	if err != nil {
		fmt.Print(err)
		return err
	}
	return nil
}

func ReadCaAddress() string {
	consul_agent_ip := beego.AppConfig.String("consul_agent_ip")
	consul_agent_port := beego.AppConfig.String("consul_agent_port")
	address := consul_agent_ip + ":" + consul_agent_port
	return address
}

func getIPaddr(ip string) string {
	ip_port := strings.Split(ip, ":")
	for i := range ip_port {
		if i == 0 {
			return ip_port[i]
		}
	}
	return ""
}
