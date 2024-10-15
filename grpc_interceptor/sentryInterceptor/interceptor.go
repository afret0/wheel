package sentryInterceptor

import (
	"context"
	"errors"
	"fmt"
	"github.com/afret0/wheel/frame/frameErr"
	"github.com/getsentry/sentry-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

//func SentryInterceptor() grpc.UnaryServerInterceptor {
//
//	return grpc_middleware.ChainUnaryServer(
//		grpc_sentry.UnaryServerInterceptor(grpc_sentry.WithRepanicOption(true)))
//}

func SentryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		opId := ""
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if val, exists := md["opid"]; exists && len(val) > 0 {
				opId = val[0]
			}
		}
		if err != nil {
			if !frameErr.IsFrameErr(err) {
				ErrInfo := fmt.Sprintf("error: %s, req: %+v, info: %+v, opId: %s", err, req, info, opId)
				sentry.CaptureException(errors.New(ErrInfo))
			}
		}

		return resp, err
	}
}
