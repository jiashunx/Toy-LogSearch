package model

import (
    "Toy-LogSearch/env"
    "encoding/json"
    "errors"
    "fmt"
    "go.uber.org/zap"
    "net/http"
    "os"
    "strings"
    "time"
    "Toy-LogSearch/log"
)

// 配置模型
type Config struct {
    Servers     []Server  `json:"servers"`
    Services    []Service `json:"services"`
}

type Server struct {
    RemoteHost  string          `json:"ip"`
    Port        uint8           `json:"port"`
    Username    string          `json:"username"`
    Password    string          `json:"password"`
}

type Service struct {
    Env         string          `json:"env"`
    Service     string          `json:"service"`
    Configs     []ServiceConfig `json:"configs"`
}

type ServiceConfig struct {
    RemoteHost  string          `json:"ip"`
    LogPaths    []string        `json:"logPaths"`
}

func (c *Config) GetServerByIp(ip string) *Server {
    for _, server := range c.Servers {
        if ip != "" && ip == server.RemoteHost {
            return &server
        }
    }
    return nil
}

func (c *Config) GetEnvServiceConfigs(e, s string) []ServiceConfig {
    sc := make([]ServiceConfig, 0)
    for _, service := range c.Services {
        if e != "" && e == service.Env && s != "" && s == service.Service {
            sc = append(sc, service.Configs...)
        }
    }
    return sc
}

func (c *Config) PrintConfigInfo() {
    for _, server := range c.Servers {
        _, _ = fmt.Printf("server: %s:%d@%s/%s\n", server.RemoteHost, server.Port, server.Username, server.Password)
    }
    for _, service := range c.Services {
        for _, sc := range service.Configs {
            _, _ = fmt.Printf("env: %s, service: %s, ip: %s, logPaths: %v\n", service.Env, service.Service, sc.RemoteHost, sc.LogPaths)
        }
    }
}

func (c *Config) verifyBeanInfo() error {
    if c == nil {
        return errors.New("config pointer can't be null")
    }
    if c.Servers == nil {
        return errors.New("servers can't be null")
    }
    for _, server := range c.Servers {
        //if server == nil {
        //    return errors.New("server can't be null")
        //}
        if server.RemoteHost == "" {
            return errors.New("server ip can't be null")
        }
    }
    if c.Services == nil {
        return errors.New("services can't be null")
    }
    for _, service := range c.Services {
        //if service == nil {
        //    return errors.New("service can't be null")
        //}
        if service.Env == "" {
            return errors.New("service env can't be null")
        }
        if service.Service == "" {
            return errors.New("service env service can't be null")
        }
        if service.Configs == nil {
            return errors.New("service configs can't be null")
        }
        for _, scfg := range service.Configs {
            //if scfg == nil {
            //    return errors.New("service config can't be null")
            //}
            if scfg.RemoteHost == "" {
                return errors.New("service config ip can't be null")
            }
            if scfg.LogPaths == nil {
                return errors.New("service config logPaths can't be null")
            }
        }
    }
    return nil
}

func (c *Config) toJson() ([]byte, error) {
    return json.MarshalIndent(c, "", "  ")
}

// 加载配置
func LoadConfig() (*Config, error) {
    if env.GetCfgSrvPath() != "" {
        config, err := resolveFromConfigServer(env.GetCfgSrvPath())
        if err != nil {
            return nil, err
        }
        _, _ = store(config)
        return config, nil
    }
    return resolveFromFile()
}

// 配置信息存储到本地
func store(config *Config) (string, error) {
    _ = os.MkdirAll("config", os.ModePerm)
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    timestamp = strings.ReplaceAll(timestamp, "-", "")
    timestamp = strings.ReplaceAll(timestamp, ":", "")
    timestamp = strings.ReplaceAll(timestamp, " ", "")
    _ = os.MkdirAll("./config", 0777)
    fileName := fmt.Sprintf("./config/config%s.json", timestamp)
    file, err := os.Create(fileName)
    if err != nil {
        log.Error(fmt.Sprintf("create file [%s] failed, error: %v", fileName, err))
        return "", err
    }
    bytes, err := config.toJson()
    if err != nil {
        log.Error(fmt.Sprintf("config marshall failed, error: %v", err))
        return "", err
    }
    _, err = file.Write(bytes)
    if err != nil {
        log.Error(fmt.Sprintf("write config to file [%s] failed, error: %v", fileName, err))
        return "", err
    }
    return fileName, nil
}

// 从远程配置服务器同步配置
func resolveFromConfigServer(cfgSrvPath string) (*Config, error) {
    fmt.Println(fmt.Sprintf("从配置服务[%s]获取配置信息", cfgSrvPath));
    url := cfgSrvPath + "/config.json"
    response, err := http.Get(cfgSrvPath + "/config.json")
    if err != nil {
        log.Error(fmt.Sprintf("从配置服务[%s]获取配置信息失败", url), zap.Error(err))
        return nil, err
    }
    log.Info(fmt.Sprintf("从配置服务[%s]获取配置信息，响应码：%d", url, response.StatusCode))
    bs := make([]byte, response.ContentLength)
    n, err := response.Body.Read(bs)
    if err != nil {
        log.Error(fmt.Sprintf("从配置服务[%s]获取配置信息，读取响应失败", url), zap.Error(err))
    }
    content := string(bs[0:n])
    log.Info(fmt.Sprintf("从配置服务[%s]获取配置信息：%s", url, content))
    return resolveFromContent(bs[0:n])
}

// 从本地配置文件加载配置
func resolveFromFile() (*Config, error) {
    cfgPath := "config.json"
    fmt.Println(fmt.Sprintf("从本地[%s]获取配置信息", cfgPath));
    bs, err := os.ReadFile(cfgPath)
    if err != nil {
        log.Error("从config.json读取配置信息失败", zap.Error(err))
        return nil, err
    }
    return resolveFromContent(bs)
}

func resolveFromContent(bs []byte) (*Config, error) {
    config := &Config{}
    if err := json.Unmarshal(bs, config); err != nil {
        log.Error("从json反序列化配置对象异常", zap.Error(err))
        return nil, err
    }
    if err := config.verifyBeanInfo(); err != nil {
        log.Error("从json反序列化配置对象属性校验失败", zap.Error(err))
        return nil, err
    }
    return config, nil
}
