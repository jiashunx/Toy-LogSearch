package main

import (
    cfg "Toy-LogSearch/config"
    "Toy-LogSearch/ssh"
    "encoding/json"
    "fmt"
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
        fmt.Println(fmt.Sprintf("SSH Test, Command: %s", response.Command))
        fmt.Println(fmt.Sprintf("SSH Test, Success: %v", response.Success))
        if response.Success {
            fmt.Println(fmt.Sprintf("SSH Test, Output: \n%s", response.OutputContent))
        } else {
            fmt.Println(fmt.Sprintf("SSH Test, Error: \n%s", response.ErrorContent))
        }
    }
    config, err := cfg.LoadConfig()
    if err != nil {
        fmt.Println(fmt.Sprintf("load config failed, error: %v", err))
        return
    }
    bs, err := json.Marshal(config)
    fmt.Println("Config: ", string(bs))
}
