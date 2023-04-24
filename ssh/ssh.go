package ssh

import (
    "bytes"
    "fmt"
    "golang.org/x/crypto/ssh"
)

type SSHRequest struct {
    RemoteHost      string
    Port            uint8
    Username        string
    Password        string
    Commands        []string
}

type SSHResponse struct {
    RemoteHost      string
    Command         string
    Success         bool
    CostMillis      int64
    OutputContent   string
    ErrorContent    string
}

func ExecuteCommands(requests []*SSHRequest) []*SSHResponse {
    resps := make([]*SSHResponse, 0)
    if requests != nil {
        for _, request := range requests {
            resps = append(resps, ExecuteCommand(request)...)
        }
    }
    return resps
}

func ExecuteCommand(request *SSHRequest) []*SSHResponse {
    resps := make([]*SSHResponse, 0)
    if request != nil && request.Commands != nil {
        for _, command := range request.Commands {
            resps = append(resps, &SSHResponse{
                RemoteHost: request.RemoteHost,
                Command: command,
            })
        }
    }
    for _, response := range resps {
        doExecute(request, response)
    }
    return resps
}

func doExecute(request *SSHRequest, response *SSHResponse) {
    config := &ssh.ClientConfig{
        User: request.Username,
        Auth: []ssh.AuthMethod{ssh.Password(request.Password)},
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }
    client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", request.RemoteHost, request.Port), config)
    if err != nil {
        response.ErrorContent = fmt.Sprintf("Create SSH Client Failed, Error: %v", err)
        return
    }
    defer client.Close()
    session, err := client.NewSession()
    if err != nil {
        response.ErrorContent = fmt.Sprintf("Create SSH Session Failed, Error: %v", err)
        return
    }
    defer session.Close()
    var stdout bytes.Buffer
    var stderr bytes.Buffer
    session.Stdout = &stdout
    session.Stderr = &stderr
    // 同一会话Run, Start, Shell, Output, or CombinedOutput 等方法仅可执行一次
    errInfo := ""
    if err = session.Run(response.Command); err != nil {
        errInfo = fmt.Sprintf("%v", err)
    }
    response.OutputContent = stdout.String()
    response.ErrorContent = stderr.String() + errInfo
    response.Success = len(response.ErrorContent) == 0
}

