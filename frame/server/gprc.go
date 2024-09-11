package server

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type GrpcOption struct {
	Logger     *logrus.Logger
	Port       string
	SvrName    string
	ConfigEnv  string
	Env        string
	GrpcServer *grpc.Server
}

func GrpcRun(ctx context.Context, opt *GrpcOption) {
	lg := opt.Logger.WithFields(logrus.Fields{"service": opt.SvrName, "serviceType": "grpc"})
	defer func() {
		lg.Infof("grpc 服务退出,  os.Exit(1)")
		os.Exit(1)
	}()

	lg.Infof("%s use config, current cfg: %s", opt.SvrName, opt.ConfigEnv)
	lg.Infof("port: %s, %s start run...", opt.Port, opt.SvrName)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", opt.Port))
	if err != nil {
		lg.Panicf("failed to listen: %v", err)
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-stopChan
		lg.Println("Received stop signal, attempting graceful shutdown...")

		// 创建一个上下文，设置超时时间为 5 秒
		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// 优雅关闭 gRPC 服务器
		opt.GrpcServer.GracefulStop()
		lg.Println("gRPC server stopped gracefully")
	}()

	lg.Infof("Starting gRPC server... at %v", lis.Addr())
	if err := opt.GrpcServer.Serve(lis); err != nil {
		lg.Panicf("failed to serve: %v", err)
	}

	lg.Infof("Server started")

	wg.Wait()
	lg.Infof("Shutting down server...")

	lg.Infof("Server exited")
}
