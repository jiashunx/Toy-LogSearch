package io.github.jiashunx.toy.LogSearch.model;

import java.util.List;

/**
 * @author jiashunx
 */
public class ServiceConfig {

    private String ip;

    private List<String> logPaths;

    public String getIp() {
        return ip;
    }

    public void setIp(String ip) {
        this.ip = ip;
    }

    public List<String> getLogPaths() {
        return logPaths;
    }

    public void setLogPaths(List<String> logPaths) {
        this.logPaths = logPaths;
    }
}
