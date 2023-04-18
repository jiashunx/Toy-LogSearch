package io.github.jiashunx.toy.LogSearch.model;

import io.github.jiashunx.masker.rest.framework.util.IOUtils;
import static org.junit.Assert.*;

import org.junit.Test;

import java.util.List;

/**
 * @author jiashunx
 */
public class AllConfigTest {

    @Test
    public void testResolveConfig() {
        String json = IOUtils.loadContentFromClasspath("config.json");
        AllConfig allConfig = AllConfig.resolveFromContent(json);
        assertNotNull(allConfig);
        List<Server> servers = allConfig.getServers();
        assertNotNull(servers);
        assertEquals(2, servers.size());
        assertEquals("192.168.43.36", servers.get(0).getIp());
        assertEquals("192.168.43.36", servers.get(1).getIp());
        List<Service> services = allConfig.getServices();
        assertEquals(4, services.size());
        assertEquals("sit2", services.get(1).getEnv());
        assertEquals("print", services.get(1).getService());
        allConfig.printConfigInfo();
    }
}
