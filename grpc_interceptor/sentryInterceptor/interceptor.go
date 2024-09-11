package sentryInterceptor

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_sentry "github.com/johnbellone/grpc-middleware-sentry"
	"google.golang.org/grpc"
)

func SentryInterceptor() grpc.UnaryServerInterceptor {

	return grpc_middleware.ChainUnaryServer(
		grpc_sentry.UnaryServerInterceptor(grpc_sentry.WithRepanicOption(true)))
}
