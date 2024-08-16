package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest/chain"
	"github.com/zeromicro/go-zero/rest/router"
	"market-api/internal/ws"
	"net/http"
	"path/filepath"

	"market-api/internal/config"
	"market-api/internal/handler"
	"market-api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/marketapi-api.yaml", "the config file")

func main() {
	flag.Parse()

	dir := filepath.Join("D:\\ms-coin-exchange-go\\", "common", "logs", "market-api")
	fmt.Println("日志位置为", dir)
	// 替换日志格式
	logx.MustSetup(logx.LogConf{ServiceName: "market-api", Stat: false, Encoding: "plain", Mode: "file", Path: dir})
	defer logx.Close()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	wsServer := ws.NewWebSocketServer(router.NewRouter(), "/socket.io")
	server := rest.MustNewServer(
		c.RestConf,
		rest.WithChain(chain.New(wsServer.ServerHandler)),
		//rest.WithRouter(wsServer),
		rest.WithCustomCors(func(header http.Header) {
			header.Set("Access-Control-Allow-Headers", "DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization,token,x-auth-token")
		}, func(writer http.ResponseWriter) {}, "http://localhost:8080"),
	)
	defer server.Stop()

	ctx := svc.NewServiceContext(c, wsServer)
	r := handler.NewRouters(server)
	handler.RegisterHandlers(r, ctx)

	handler.RegisterWsHandlers(r, ctx)
	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	group := service.NewServiceGroup()
	group.Add(server)
	group.Add(wsServer)
	group.Start()
}
