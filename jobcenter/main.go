package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"jobcenter/internal/config"
	"jobcenter/internal/svc"
	"jobcenter/internal/task"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var configFile = flag.String("f", "etc/conf.yaml", "the config file")

func main() {
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	t := task.NewTask(ctx)
	t.Run()

	go func() {
		exit := make(chan os.Signal)
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-exit:
			log.Println("监听到中断信号，终止程序")
			ctx.MongoClient.Disconnect()
			t.Stop()
		}
	}()
	t.StartBlocking()
}
