package main

import (
    "bytes"
    "fmt"
    "golang.org/x/crypto/ssh"
)

func main() {
    commands := []string { "/bin/ps", "cat /etc/profileinfo" }
    defer func() {
        if v := recover(); v != nil {
            fmt.Println(fmt.Sprintf("SSH Test Occur Panic Error: %v", v))
        }
    }()
    config := &ssh.ClientConfig{
        User: "jiashunx",
        Auth: []ssh.AuthMethod{ssh.Password("1234.abcd")},
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }
    client, err := ssh.Dial("tcp", "192.168.183.1:22", config)
    if err != nil {
        fmt.Println(fmt.Sprintf("Create SSH Client Failed, Error: %v", err))
        return
    }
    defer func() {
        err := client.Close()
        if err != nil {
            fmt.Println(fmt.Sprintf("Close SSH Client Failed, Error: %v", err))
        }
    }()
    session, err := client.NewSession()
    if err != nil {
        fmt.Println(fmt.Sprintf("Create SSH Session Failed, Error: %v", err))
        return
    }
    defer func() {
        err := session.Close()
        if err != nil {
            fmt.Println(fmt.Sprintf("Close SSH Session Failed, Error: %v", err))
        }
    }()
    for _, command := range commands {
        var stdout bytes.Buffer
        var stderr bytes.Buffer
        session.Stdout = &stdout
        session.Stderr = &stderr
        // TODO 调整为同一会话支持执行多个命令
        // 同一会话Run, Start, Shell, Output, or CombinedOutput 等方法仅可执行一次
        if err = session.Run(command); err != nil {
            fmt.Println(fmt.Sprintf("SSH Test, Execute Command [%s] Failed, Error: %v", command, err))
        }
        fmt.Println(fmt.Sprintf("SSH Test, Execute Command [%s] StdOut: %v", command, stdout.String()))
        fmt.Println(fmt.Sprintf("SSH Test, Execute Command [%s] StdErr: %v", command, stderr.String()))
    }

}
