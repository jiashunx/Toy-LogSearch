package io.github.jiashunx.toy.LogSearch.server;

import io.github.jiashunx.masker.rest.framework.MRestServer;
import io.github.jiashunx.toy.LogSearch.model.AllConfig;
import io.github.jiashunx.toy.LogSearch.model.RuntimeEnv;
import io.netty.handler.codec.http.HttpResponseStatus;

import java.util.Objects;
import java.util.concurrent.atomic.AtomicReference;

/**
 * @author jiashunx
 */
public class EchoServer {

    private AtomicReference<AllConfig> configRef;

    public EchoServer() {}

    public EchoServer(AtomicReference<AllConfig> configRef) {
        this.configRef = Objects.requireNonNull(configRef);
    }

    public void start() {
        new Thread(() -> {
            new MRestServer(RuntimeEnv.getServerPort(), "config-server")
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
                    if (configRef.get() == null) {
                        response.write(HttpResponseStatus.NOT_FOUND);
                        return;
                    }
                    response.write(configRef.get());
                })
                .getRestServer()
                .start();
        }, "config-server").start();
    }

}
