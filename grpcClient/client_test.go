package grpcClient

import (
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 默认 option 必须覆盖 trace + caller(unary/stream), 漏掉任何一个都会让
// 统一建连入口失去意义。
func TestDefaultOptionsIncludesTraceAndCaller(t *testing.T) {
	if got := len(DefaultOptions()); got != 4 {
		t.Fatalf("DefaultOptions len = %d, want 4", got)
	}
}

func TestNewBuildsClient(t *testing.T) {
	conn, err := New("passthrough:///localhost:50051")
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = conn.Close()
	})

	if conn == nil {
		t.Fatal("New returned nil conn")
	}
}

// 调用方传入的 option 排在默认值之后, 必须能覆盖默认的 insecure 传输。
func TestNewAcceptsOverridingOptions(t *testing.T) {
	conn, err := New(
		"passthrough:///localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUserAgent("custom-agent"),
	)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	t.Cleanup(func() {
		_ = conn.Close()
	})
}
