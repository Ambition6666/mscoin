package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"path/filepath"

	"ucenter-api/internal/config"
	"ucenter-api/internal/handler"
	"ucenter-api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/ucenterapi-api.yaml", "the config file")

func main() {
	flag.Parse()

	dir := filepath.Join("D:\\ms-coin-exchange-go\\", "common", "logs", "ucenter-api")
	fmt.Println("日志位置为", dir)

	// 替换日志格式
	logx.MustSetup(logx.LogConf{ServiceName: "ucenter-api", Stat: false, Encoding: "plain", Mode: "file", Path: dir})
	defer logx.Close()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithCustomCors(func(header http.Header) {
		header.Set("Access-Control-Allow-Headers", "DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization,token,x-auth-token")
	}, func(writer http.ResponseWriter) {}, "http://localhost:8080"))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)

	router := handler.NewRouters(server)
	handler.RegisterHandlers(router, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
