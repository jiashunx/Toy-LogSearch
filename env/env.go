package env

import (
    "Toy-LogSearch/log"
    "fmt"
    "os"
)

var address string
var cfgSrvPath string

func init() {
    args := os.Args
    address0 := ":38888"
    if len(args) >= 2 {
        address0 = args[1]
    }
    address = address0
    log.Info(fmt.Sprintf("运行环境：配置服务监听端口：[%s]", address))
    cfgSrvPath0 := ""
    if (len(args) >= 3) {
        cfgSrvPath0 = args[2]
    }
    cfgSrvPath = cfgSrvPath0
    log.Info(fmt.Sprintf("运行环境：远程配置服务地址：[%s]", cfgSrvPath))
}

func GetServerAddress() string {
    return address
}

func GetCfgSrvPath() string {
    return cfgSrvPath
}
