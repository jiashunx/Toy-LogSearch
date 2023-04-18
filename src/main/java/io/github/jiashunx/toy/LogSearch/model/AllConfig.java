package io.github.jiashunx.toy.LogSearch.model;

import com.alibaba.fastjson.JSON;
import io.github.jiashunx.masker.rest.framework.util.StringUtils;
import org.apache.commons.io.IOUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.ByteArrayOutputStream;
import java.io.File;
import java.io.FileInputStream;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;
import java.util.Objects;

/**
 * @author jiashunx
 */
public class AllConfig {

    private static final Logger logger = LoggerFactory.getLogger(AllConfig.class);

    private List<Server> servers;

    private List<Service> services;

    public List<Server> getServers() {
        return servers;
    }

    public void setServers(List<Server> servers) {
        this.servers = servers;
    }

    public List<Service> getServices() {
        return services;
    }

    public void setServices(List<Service> services) {
        this.services = services;
    }

    public static AllConfig resolveFromFile(File jsonFile) {
        try (FileInputStream inputStream = new FileInputStream(jsonFile.getAbsolutePath());
             ByteArrayOutputStream outputStream = new ByteArrayOutputStream();) {
            IOUtils.copy(inputStream, outputStream);
            String jsonContent = new String(outputStream.toByteArray(), StandardCharsets.UTF_8);
            return resolveFromContent(jsonContent);
        } catch (Throwable throwable) {
            if (logger.isErrorEnabled()) {
                logger.error("解析配置文件异常", throwable);
            }
        }
        return null;
    }

    public static AllConfig resolveFromContent(String json) {
        try {
            return JSON.parseObject(json, AllConfig.class);
        } catch (Throwable throwable) {
            if (logger.isErrorEnabled()) {
                logger.error("从json反序列化配置对象异常", throwable);
            }
        }
        return null;
    }

    public void verifyBeanInfo() {
        Objects.requireNonNull(servers, "servers can't be null");
        for (Server server: servers) {
            Objects.requireNonNull(server, "server can't be null");
            Objects.requireNonNull(server.getIp(), "server ip can't be null");
        }
        Objects.requireNonNull(services, "services can't be null");
        for (Service service: services) {
            Objects.requireNonNull(service.getEnv(), "service env can't be null");
            Objects.requireNonNull(service.getEnv(), "service env service can't be null");
            Objects.requireNonNull(service.getConfigs(), "service configs can't be null");
            for (ServiceConfig config: service.getConfigs()) {
                Objects.requireNonNull(config.getIp(), "service config ip can't be null");
                Objects.requireNonNull(config.getLogPaths(), "service config logPaths can't be null");
            }
        }
    }

    public Server getServerByIp(String ip) {
        for (Server server: servers) {
            if (server.getIp().equals(ip)) {
                return server;
            }
        }
        return null;
    }

    public List<ServiceConfig> getEnvServiceConfigs(String env, String service) {
        List<ServiceConfig> configs = new ArrayList<>();
        for (Service s: services) {
            if (s.getEnv().equals(env) && s.getService().equals(service)) {
                configs.addAll(s.getConfigs());
            }
        }
        return configs;
    }

    public void printConfigInfo() {
        for (Server server: servers) {
            System.out.println(String.format("server: %s:%d@%s/%s", server.getIp(), server.getPort(), server.getUsername(), server.getPassword()));
        }
        for (Service service: services) {
            List<ServiceConfig> configs = service.getConfigs();
            for (ServiceConfig config: configs) {
                System.out.println(String.format("env: %s, service: %s, ip: %s, logPaths: %s", service.getEnv(), service.getService(), config.getIp(), JSON.toJSON(config.getLogPaths())));
            }
        }
    }

}
