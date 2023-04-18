package io.github.jiashunx.toy.LogSearch.model;

import java.util.List;

/**
 * @author jiashunx
 */
public class Service {

    private String env;

    private String service;

    private List<ServiceConfig> configs;

    public String getEnv() {
        return env;
    }

    public void setEnv(String env) {
        this.env = env;
    }

    public String getService() {
        return service;
    }

    public void setService(String service) {
        this.service = service;
    }

    public List<ServiceConfig> getConfigs() {
        return configs;
    }

    public void setConfigs(List<ServiceConfig> configs) {
        this.configs = configs;
    }
}
