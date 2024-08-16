package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc-common/market/types/market"
	"grpc-common/market/types/rate"
	"market/internal/config"
	rateServer "market/internal/server/exchange_rate"
	marketServer "market/internal/server/market"
	"market/internal/svc"
	"path/filepath"
)

var configFile = flag.String("f", "etc/market.yaml", "the config file")

func main() {
	flag.Parse()

	dir := filepath.Join("D:\\ms-coin-exchange-go\\", "common", "logs", "market")
	fmt.Println("日志位置为", dir)

	// 替换日志格式
	logx.MustSetup(logx.LogConf{ServiceName: "market", Stat: false, Encoding: "plain", Mode: "file", Path: dir})
	defer logx.Close()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		rate.RegisterExchangeRateServer(grpcServer, rateServer.NewExchangeRateServer(ctx))
		market.RegisterMarketServer(grpcServer, marketServer.NewMarketServer(ctx))
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
