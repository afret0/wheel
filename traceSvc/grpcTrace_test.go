package traceSvc

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestGrpcCallerUnaryClientInterceptorOverwritesStaleCaller(t *testing.T) {
	t.Setenv("HOSTNAME", "order-service-6d8f7b9c4d-x2abc")

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
	t.Setenv("HOSTNAME", "order-service-6d8f7b9c4d-x2abc")

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
	t.Setenv("HOSTNAME", "order-service-6d8f7b9c4d-x2abc")

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

func TestGrpcCallerUnaryClientInterceptorKeepsOtherOutgoingKeys(t *testing.T) {
	t.Setenv("HOSTNAME", "order-service-6d8f7b9c4d-x2abc")

	ctx := metadata.NewOutgoingContext(
		context.Background(),
		metadata.Pairs("_uid", "10086", "caller", "upstream-service"),
	)

	err := grpcCallerUnaryClientInterceptor(
		ctx,
		"/test.Service/Method",
		nil,
		nil,
		nil,
		func(ctx context.Context, _ string, _, _ any, _ *grpc.ClientConn, _ ...grpc.CallOption) error {
			md, _ := metadata.FromOutgoingContext(ctx)
			if got := md.Get("_uid"); len(got) != 1 || got[0] != "10086" {
				t.Fatalf("_uid = %v, want [10086]", got)
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
