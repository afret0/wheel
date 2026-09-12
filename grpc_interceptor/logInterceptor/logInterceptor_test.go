package logInterceptor

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/afret0/wheel/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestInterceptorLogsCallerAndLatencyWhenEnabled(t *testing.T) {
	t.Setenv("ENABLE_CALLER", "true")

	var output bytes.Buffer
	logger := log.GetMiddleWareLogger()
	originalOutput := logger.Out
	logger.SetOutput(&output)
	t.Cleanup(func() {
		logger.SetOutput(originalOutput)
	})

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs("caller", "gateway-service"),
	)
	_, err := Interceptor()(ctx, "request", &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/Method",
	}, func(context.Context, any) (any, error) {
		return "response", nil
	})
	if err != nil {
		t.Fatalf("interceptor returned error: %v", err)
	}

	var entry map[string]any
	if err := json.Unmarshal(output.Bytes(), &entry); err != nil {
		t.Fatalf("unmarshal log entry: %v", err)
	}
	if got := entry["caller"]; got != "gateway-service" {
		t.Fatalf("caller = %v, want gateway-service", got)
	}
	if _, ok := entry["latency"]; !ok {
		t.Fatal("latency field is missing")
	}
	if _, ok := entry["latencyT"]; ok {
		t.Fatal("legacy latencyT field should not be logged")
	}
}

func TestInterceptorOmitsCallerAndLatencyWhenDisabled(t *testing.T) {
	t.Setenv("ENABLE_CALLER", "false")

	var output bytes.Buffer
	logger := log.GetMiddleWareLogger()
	originalOutput := logger.Out
	logger.SetOutput(&output)
	t.Cleanup(func() {
		logger.SetOutput(originalOutput)
	})

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs("caller", "gateway-service"),
	)
	_, err := Interceptor()(ctx, "request", &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/Method",
	}, func(context.Context, any) (any, error) {
		return "response", nil
	})
	if err != nil {
		t.Fatalf("interceptor returned error: %v", err)
	}

	var entry map[string]any
	if err := json.Unmarshal(output.Bytes(), &entry); err != nil {
		t.Fatalf("unmarshal log entry: %v", err)
	}
	if _, ok := entry["caller"]; ok {
		t.Fatal("caller should not be logged")
	}
	if _, ok := entry["latency"]; ok {
		t.Fatal("latency should not be logged")
	}
}
