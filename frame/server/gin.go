package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Option struct {
	Logger    *logrus.Logger
	Port      string
	SvrName   string
	ConfigEnv string
	Env       string
	Engine    *gin.Engine
}

func GinRun(ctx context.Context, opt *Option) {
	lg := opt.Logger.WithFields(logrus.Fields{"service": opt.SvrName, "serviceType": "gin"})
	defer func() {
		lg.Infof("http 服务退出,  os.Exit(1)")
		os.Exit(1)
	}()
	lg.Infof("%s use config, current cfg: %s", opt.SvrName, opt.ConfigEnv)

	lg.Infof("port: %s, %s start run...", opt.Port, opt.SvrName)
	if opt.Env == "pro" {
		gin.SetMode(gin.ReleaseMode)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", opt.Port),
		Handler: opt.Engine,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			lg.Panicf("listen: %s\n", err)
		}
	}()
	lg.Infof("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	lg.Infof("Shutting down server...")

	if err := srv.Shutdown(ctx); err != nil {
		lg.Panicf("Server forced to shutdown: %s", err)
	}

	lg.Println("Server exited")

}
