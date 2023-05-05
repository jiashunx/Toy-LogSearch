package main

import (
    _ "Toy-LogSearch/env"
    "Toy-LogSearch/log"
    "Toy-LogSearch/model"
    "Toy-LogSearch/server"
    "Toy-LogSearch/ssh"
    "Toy-LogSearch/utils"
    "bufio"
    "fmt"
    "go.uber.org/zap"
    "os"
    "strings"
    "time"
)

// go run main.go :38889 http://127.0.0.1:38888
func main() {
    s := &server.EchoServer{ConfigRef: nil}
    go func() {
        s.StartServer()
    }()
    time.Sleep(time.Second * 3)
    configRef, err := model.LoadConfig()
    s.ConfigRef = configRef
    if err != nil {
        fmt.Println("配置信息解析异常")
    }
    if configRef != nil {
        configRef.PrintConfigInfo()
    }
    utils.PrintHelpInfo()
    reader := bufio.NewReader(os.Stdin)
    text := ""
    for {
        fmt.Print("请输入==> ")
        bytes, _, err := reader.ReadLine()
        if err != nil {
            log.Error("Read Content From Console Failed", zap.Error(err))
            fmt.Println("录入参数处理异常，无法匹配命令，请重新输入！")
            utils.PrintHelpInfo()
            continue
        }

        text = string(bytes)
        command := utils.GetCommand(text)
        if command == nil {
            fmt.Println("录入参数有误，无法匹配命令，请重新输入！")
            utils.PrintHelpInfo()
            continue
        }

        cmdArgs := utils.GetCommandArgs(text)
        if len(cmdArgs) == 0 {
            fmt.Println("录入参数解析失败，请重新输入！")
            utils.PrintHelpInfo()
            continue
        }

        if command == utils.CHelp {
            utils.PrintHelpInfo()
            continue
        }

        if command == utils.CExit || command == utils.CQuit {
            log.Info("结束main线程")
            os.Exit(0)
        }

        if command == utils.CReload {
            tmpConfig, err := model.LoadConfig()
            if err != nil {
                fmt.Println("重新加载配置文件，获取配置信息为null，不更新配置信息！")
                if configRef != nil {
                    configRef.PrintConfigInfo()
                }
                utils.PrintHelpInfo()
                continue
            }
            configRef = tmpConfig
            // 更新配置信息
            s.ConfigRef = configRef
            configRef.PrintConfigInfo()
            utils.PrintHelpInfo()
            continue
        }

        if configRef == nil {
            fmt.Println("当前配置对象为空，请同步配置后执行相应命令！")
            utils.PrintHelpInfo()
            continue
        }

        env := cmdArgs[0]
        service := cmdArgs[1]
        scs := configRef.GetEnvServiceConfigs(env, service)
        if (len(scs) == 0) {
            fmt.Println(fmt.Sprintf("未找到env[%s], service[%s] 对应服务配置，请重新输入！", env, service))
            utils.PrintHelpInfo()
            continue
        }
        serverMap := make(map[string]*model.Server)
        for _, sc := range scs {
            server := configRef.GetServerByIp(sc.RemoteHost)
            if server != nil {
                serverMap[sc.RemoteHost] = server
            }
        }
        if len(serverMap) == 0 {
            fmt.Println(fmt.Sprintf("未找到env[%s], service[%s] 对应服务器配置，请重新输入！", env, service))
            utils.PrintHelpInfo()
            continue
        }

        cmd := ""
        // 根据配置文件中配置的日志路径进行检索
        if command == utils.CT1 {
            var cmdBuilder strings.Builder
            cmdBuilder.WriteString(" | grep ")
            for i, l := 2, len(cmdArgs); i < l; i++ {
                cmdBuilder.WriteString(cmdArgs[i])
                cmdBuilder.WriteString(" ")
                if i < l - 1 {
                    cmdBuilder.WriteString("| grep ")
                }
            }
            cmd = cmdBuilder.String()
        }
        // 根据配置文件中配置的日志路径进行自定义grep命令检索
        if command == utils.CT2 {
            cmd = " | grep " + cmdArgs[3]
        }
        // 执行自定义脚本命令
        if command == utils.CT3 {
            cmd = cmdArgs[3]
        }

        for _, sc := range scs {
            if server, ok := serverMap[sc.RemoteHost]; ok {
                // 根据配置文件中配置的日志路径进行自定义grep命令检索
                if command == utils.CT1 || command == utils.CT2 {
                    executeLogQuery(server, &sc, cmd)
                }
                // 执行自定义脚本命令
                if command == utils.CT3 {
                    executeBashCommand(server, &sc, cmd)
                }
            }
        }
    }
}

func executeLogQuery(server *model.Server, configRef *model.ServiceConfig, commandSuffix string) {
    logPaths := configRef.LogPaths
    commands := make([]string, len(logPaths))
    for i, logPath := range logPaths {
        commands[i] = "cat " + logPath + commandSuffix
    }
    sshRequest := &ssh.SSHRequest{
        RemoteHost: server.RemoteHost,
        Port: server.Port,
        Username: server.Username,
        Password: server.Password,
        Commands: commands,
    }
    sshResponses := ssh.ExecuteCommand(sshRequest)
    for _, sshResponse := range sshResponses {
        fmt.Println("BGN ===================================================")
        fmt.Println(fmt.Sprintf("%s -> %s", sshResponse.RemoteHost, sshResponse.Command))
        fmt.Println(fmt.Sprintf("日志查询是否成功? %t", sshResponse.Success))
        if sshResponse.Success {
            fmt.Println("日志查询返回结果: ")
            fmt.Println(sshResponse.OutputContent)
        } else {
            fmt.Println("日志查询返回异常: ")
            fmt.Println(sshResponse.ErrorContent)
        }
        fmt.Println("FIN ===================================================")
    }
}

func executeBashCommand(server *model.Server, configRef *model.ServiceConfig, command string) {
    sshRequest := &ssh.SSHRequest{
        RemoteHost: server.RemoteHost,
        Port: server.Port,
        Username: server.Username,
        Password: server.Password,
        Commands: []string{ command },
    }
    sshResponses := ssh.ExecuteCommand(sshRequest)
    for _, sshResponse := range sshResponses {
        fmt.Println("BGN ===================================================")
        fmt.Println(fmt.Sprintf("%s -> %s", sshResponse.RemoteHost, sshResponse.Command))
        fmt.Println(fmt.Sprintf("命令执行是否成功? %t", sshResponse.Success))
        if sshResponse.Success {
            fmt.Println("命令执行返回结果: ")
            fmt.Println(sshResponse.OutputContent)
        } else {
            fmt.Println("命令执行返回异常: ")
            fmt.Println(sshResponse.ErrorContent)
        }
        fmt.Println("FIN ===================================================")
    }
}
