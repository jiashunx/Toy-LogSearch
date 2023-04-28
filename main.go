package main

import (
    cfg "Toy-LogSearch/config"
    log "Toy-LogSearch/log"
    "Toy-LogSearch/ssh"
    "encoding/json"
    "fmt"
    "go.uber.org/zap"
)

func main() {
    request := &ssh.SSHRequest{
        RemoteHost: "192.168.183.1",
        Port: 22,
        Username: "jiashunx",
        Password: "1234.abcd",
        Commands: []string { "/bin/ps", "cat /etc/profileinfo", "/bin/ps -ef" },
    }
    resps := ssh.ExecuteCommand(request)
    for _, response := range resps {
        log.Info(fmt.Sprintf("SSH Test, Command: %s", response.Command))
        log.Info(fmt.Sprintf("SSH Test, Success: %v", response.Success))
        if response.Success {
            log.Info(fmt.Sprintf("SSH Test, Output: \n%s", response.OutputContent))
        } else {
            log.Info(fmt.Sprintf("SSH Test, Error: \n%s", response.ErrorContent))
        }
    }
    config, err := cfg.LoadConfig()
    if err != nil {
        log.Info(fmt.Sprintf("load config failed, error: %v", err))
        return
    }
    bs, err := json.Marshal(config)
    log.Info(fmt.Sprintf("Config: %s", string(bs)), zap.String("k", "value"))
}
