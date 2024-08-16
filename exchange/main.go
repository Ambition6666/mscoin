package main

import (
	"exchange/internal/config"
	orderserver "exchange/internal/server/order"
	"exchange/internal/svc"
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc-common/exchange/types/order"
	"path/filepath"
)

var configFile = flag.String("f", "etc/exchange.yaml", "the config file")

func main() {
	flag.Parse()

	dir := filepath.Join("D:\\ms-coin-exchange-go\\", "common", "logs", "exchange")
	fmt.Println("日志位置为", dir)

	// 替换日志格式
	logx.MustSetup(logx.LogConf{ServiceName: "ucenter", Stat: false, Encoding: "plain", Mode: "file", Path: dir})
	defer logx.Close()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		order.RegisterOrderServer(grpcServer, orderserver.NewOrderServer(ctx))
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
