package sentryInterceptor

import (
	"context"
	"errors"
	"fmt"
	"github.com/afret0/wheel/frame/frameErr"
	"github.com/afret0/wheel/log"
	"github.com/getsentry/sentry-go"
	"google.golang.org/grpc"
)

// func SentryInterceptor() grpc.UnaryServerInterceptor {
//
//	return grpc_middleware.ChainUnaryServer(
//		grpc_sentry.UnaryServerInterceptor(grpc_sentry.WithRepanicOption(true)))
func SentryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			if !frameErr.IsFrameErr(err) {
				ErrInfo := fmt.Sprintf("error: %s, req: %+v, info: %+v, opId: %s", err, req, info, log.OpId(ctx))
				sentry.CaptureException(errors.New(ErrInfo))
			}
		}

		return resp, err
	}
}
