package server

import (
    "Toy-LogSearch/model"
    "encoding/json"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "net/http"
)

type EchoServer struct {
    Config *model.Config
    Address string
}

func (server *EchoServer) StartServer() {
    e := echo.New()
    e.Debug = true
    // e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.GET("/config.json", func(c echo.Context) error {
        if server.Config != nil {
            bs, err := json.Marshal(server.Config)
            if err != nil {
                return err
            }
            return c.String(http.StatusOK, string(bs))
        }
        return c.String(http.StatusNotFound, "")
    })
    e.Logger.Fatal(e.Start(server.Address))
}
