package tool

import (
	"context"
	"strings"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestGrpcCtxInjectsCaller(t *testing.T) {
	t.Setenv("HOSTNAME", "order-6d8f7b9c4d-x2abc")

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs(MetaCaller, "upstream-service"),
	)

	md, ok := metadata.FromOutgoingContext(GrpcCtx(ctx))
	if !ok {
		t.Fatal("outgoing metadata is missing")
	}

	if got := md.Get(MetaCaller); len(got) != 1 || got[0] != "order" {
		t.Fatalf("caller = %v, want [order]", got)
	}
}

func TestGrpcCtxInjectsCallerWhenHostnameMissing(t *testing.T) {
	t.Setenv("HOSTNAME", "")

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs(MetaCaller, "upstream-service"),
	)

	md, ok := metadata.FromOutgoingContext(GrpcCtx(ctx))
	if !ok {
		t.Fatal("outgoing metadata is missing")
	}

	got := md.Get(MetaCaller)
	if len(got) == 1 && got[0] == "upstream-service" {
		t.Fatalf("caller = %v, 沿用了上游服务名", got)
	}
	if len(got) == 1 && isGeneratedUuid(got[0]) {
		t.Fatalf("caller = %v, 写入了 HostId 兜底的随机 uuid", got)
	}
}

func TestGrpcCtxKeepsIncomingMetadataUntouched(t *testing.T) {
	t.Setenv("HOSTNAME", "order-6d8f7b9c4d-x2abc")

	incoming := metadata.Pairs(MetaCaller, "upstream-service", "opid", "op-1")
	ctx := metadata.NewIncomingContext(context.Background(), incoming)

	if _, ok := metadata.FromOutgoingContext(GrpcCtx(ctx)); !ok {
		t.Fatal("outgoing metadata is missing")
	}

	if got := incoming.Get(MetaCaller); len(got) != 1 || got[0] != "upstream-service" {
		t.Fatalf("incoming caller = %v, want [upstream-service]", got)
	}
	if got := incoming.Get("opid"); len(got) != 1 || got[0] != "op-1" {
		t.Fatalf("incoming opid = %v, want [op-1]", got)
	}
}

func TestGrpcCtxCarriesOpId(t *testing.T) {
	t.Setenv("HOSTNAME", "order-6d8f7b9c4d-x2abc")

	ctx := context.WithValue(context.Background(), "opId", "op-1")

	md, ok := metadata.FromOutgoingContext(GrpcCtx(ctx))
	if !ok {
		t.Fatal("outgoing metadata is missing")
	}

	got := md.Get("opid")
	if len(got) != 1 || got[0] == "" {
		t.Fatalf("opid = %v, want 单值", got)
	}
	if !strings.HasPrefix(got[0], "op-1-") {
		t.Fatalf("opid = %q, 期望在入参 opId 上追加 spanId", got[0])
	}
}
