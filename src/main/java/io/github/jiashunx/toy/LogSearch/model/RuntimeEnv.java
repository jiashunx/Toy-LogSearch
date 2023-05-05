package io.github.jiashunx.toy.LogSearch.model;

import io.github.jiashunx.masker.rest.framework.util.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * @author jiashunx
 */
public class RuntimeEnv {

    private static final Logger logger = LoggerFactory.getLogger(RuntimeEnv.class);

    private static final int port;
    private static final String cfgSrvPath;
    static {
        String port0 = System.getProperty("server.port");
        if (StringUtils.isEmpty(port0)) {
            port0 = "38888";
        }
        port = Integer.parseInt(port0);
        if (logger.isInfoEnabled()) {
            logger.info("运行环境：配置服务监听端口：[{}]", port);
        }
        String cfgSrvPath0 = System.getProperty("config.server");
        if (StringUtils.isNotEmpty(cfgSrvPath0)) {
            cfgSrvPath = cfgSrvPath0;
        } else {
            cfgSrvPath = "";
        }
        if (logger.isInfoEnabled()) {
            logger.info("运行环境：远程配置服务地址：[{}]", cfgSrvPath);
        }
    }

    public static void init() {}

    public static int getServerPort() {
        return port;
    }

    public static String getCfgSrvPath() {
        return cfgSrvPath;
    }

}
