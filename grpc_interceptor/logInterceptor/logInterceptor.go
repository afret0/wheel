package logInterceptor

import (
	"context"
	"errors"
	"fmt"
	"github.com/getsentry/sentry-go"
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

type Option struct {
	Service        string `json:"service"`
	ReportToSentry bool   `json:"reportToSentry"`
	RePanic        bool   `json:"rePanic"`
	Debug          bool   `json:"debug"`
}

type Opt = Option

func Interceptor(opts ...*Option) grpc.UnaryServerInterceptor {
	opt := new(Option)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
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
			if opt.Debug {
				return
			}

			if r := recover(); r != nil {
				// 记录panic信息
				stack := string(debug.Stack())
				lg.WithFields(logrus.Fields{
					"panic": r,
					"stack": stack,
				}).Error("Panic occurred")

				if opt.ReportToSentry {
					go sentry.CaptureException(errors.New(fmt.Sprintf("Panic occurred: %s", stack)))
				}

				if opt.RePanic {
					panic(fmt.Sprintf("Panic occurred: %s, \nstack: %s", r, stack))
				}

				err = status.Errorf(codes.Internal, "Panic occurred: %#v, stack: %s", r, stack)
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
}

//func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
//	opId := strings.ReplaceAll(uuid.New().String(), "-", "")
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
//	ctx = context.WithValue(ctx, "opId", opId)
//
//	clientIP := ""
//	if p, ok := peer.FromContext(ctx); ok {
//		clientIP = p.Addr.String()
//	}
//
//	startT := time.Now()
//
//	lg := log.GetMiddleWareLogger().WithFields(logrus.Fields{"opId": opId, "info": info, "clientIP": clientIP, "req": req, "_uid": uid})
//
//	defer func() {
//		if r := recover(); r != nil {
//			// 记录panic信息
//			stack := string(debug.Stack())
//			lg.WithFields(logrus.Fields{
//				"panic": r,
//				"stack": stack,
//			}).Error("Panic occurred")
//
//			go sentry.CaptureException(errors.New(fmt.Sprintf("Panic occurred: %s", stack)))
//
//			err = status.Errorf(codes.Internal, "Panic occurred: %v", r)
//			resp = nil
//		}
//	}()
//
//	resp, err = handler(ctx, req)
//
//	endT := time.Now()
//	latencyT := endT.Sub(startT)
//	lg.WithFields(logrus.Fields{
//		"latencyT": latencyT.Milliseconds(),
//		"res":      resp,
//		"err":      err,
//	}).Info("请求日志")
//	return resp, err
//}
