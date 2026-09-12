package traceSvc

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestGrpcCallerUnaryClientInterceptorOverwritesStaleCaller(t *testing.T) {
	t.Setenv("APP_NAME", "order-service")

	// 模拟直接透传入站 metadata 的场景: caller 仍是上游服务名。
	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("opid", "existing-op", "caller", "upstream-service"),
	)

	err := grpcCallerUnaryClientInterceptor(
		ctx,
		"/test.Service/Method",
		nil,
		nil,
		nil,
		func(ctx context.Context, _ string, _, _ any, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
			md, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				t.Fatal("outgoing metadata is missing")
			}
			if got := md.Get("caller"); len(got) != 1 || got[0] != "order-service" {
				t.Fatalf("caller = %v, want [order-service]", got)
			}
			if got := md.Get("opid"); len(got) != 1 || got[0] != "existing-op" {
				t.Fatalf("opid = %v, want [existing-op]", got)
			}
			return nil
		},
	)
	if err != nil {
		t.Fatalf("interceptor returned error: %v", err)
	}
}

func TestGrpcCallerUnaryClientInterceptorWithoutOutgoingMetadata(t *testing.T) {
	t.Setenv("APP_NAME", "order-service")

	err := grpcCallerUnaryClientInterceptor(
		context.Background(),
		"/test.Service/Method",
		nil,
		nil,
		nil,
		func(ctx context.Context, _ string, _, _ any, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
			md, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				t.Fatal("outgoing metadata is missing")
			}
			if got := md.Get("caller"); len(got) != 1 || got[0] != "order-service" {
				t.Fatalf("caller = %v, want [order-service]", got)
			}
			return nil
		},
	)
	if err != nil {
		t.Fatalf("interceptor returned error: %v", err)
	}
}

func TestGrpcCallerStreamClientInterceptor(t *testing.T) {
	t.Setenv("APP_NAME", "order-service")

	_, err := grpcCallerStreamClientInterceptor(
		context.Background(),
		&grpc.StreamDesc{},
		nil,
		"/test.Service/Stream",
		func(ctx context.Context, _ *grpc.StreamDesc, _ *grpc.ClientConn, _ string, _ ...grpc.CallOption) (grpc.ClientStream, error) {
			md, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				t.Fatal("outgoing metadata is missing")
			}
			if got := md.Get("caller"); len(got) != 1 || got[0] != "order-service" {
				t.Fatalf("caller = %v, want [order-service]", got)
			}
			return nil, nil
		},
	)
	if err != nil {
		t.Fatalf("interceptor returned error: %v", err)
	}
}

// APP_NAME 缺失时无法得知自身身份, 保留上游 caller 会产生误导性日志。
func TestGrpcCallerUnaryClientInterceptorDropsStaleCallerWithoutAppName(t *testing.T) {
	t.Setenv("APP_NAME", "")

	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("caller", "upstream-service"),
	)

	err := grpcCallerUnaryClientInterceptor(
		ctx,
		"/test.Service/Method",
		nil,
		nil,
		nil,
		func(ctx context.Context, _ string, _, _ any, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
			md, _ := metadata.FromOutgoingContext(ctx)
			if got := md.Get("caller"); len(got) != 0 {
				t.Fatalf("caller = %v, want empty", got)
			}
			return nil
		},
	)
	if err != nil {
		t.Fatalf("interceptor returned error: %v", err)
	}
}

func TestGrpcCallerOptionsCoversUnaryAndStream(t *testing.T) {
	if got := len(GrpcCallerOptions()); got != 2 {
		t.Fatalf("GrpcCallerOptions len = %d, want 2", got)
	}
}
