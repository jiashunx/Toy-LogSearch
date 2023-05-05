package utils

import (
    "Toy-LogSearch/log"
    "fmt"
    "go.uber.org/zap"
    "regexp"
    "strings"
)

type Command struct {
    Name    string
    Regex   string
}

func (c *Command) String() string {
    return fmt.Sprintf("CommandType=%s{regex=[%s]}", c.Name, c.Regex)
}

var (
    // 1. [env] [service] [condition1] [condition2] [condition3]
    // [condition]不可出现空字符串
    CT1       = &Command{Name: "T1", Regex: "^\\S+\\s+\\S+\\s+\\S+(\\s|\\S)*$"}
    // 2. [env] [service] grep [condition1] [condition2] [condition3]
    // [condition]可出现空字符串
    CT2       = &Command{Name: "T2", Regex:"^\\S+\\s+\\S+\\s+(grep){1}\\s+(\\s|\\S)+$"}
    // 3. [env] [service] bash [command1] [command2] [command3]
    // [command]可出现空字符串
    CT3       = &Command{Name: "T3", Regex: "^\\S+\\s+\\S+\\s+(bash){1}\\s+(\\s|\\S)+$"}
    // 重新加载配置信息
    CReload   = &Command{Name: "RELOAD", Regex: "^(reload){1}$"}
    // help信息
    CHelp     = &Command{Name: "HELP", Regex: "^(help){1}$"}
    // 退出进程
    CQuit     = &Command{Name: "QUIT", Regex: "^(quit){1}$"}
    // 退出进程
    CExit     = &Command{Name: "EXIT", Regex: "^(exit){1}$"}
)

var Commands = []*Command{ CT1, CT2, CT3, CReload, CHelp, CQuit, CExit }

func PrintHelpInfo() {
    fmt.Println("1.日志查询命令：[env] [service] [condition1] [condition2]...[conditionN]（condition不可出现空格）");
    fmt.Println("  日志查询样例: sit2 newcore 2023-03-19 17\\:17\\:17 202304181703222101000066");
    fmt.Println("2.日志查询命令：[env] [service] grep [condition1] [condition2]...[conditionN]（condition可出现空格）");
    fmt.Println("  日志查询样例: sit2 newcore grep 2023-03-19 | grep 17\\:17\\:17 | grep 202304181703222101000066");
    fmt.Println("3.自定义脚本命令：[env] [service] bash [command1] [command2]...[commandN]（command可出现空格）");
    fmt.Println("  自定义脚本样例: sit2 newcore bash cat /log/print.log | grep 2023-03-19 | grep 17\\:17\\:17 | grep 202304181703222101000066");
    fmt.Println("4.其他命令：help（输出帮助信息）、reload（更新配置）、quit（退出进程）、exit（退出进程）");
}

func GetCommand(text string) *Command {
    content := strings.TrimSpace(text)
    var command *Command
    for _, c := range Commands {
        b, err := regexp.MatchString(c.Regex, content)
        if err != nil {
            log.Error(fmt.Sprintf("正则匹配处理失败, 正则表达式: [%s], 原字符串: [%s]", c.Regex, content), zap.Error(err))
            continue
        }
        // 循环匹配（以最后一个匹配到的为准）
        if b {
            command = c
        }
    }
    if command == nil {
        log.Info(fmt.Sprintf("命令[%s]对应命令类型: nil", content))
    } else {
        log.Info(fmt.Sprintf("命令[%s]对应命令类型: %s", content, command.String()))
    }
    return command
}

func GetCommandArgs(text string) []string {
    args := make([]string, 0)
    var word string
    content := strings.TrimSpace(text)
    command := GetCommand(text)
    if command != nil {
        switch command {
        case CT1:
            content1 := content
            word = GetFirstWord(content1)
            for ; word != ""; {
                args = append(args, word)
                content1 = content1[(strings.Index(content1, word) + len(word)):]
                word = GetFirstWord(content1)
            }
        case CT2, CT3:
            content23 := content
            loopCount := 0
            word = GetFirstWord(content23)
            for ; word != "" && loopCount < 3; {
                args = append(args, word)
                content23 = content23[(strings.Index(content23, word) + len(word)):]
                loopCount++
                if loopCount == 3 {
                    args = append(args, strings.TrimSpace(content23))
                    break
                }
                word = GetFirstWord(content23)
            }
        case CReload, CHelp, CQuit, CExit:
            args = append(args, GetFirstWord(content))
        default:
            // do nothing.
        }
        log.Info(fmt.Sprintf("命令[%s]解析拆分参数: %v", content, args))
    }
    return args
}

func GetFirstWord(text string) string {
    content := strings.TrimSpace(text)
    var builder strings.Builder
    for _, r := range []rune(content) {
        if r == ' ' {
            break
        }
        builder.WriteRune(r)
    }
    return builder.String()
}
