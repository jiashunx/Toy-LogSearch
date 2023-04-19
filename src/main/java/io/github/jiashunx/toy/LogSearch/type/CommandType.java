package io.github.jiashunx.toy.LogSearch.type;

/**
 * @author jiashunx
 */
public enum CommandType {

    // 1. [env] [service] [condition1] [condition2] [condition3]
    // [condition]不可出现空字符串
    T1("^\\S+\\s+\\S+\\s+\\S+(\\s|\\S)*$"),
    // 2. [env] [service] grep [condition1] [condition2] [condition3]
    // [condition]可出现空字符串
    T2("^\\S+\\s+\\S+\\s+(grep){1}\\s+(\\s|\\S)+$"),
    // 3. [env] [service] bash [command1] [command2] [command3]
    // [command]可出现空字符串
    T3("^\\S+\\s+\\S+\\s+(bash){1}\\s+(\\s|\\S)+$"),
    // 重新加载配置信息
    RELOAD("^(reload){1}$"),
    // help信息
    HELP("^(help){1}$"),
    // 退出进程
    QUIT("^(quit){1}$"),
    // 退出进程
    EXIT("^(exit){1}$");

    public String regex;

    CommandType(String regex) {
        this.regex = regex;
    }

    @Override
    public String toString() {
        return String.format("CommandType=%s{regex=[%s]}", name(), regex);
    }
}
