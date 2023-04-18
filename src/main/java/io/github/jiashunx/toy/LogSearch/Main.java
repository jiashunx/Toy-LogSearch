package io.github.jiashunx.toy.LogSearch;

import com.alibaba.fastjson.JSON;
import io.github.jiashunx.tools.jsch.SSHExecutor;
import io.github.jiashunx.tools.jsch.SSHRequest;
import io.github.jiashunx.tools.jsch.SSHResponse;
import org.apache.commons.io.IOUtils;

import java.io.ByteArrayOutputStream;
import java.io.File;
import java.io.FileInputStream;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Scanner;

/**
 * @author jiashunx
 */
public class Main {

    public static void main(String[] args) throws Throwable {
        Map<String, List<SSHRequest>> configMap = new HashMap<>();
        String configFilePath = System.getProperty("user.dir") + File.separator + "config.json";
        FileInputStream inputStream = new FileInputStream(configFilePath);
        ByteArrayOutputStream outputStream = new ByteArrayOutputStream();
        IOUtils.copy(inputStream, outputStream);
        String jsonContent = new String(outputStream.toByteArray(), StandardCharsets.UTF_8);
        Map<String, Object> jsonMap = JSON.parseObject(jsonContent);
        jsonMap.forEach((key0, value0) -> {
            String envName = key0;
            Map<String, Object> envInfo = (Map<String, Object>) value0;
            envInfo.forEach((key1, value1) -> {
                String serviceName = key1;
                List<Object> serviceList = (List<Object>) value1;
                serviceList.forEach(value2 -> {
                    Map<String, Object> serverInfo = (Map<String, Object>) value2;
                    String ip = (String) serverInfo.get("ip");
                    int port = Integer.parseInt(String.valueOf(serverInfo.get("port")));
                    String username = (String) serverInfo.get("username");
                    String password = (String) serverInfo.get("password");
                    List<String> logPathList = (List<String>) serverInfo.get("log_path");
                    List<String> commands = new ArrayList<>(logPathList.size());
                    commands.addAll(logPathList);
                    System.out.println(String.format("%s -> %s -> %s:%d -> %s/%s -> %s", envName, serviceName, ip, port, username, password, logPathList.toString()));
                    String configId = envName + "___________" + serviceName;
                    configMap.computeIfAbsent(configId, k -> new ArrayList<>()).add(
                            new SSHRequest(ip, port, username, password, commands.toArray(new String[0]))
                    );
                });
            });
        });

        /*
        System.out.println(System.getProperty("user.dir"));
        SSHRequest sshRequest = new SSHRequest("192.168.43.36", 22, "jiashunx", "1234.abcd", "ls");
        SSHResponse sshResponse = SSHExecutor.execCommand(sshRequest);
        System.out.println(sshResponse.getErrorContent());
        System.out.println(sshResponse.getOutputContent());
        */
        Scanner scanner = new Scanner(System.in);
        String inputLine = null;
        while (inputLine == null || inputLine.trim().isEmpty()) {
            try {
                System.out.println("");
                System.out.println("1.命令执行参数格式：[环境名] [服务名] bash [待执行命令(可有空格)]");
                System.out.println("  命令执行参数样例: t3 print bash cat /log/print.log | grep 哈哈哈哈哈哈哈");
                System.out.println("2.日志查询参数格式：[环境名] [服务名] grep [待执行查询条件(自定义grep命令)]");
                System.out.println("  日志查询参数样例: sit2 newcore grep 就将计就计");
                System.out.println("3.日志查询参数格式：[环境名] [服务名] [查询条件1] [查询条件2] [查询条件3]");
                System.out.println("  日志查询参数样例: sit2 newcore 流水号");
                System.out.print("请输入==> ");
//                inputLine = new String(scanner.nextLine().getBytes("GBK"), "UTF-8");
                inputLine = scanner.nextLine();
                System.out.println("录入参数: " + inputLine);
                String[] queryArgs = inputLine.split(" ");
                int argsNum = queryArgs.length;
                if (argsNum <= 2 || argsNum == 3 && ("bash".equals(queryArgs[2]) || "grep".equals(queryArgs[2]))) {
                    System.err.println("查询条件有误！");
                    inputLine = null;
                    continue;
                }
                String envName = queryArgs[0];
                String serviceName = queryArgs[1];
                String configId = envName + "___________" + serviceName;
                if ("bash".equals(queryArgs[2])) {// 第三个参数是bash, 则执行指定命令
                    StringBuilder commandBuilder = new StringBuilder();
                    for (int i = 3; i < argsNum; i++) {
                        commandBuilder.append(queryArgs[i]).append(" ");
                    }
                    configMap.computeIfAbsent(configId, k -> new ArrayList<>()).forEach(r -> {
                        SSHRequest sshRequest = new SSHRequest(r.getRemoteHost(), r.getSshPort(), r.getUsername(), r.getPassword(), commandBuilder.toString());
                        SSHResponse sshResponse = SSHExecutor.execCommand(sshRequest);
                        System.out.println("START ===================================================");
                        System.out.println(String.format("%s -> %s", sshResponse.getRemoteHost(), sshResponse.getCommand()));
                        System.out.println("命令执行是否成功? " + sshResponse.isSuccess());
                        System.out.println("命令执行返回结果: " + sshResponse.getOutputContent());
                        System.out.println("命令执行返回异常: " + sshResponse.getErrorContent());
                        System.out.println("FIN =====================================================");
                    });
                } else if ("grep".equals(queryArgs[2])) {
                    StringBuilder commandSuffix = new StringBuilder(" | grep ");
                    for (int i = 3; i < argsNum; i++) {
                        commandSuffix.append(queryArgs[i]).append(" ");
                    }
                    configMap.computeIfAbsent(configId, k -> new ArrayList<>()).forEach(r -> {
                        doLogQuery(r, commandSuffix.toString());
                    });
                } else {
                    StringBuilder commandSuffix = new StringBuilder();
                    for (int i = 2; i < argsNum; i++) {
                        commandSuffix.append(" | grep ").append(queryArgs[i]);
                    }
                    configMap.computeIfAbsent(configId, k -> new ArrayList<>()).forEach(r -> {
                        doLogQuery(r, commandSuffix.toString());
                    });
                }
            } catch (Throwable throwable) {
                throwable.printStackTrace();
            } finally {
                inputLine = null;
            }
        }
    }

    private static void doLogQuery(SSHRequest r, String commandSuffix) {
        String[] commands = r.getCommandArr();
        String[] newCommands = new String[commands.length];
        for (int i = 0; i < commands.length; i++) {
            newCommands[i] = "cat " + commands[i] + commandSuffix;
        }
        SSHRequest sshRequest = new SSHRequest(r.getRemoteHost(), r.getUsername(), r.getPassword(), newCommands);
        List<SSHResponse> sshResponseList = SSHExecutor.execMultiCommand(sshRequest);
        for (SSHResponse sshResponse: sshResponseList) {
            System.out.println("START ===================================================");
            System.out.println(String.format("%s -> %s", sshResponse.getRemoteHost(), sshResponse.getCommand()));
            System.out.println("日志查询是否成功? " + sshResponse.isSuccess());
            System.out.println("日志查询返回结果: " + sshResponse.getOutputContent());
            System.out.println("日志查询返回异常: " + sshResponse.getErrorContent());
            System.out.println("FIN =====================================================");
        }
    }

}
