package config

import (
    "encoding/json"
    "errors"
    "fmt"
    "os"
    "strings"
    "time"
)

// 配置模型
type Config struct {
    Servers [] struct {
        RemoteHost  string  `json:"ip"`
        Port        uint8   `json:"port"`
        Username    string  `json:"username"`
        Password    string  `json:"password"`
    }                       `json:"servers"`
    Services [] struct {
        Env         string  `json:"env"`
        Service     string  `json:"service"`
        Configs []struct{
            RemoteHost  string      `json:"ip"`
            LogPaths    []string    `json:"logPaths"`
        }                   `json:"configs"`
    }                       `json:"services"`
}

func (c *Config) verify() error {
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
    args := os.Args
    if len(args) >= 3 {
        config, err := resolveFromConfigServer(args[2])
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
    fileName := fmt.Sprintf("config%s.json", timestamp)
    file, err := os.Create(fileName)
    if err != nil {
        fmt.Println(fmt.Sprintf("create file [%s] failed, error: %v", fileName, err))
        return "", err
    }
    bytes, err := config.toJson()
    if err != nil {
        fmt.Println(fmt.Sprintf("config marshall failed, error: %v", err))
        return "", err
    }
    _, err = file.Write(bytes)
    if err != nil {
        fmt.Println(fmt.Sprintf("write config to file [%s] failed, error: %v", fileName, err))
        return "", err
    }
    return fileName, nil
}

// 从远程配置服务器同步配置
func resolveFromConfigServer(configServerPath string) (*Config, error) {
    return nil, nil
}

// 从本地配置文件加载配置
func resolveFromFile() (*Config, error) {
    bs, err := os.ReadFile("config.json")
    if err != nil {
        return nil, err
    }
    config := &Config{}
    if err := json.Unmarshal(bs, config); err != nil {
        return nil, err
    }
    if err := config.verify(); err != nil {
        return nil, err
    }
    return config, nil
}
