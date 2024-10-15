package logInterceptor

import (
	"context"
	"runtime/debug"
	"strings"
	"time"

	"github.com/afret0/wheel/log"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

//func Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
//	opId := tool.UUIDWithoutHyphen()
//	uid := ""
//	if md, ok := metadata.FromIncomingContext(ctx); ok {
//		if val, exists := md["opid"]; exists && len(val) > 0 {
//			opId = val[0]
//		}
//
//		if val, exists := md["_uid"]; exists && len(val) > 0 {
//			uid = val[0]
//		}
//	}
//
//	clientIP := ""
//	if p, ok := peer.FromContext(ctx); ok {
//		clientIP = p.Addr.String()
//	}
//
//	startT := time.Now()
//	//reqAt := startT.Format("2006-01-02 15:04:05")
//
//	lg := log.GetMiddleWareLogger().WithFields(logrus.Fields{"opId": opId, "info": info, "clientIP": clientIP, "req": req, "_uid": uid})
//	resp, err := handler(ctx, req)
//
//	endT := time.Now()
//	latencyT := endT.Sub(startT)
//	lg.WithFields(logrus.Fields{
//		"latencyT": latencyT.Milliseconds(),
//		"res":      resp,
//		"err":      err,
//	}).Info("请求日志")
//
//	return resp, err
//}

func Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	opId := strings.ReplaceAll(uuid.New().String(), "-", "")
	uid := ""
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val, exists := md["opid"]; exists && len(val) > 0 {
			opId = val[0]
		}

		if val, exists := md["_uid"]; exists && len(val) > 0 {
			uid = val[0]
		}
	}
	ctx = context.WithValue(ctx, "opId", opId)

	clientIP := ""
	if p, ok := peer.FromContext(ctx); ok {
		clientIP = p.Addr.String()
	}

	startT := time.Now()

	lg := log.GetMiddleWareLogger().WithFields(logrus.Fields{"opId": opId, "info": info, "clientIP": clientIP, "req": req, "_uid": uid})

	defer func() {
		if r := recover(); r != nil {
			// 记录panic信息
			stack := string(debug.Stack())
			lg.WithFields(logrus.Fields{
				"panic": r,
				"stack": stack,
			}).Error("Panic occurred")

			err = status.Errorf(codes.Internal, "Panic occurred: %v", r)
			resp = nil
		}
	}()

	resp, err = handler(ctx, req)

	endT := time.Now()
	latencyT := endT.Sub(startT)
	lg.WithFields(logrus.Fields{
		"latencyT": latencyT.Milliseconds(),
		"res":      resp,
		"err":      err,
	}).Info("请求日志")
	return resp, err
}
