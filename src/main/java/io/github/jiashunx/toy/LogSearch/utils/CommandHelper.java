package io.github.jiashunx.toy.LogSearch.utils;

import com.alibaba.fastjson.JSON;
import io.github.jiashunx.masker.rest.framework.util.StringUtils;
import io.github.jiashunx.toy.LogSearch.type.CommandType;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.List;

/**
 * @author jiashunx
 */
public class CommandHelper {

    private static final Logger logger = LoggerFactory.getLogger(CommandHelper.class);

    public static void printHelpInfo() {
        System.out.println("1.日志查询命令：[env] [service] [condition1] [condition2]...[conditionN]（condition不可出现空格）");
        System.out.println("  日志查询样例: sit2 newcore 2023-03-19 17\\:17\\:17 202304181703222101000066");
        System.out.println("2.日志查询命令：[env] [service] grep [condition1] [condition2]...[conditionN]（condition可出现空格）");
        System.out.println("  日志查询样例: sit2 newcore grep 2023-03-19 | grep 17\\:17\\:17 | grep 202304181703222101000066");
        System.out.println("3.自定义脚本命令：[env] [service] bash [command1] [command2]...[commandN]（command可出现空格）");
        System.out.println("  自定义脚本样例: sit2 newcore bash cat /log/print.log | grep 2023-03-19 | grep 17\\:17\\:17 | grep 202304181703222101000066");
        System.out.println("4.其他命令：help（输出帮助信息）、reload（更新配置）、quit（退出进程）、exit（退出进程）");
    }

    /**
     * 校验命令合法性
     * @param text 命令行内容
     * @return CommandType
     */
    public static CommandType getCommandType(String text) {
        String content = String.valueOf(text).trim();
        CommandType[] types = CommandType.values();
        CommandType matchedType = null;
        for (CommandType type : types) {
            if (content.matches(type.regex)) { // 循环匹配（以最后一个匹配到的为准）
                matchedType = type;
            }
        }
        if (logger.isInfoEnabled()) {
            logger.info("命令[{}]对应命令类型: {}", content, matchedType == null ? "" : matchedType.toString());
        }
        return matchedType;
    }

    public static String[] getCommandArgs(String text) {
        String[] args = null;
        List<String> words = new ArrayList<>();
        String word = "";
        String content = String.valueOf(text).trim();
        CommandType type = getCommandType(content);
        if (type != null) {
            switch (type) {
                case T1:
                    String content1 = content;
                    word = getFirstWord(content1);
                    while (StringUtils.isNotEmpty(word)) {
                        words.add(word);
                        content1 = content1.substring(content1.indexOf(word) + word.length());
                        word = getFirstWord(content1);
                    }
                    args = words.toArray(new String[0]);
                    break;
                case T2:
                case T3:
                    String content23 = content;
                    int loopCount = 0;
                    word = getFirstWord(content23);
                    while (StringUtils.isNotEmpty(word) && loopCount < 3) {
                        words.add(word);
                        loopCount++;
                        content23 = content23.substring(content23.indexOf(word) + word.length());
                        if (loopCount == 3) {
                            words.add(content23.trim());
                            break;
                        }
                        word = getFirstWord(content23);
                    }
                    args = words.toArray(new String[0]);
                    break;
                case RELOAD:
                case HELP:
                case QUIT:
                case EXIT:
                    args = new String[] { getFirstWord(content) };
                    break;
                default:
            }
        }
        if (logger.isInfoEnabled()) {
            logger.info("命令[{}]解析拆分参数: {}", content, args == null ? "null" : JSON.toJSONString(args));
        }
        return args;
    }

    private static String getFirstWord(String text) {
        String content = String.valueOf(text).trim();
        StringBuilder builder = new StringBuilder();
        for (int i = 0; i < content.length(); i++) {
            if (content.charAt(i) == ' ') {
                break;
            }
            builder.append(content.charAt(i));
        }
        return builder.toString();
    }

}
