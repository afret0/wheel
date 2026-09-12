package tool

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

// GrpcCtx 是发起 RPC 前的必经路径, 即使建连时漏加 DialOption 也应带上 caller。
func TestGrpcCtxInjectsCaller(t *testing.T) {
	t.Setenv("APP_NAME", "order-service")

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs("caller", "upstream-service"),
	)

	md, ok := metadata.FromOutgoingContext(GrpcCtx(ctx))
	if !ok {
		t.Fatal("outgoing metadata is missing")
	}
	if got := md.Get("caller"); len(got) != 1 || got[0] != "order-service" {
		t.Fatalf("caller = %v, want [order-service]", got)
	}
}

func TestGrpcCtxDropsStaleCallerWithoutAppName(t *testing.T) {
	t.Setenv("APP_NAME", "")

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs("caller", "upstream-service"),
	)

	md, ok := metadata.FromOutgoingContext(GrpcCtx(ctx))
	if !ok {
		t.Fatal("outgoing metadata is missing")
	}
	if got := md.Get("caller"); len(got) != 0 {
		t.Fatalf("caller = %v, want empty", got)
	}
}
