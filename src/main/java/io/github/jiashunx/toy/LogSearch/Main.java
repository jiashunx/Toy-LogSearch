package io.github.jiashunx.toy.LogSearch;

import com.alibaba.fastjson.JSON;
import io.github.jiashunx.masker.rest.framework.MRestServer;
import io.github.jiashunx.masker.rest.framework.util.FileUtils;
import io.github.jiashunx.masker.rest.framework.util.IOUtils;
import io.github.jiashunx.masker.rest.framework.util.StringUtils;
import io.github.jiashunx.tools.jsch.SSHExecutor;
import io.github.jiashunx.tools.jsch.SSHRequest;
import io.github.jiashunx.tools.jsch.SSHResponse;
import io.github.jiashunx.toy.LogSearch.model.AllConfig;
import io.github.jiashunx.toy.LogSearch.model.Server;
import io.github.jiashunx.toy.LogSearch.model.ServiceConfig;
import io.github.jiashunx.toy.LogSearch.type.CommandType;
import io.github.jiashunx.toy.LogSearch.utils.CommandHelper;
import io.netty.handler.codec.http.HttpResponseStatus;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Scanner;
import java.util.concurrent.atomic.AtomicReference;

/**
 * -Dserver.port=38889
 * -Dconfig.server=http://127.0.0.1:38888
 * @author jiashunx
 */
public class Main {

    private static final Logger logger = LoggerFactory.getLogger(Main.class);

    public static void main(String[] args) throws Throwable {
        AtomicReference<AllConfig> allConfigRef = new AtomicReference<>(loadConfig(args));
        if (allConfigRef.get() == null) {
            System.out.println("配置信息解析异常");
        }
        new Thread(() -> {
            String port = System.getProperty("server.port");
            if (StringUtils.isEmpty(port)) {
                port = "38888";
            }
            new MRestServer(Integer.parseInt(port), "config-server")
                .context("/")
                .filter("/*", (request, response, filterChain) -> {
                    String requestUrl = request.getUrl();
                    if (requestUrl.equals("/config.json")) {
                        filterChain.doFilter(request, response);
                        return;
                    }
                    response.write(HttpResponseStatus.NOT_FOUND);
                })
                .get("/config.json", (request, response) -> {
                    if (allConfigRef.get() == null) {
                        response.write(HttpResponseStatus.NOT_FOUND);
                        return;
                    }
                    response.write(allConfigRef.get());
                })
                .getRestServer()
                .start();
        }, "config-server").start();
        Scanner scanner = new Scanner(System.in);
        String inputLine = null;
        if (allConfigRef.get() != null) {
            allConfigRef.get().printConfigInfo();
        }
        CommandHelper.printHelpInfo();
        while (true) {
            try {
                System.out.print("请输入==> ");
                // inputLine = new String(scanner.nextLine().getBytes("GBK"), "UTF-8");
                inputLine = scanner.nextLine();

                CommandType commandType = CommandHelper.getCommandType(inputLine);
                if (commandType == null) {
                    System.out.println("录入参数有误，无法匹配命令，请重新输入！");
                    CommandHelper.printHelpInfo();
                    continue;
                }

                String[] commandArgs = CommandHelper.getCommandArgs(inputLine);
                if (commandArgs == null) {
                    System.out.println("录入参数解析失败，请重新输入！");
                    CommandHelper.printHelpInfo();
                    continue;
                }

                if (commandType == CommandType.HELP) {
                    CommandHelper.printHelpInfo();
                    continue;
                }

                if (commandType == CommandType.EXIT || commandType == CommandType.QUIT) {
                    if (logger.isWarnEnabled()) {
                        logger.warn("结束main线程.");
                    }
                    System.exit(1);
                    break;
                }

                if (commandType == CommandType.RELOAD) {
                    AllConfig allConfig = loadConfig(args);
                    if (allConfig == null) {
                        System.out.println("重新加载配置文件，获取配置信息为null，不更新配置信息！");
                        if (allConfigRef.get() != null) {
                            allConfigRef.get().printConfigInfo();
                        }
                        CommandHelper.printHelpInfo();
                        continue;
                    }
                    allConfigRef.set(allConfig);
                    allConfigRef.get().printConfigInfo();
                    CommandHelper.printHelpInfo();
                    continue;
                }

                if (allConfigRef.get() == null) {
                    System.out.println("当前配置对象为空，请同步配置后执行相应命令！");
                    CommandHelper.printHelpInfo();
                    continue;
                }

                String env = commandArgs[0];
                String service = commandArgs[1];
                List<ServiceConfig> configs = allConfigRef.get().getEnvServiceConfigs(env, service);
                if (configs.isEmpty()) {
                    System.out.println(String.format("未找到env[%s], service[%s] 对应服务配置，请重新输入！", env, service));
                    CommandHelper.printHelpInfo();
                    continue;
                }
                Map<String, Server> serverMap = new HashMap<>();
                for (ServiceConfig config: configs) {
                    Server server = allConfigRef.get().getServerByIp(config.getIp());
                    if (server != null) {
                        serverMap.put(config.getIp(), server);
                    }
                }
                if (serverMap.isEmpty()) {
                    System.out.println(String.format("未找到env[%s], service[%s] 对应服务器配置，请重新输入！", env, service));
                    CommandHelper.printHelpInfo();
                    continue;
                }

                String command = "";
                // 根据配置文件中配置的日志路径进行检索
                if (commandType == CommandType.T1) {
                    StringBuilder commandBuilder = new StringBuilder(" | grep ");
                    for (int i = 2; i < commandArgs.length; i++) {
                        commandBuilder.append(commandArgs[i]).append(" ");
                        if (i < commandArgs.length - 1) {
                            commandBuilder.append("| grep ");
                        }
                    }
                    command = commandBuilder.toString();
                }
                // 根据配置文件中配置的日志路径进行自定义grep命令检索
                if (commandType == CommandType.T2) {
                    command = " | grep " + commandArgs[3];
                }
                // 执行自定义脚本命令
                if (commandType == CommandType.T3) {
                    command = commandArgs[3];
                }

                for (ServiceConfig config: configs) {
                    if (serverMap.containsKey(config.getIp())) {
                        Server server = serverMap.get(config.getIp());
                        // 根据配置文件中配置的日志路径进行自定义grep命令检索
                        if (commandType == CommandType.T1 || commandType == CommandType.T2) {
                            executeLogQuery(server, config, command);
                        }
                        // 执行自定义脚本命令
                        if (commandType == CommandType.T3) {
                            executeBashCommand(server, config, command);
                        }
                    }
                }
            } catch (Throwable throwable) {
                if (logger.isErrorEnabled()) {
                    logger.error("main线程处理异常", throwable);
                }
                throwable.printStackTrace();
            }
        }
    }

    private static void executeLogQuery(Server server, ServiceConfig config, String commandSuffix) {
        List<String> logPaths = config.getLogPaths();
        String[] commands = new String[logPaths.size()];
        for (int i = 0, size = logPaths.size(); i < size; i++) {
            commands[i] = "cat " + logPaths.get(i) + commandSuffix;
        }
        SSHRequest sshRequest = new SSHRequest(server.getIp(), server.getPort(), server.getUsername(), server.getPassword(), commands);
        List<SSHResponse> sshResponseList = SSHExecutor.execMultiCommand(sshRequest);
        for (SSHResponse sshResponse: sshResponseList) {
            System.out.println("BGN ===================================================");
            System.out.println(String.format("%s -> %s", sshResponse.getRemoteHost(), sshResponse.getCommand()));
            System.out.println("日志查询是否成功? " + sshResponse.isSuccess());
            if (sshResponse.isSuccess()) {
                System.out.println("日志查询返回结果: \n" + sshResponse.getOutputContent());
            } else {
                System.out.println("日志查询返回异常: \n" + sshResponse.getErrorContent());
            }
            System.out.println("FIN ===================================================");
        }
    }

    private static void executeBashCommand(Server server, ServiceConfig config, String command) {
        SSHRequest sshRequest = new SSHRequest(server.getIp(), server.getPort(), server.getUsername(), server.getPassword(), command);
        SSHResponse sshResponse = SSHExecutor.execCommand(sshRequest);
        System.out.println("BGN ===================================================");
        System.out.println(String.format("%s -> %s", sshResponse.getRemoteHost(), sshResponse.getCommand()));
        System.out.println("命令执行是否成功? " + sshResponse.isSuccess());
        if (sshResponse.isSuccess()) {
            System.out.println("命令执行返回结果: \n" + sshResponse.getOutputContent());
        } else {
            System.out.println("命令执行返回异常: \n" + sshResponse.getErrorContent());
        }
        System.out.println("FIN ===================================================");
    }

    private static AllConfig loadConfig(String[] args) {
        String cfgSrvPath = System.getProperty("config.server");
        if (StringUtils.isNotEmpty(cfgSrvPath)) {
            if (logger.isInfoEnabled()) {
                logger.info("从配置服务[{}]获取配置信息", cfgSrvPath);
            }
            System.out.println(String.format("从配置服务[%s]获取配置信息", cfgSrvPath));
            AllConfig allConfig = AllConfig.resolveFromConfigServer(cfgSrvPath);
            if (allConfig != null) {
                try {
                    File tmpFile = FileUtils.newFile(String.format("%sconfig%s.json"
                            , System.getProperty("user.dir") + File.separator + "config" + File.separator
                            , new SimpleDateFormat("yyyyMMddHHmmssSSS").format(new Date())));
                    String json = JSON.toJSONString(allConfig, true);
                    IOUtils.write(json, tmpFile);
                    if (logger.isInfoEnabled()) {
                        logger.info("从配置服务[{}]获取配置信息持久化到配置文件: {}", cfgSrvPath, tmpFile.getAbsolutePath());
                    }
                } catch (Throwable throwable) {
                    if (logger.isErrorEnabled()) {
                        logger.error("从配置服务[{}]获取配置信息持久化到配置文件失败", cfgSrvPath, throwable);
                    }
                }
            }
            return allConfig;
        }
        String cfgPath = System.getProperty("user.dir") + File.separator + "config.json";
        if (logger.isInfoEnabled()) {
            logger.info("从本地[{}]获取配置信息", cfgPath);
        }
        System.out.println(String.format("从本地[%s]获取配置信息", cfgPath));
        return AllConfig.resolveFromFile(new File(cfgPath));
    }

}
