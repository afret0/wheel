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

func TestInterceptorLogsCallerOnPanic(t *testing.T) {
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
	resp, err := Interceptor()(ctx, "request", &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/Method",
	}, func(context.Context, any) (any, error) {
		panic("boom")
	})
	if err == nil {
		t.Fatal("expected error from panicking handler")
	}
	if resp != nil {
		t.Fatalf("resp = %v, want nil", resp)
	}

	var entry map[string]any
	if err := json.Unmarshal(output.Bytes(), &entry); err != nil {
		t.Fatalf("unmarshal log entry: %v, raw: %s", err, output.String())
	}
	if got := entry["caller"]; got != "gateway-service" {
		t.Fatalf("caller = %v, want gateway-service", got)
	}
	if got := entry["panic"]; got != "boom" {
		t.Fatalf("panic = %v, want boom", got)
	}
}

func TestInterceptorKeepsIncomingMetadataUntouched(t *testing.T) {
	t.Setenv("ENABLE_CALLER", "true")

	var output bytes.Buffer
	logger := log.GetMiddleWareLogger()
	originalOutput := logger.Out
	logger.SetOutput(&output)
	t.Cleanup(func() {
		logger.SetOutput(originalOutput)
	})

	incoming := metadata.Pairs("caller", "gateway-service")
	ctx := metadata.NewIncomingContext(context.Background(), incoming)

	var outGoingCaller []string
	if _, err := Interceptor()(ctx, "request", &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/Method",
	}, func(ctx context.Context, _ any) (any, error) {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			t.Fatal("handler 看不到出站 metadata, 下游调用会丢 opId")
		}
		if got := md.Get("opid"); len(got) != 1 || got[0] == "" {
			t.Fatalf("outgoing opid = %v, want 单值", got)
		}
		outGoingCaller = md.Get("caller")
		return "response", nil
	}); err != nil {
		t.Fatalf("interceptor returned error: %v", err)
	}

	if got := incoming.Get("opid"); len(got) != 0 {
		t.Fatalf("入站 metadata 被写入 opid = %v, 污染了上游上下文", got)
	}
	if len(outGoingCaller) != 0 {
		t.Fatalf("outgoing caller = %v, 不应把上游服务名带给出站 metadata", outGoingCaller)
	}
}
