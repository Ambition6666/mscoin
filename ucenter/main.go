package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"grpc-common/ucenter/types/asset"
	"grpc-common/ucenter/types/login"
	"grpc-common/ucenter/types/member"
	"grpc-common/ucenter/types/register"
	"path/filepath"
	"ucenter/internal/config"
	assetServer "ucenter/internal/server/asset"
	loginServer "ucenter/internal/server/login"
	memberServer "ucenter/internal/server/member"
	registerServer "ucenter/internal/server/register"
	"ucenter/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/register.yaml", "the config file")

func main() {
	flag.Parse()

	dir := filepath.Join("D:\\ms-coin-exchange-go\\", "common", "logs", "ucenter")
	fmt.Println("日志位置为", dir)

	// 替换日志格式
	logx.MustSetup(logx.LogConf{ServiceName: "ucenter", Stat: false, Encoding: "plain", Mode: "file", Path: dir})
	defer logx.Close()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		register.RegisterRegisterServer(grpcServer, registerServer.NewRegisterServer(ctx))
		login.RegisterLoginServer(grpcServer, loginServer.NewLoginServer(ctx))
		asset.RegisterAssetServer(grpcServer, assetServer.NewAssetServer(ctx))
		member.RegisterMemberServer(grpcServer, memberServer.NewMemberServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
