package grpcClient_test

import (
	"context"
	"encoding/json"
	"net"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"

	"github.com/afret0/wheel/grpcClient"
	"github.com/afret0/wheel/grpc_interceptor/logInterceptor"
	"github.com/afret0/wheel/log"
	"github.com/afret0/wheel/tool"
)

func TestCallerFlowsFromClientToServerLog(t *testing.T) {
	t.Setenv("HOSTNAME", "service-a-6d8f7b9c4d-x2abc")
	t.Setenv("ENABLE_CALLER", "true")

	logs := captureMiddlewareLog(t)
	target, dialer := startHealthServer(t)

	conn, err := grpcClient.New(target, grpc.WithContextDialer(dialer))
	if err != nil {
		t.Fatalf("grpcClient.New: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := grpc_health_v1.NewHealthClient(conn).Check(ctx, &grpc_health_v1.HealthCheckRequest{}); err != nil {
		t.Fatalf("health check: %v", err)
	}

	want := "service-a"

	entry := lastLogEntry(t, logs.String())
	if got := entry["caller"]; got != want {
		t.Fatalf("caller = %v, want %v", got, want)
	}
	if _, ok := entry["latency"]; !ok {
		t.Fatal("latency field is missing")
	}
}

// 即使建连时漏掉了 wheel 的标准 DialOption, 只要走 tool.GrpcCtx 构造 ctx,
// caller 依然会被带上。
func TestCallerFlowsWithoutDialOptions(t *testing.T) {
	t.Setenv("HOSTNAME", "service-a-6d8f7b9c4d-x2abc")
	t.Setenv("ENABLE_CALLER", "true")

	logs := captureMiddlewareLog(t)
	target, dialer := startHealthServer(t)

	conn, err := grpc.NewClient(target, grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc.NewClient: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ctx = tool.GrpcCtx(ctx)

	if _, err := grpc_health_v1.NewHealthClient(conn).Check(ctx, &grpc_health_v1.HealthCheckRequest{}); err != nil {
		t.Fatalf("health check: %v", err)
	}

	want := "service-a"

	entry := lastLogEntry(t, logs.String())
	if got := entry["caller"]; got != want {
		t.Fatalf("caller = %v, want %v", got, want)
	}
}

func startHealthServer(t *testing.T) (string, func(context.Context, string) (net.Conn, error)) {
	t.Helper()

	lis := bufconn.Listen(1024 * 1024)
	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(logInterceptor.Interceptor()))
	grpc_health_v1.RegisterHealthServer(srv, health.NewServer())

	go func() {
		_ = srv.Serve(lis)
	}()
	t.Cleanup(srv.Stop)

	return "passthrough:///bufnet", func(ctx context.Context, _ string) (net.Conn, error) {
		return lis.DialContext(ctx)
	}
}

func captureMiddlewareLog(t *testing.T) *strings.Builder {
	t.Helper()

	buf := new(strings.Builder)
	logger := log.GetMiddleWareLogger()
	original := logger.Out
	logger.SetOutput(buf)
	t.Cleanup(func() { logger.SetOutput(original) })

	return buf
}

// 服务端拦截器只输出一条 "request log", 但连接建立过程可能夹杂其它日志,
// 因此取最后一条可解析的 JSON 记录。
func lastLogEntry(t *testing.T, out string) map[string]any {
	t.Helper()

	lines := strings.Split(strings.TrimSpace(out), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		var entry map[string]any
		if err := json.Unmarshal([]byte(lines[i]), &entry); err != nil {
			continue
		}
		if entry["type"] == "interceptorLog" {
			return entry
		}
	}

	t.Fatalf("no interceptor log entry found in output: %q", out)
	return nil
}

func TestCallerFlowsWithAdditionalChainUnaryInterceptor(t *testing.T) {
	t.Setenv("HOSTNAME", "service-a-6d8f7b9c4d-x2abc")
	t.Setenv("ENABLE_CALLER", "true")

	logs := captureMiddlewareLog(t)
	target, dialer := startHealthServer(t)

	called := false
	conn, err := grpcClient.New(target,
		grpc.WithContextDialer(dialer),
		grpc.WithChainUnaryInterceptor(func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			called = true
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
	)
	if err != nil {
		t.Fatalf("grpcClient.New: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := grpc_health_v1.NewHealthClient(conn).Check(ctx, &grpc_health_v1.HealthCheckRequest{}); err != nil {
		t.Fatalf("health check: %v", err)
	}

	if !called {
		t.Fatal("调用方自己的 chain 拦截器没执行")
	}
	entry := lastLogEntry(t, logs.String())
	if got := entry["caller"]; got != "service-a" {
		t.Fatalf("caller = %v, want service-a", got)
	}
}

func TestCallerFlowsWithLegacyUnaryInterceptor(t *testing.T) {
	t.Setenv("HOSTNAME", "service-a-6d8f7b9c4d-x2abc")
	t.Setenv("ENABLE_CALLER", "true")

	logs := captureMiddlewareLog(t)
	target, dialer := startHealthServer(t)

	called := false
	conn, err := grpcClient.New(target,
		grpc.WithContextDialer(dialer),
		grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			called = true
			return invoker(ctx, method, req, reply, cc, opts...)
		}),
	)
	if err != nil {
		t.Fatalf("grpcClient.New: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := grpc_health_v1.NewHealthClient(conn).Check(ctx, &grpc_health_v1.HealthCheckRequest{}); err != nil {
		t.Fatalf("health check: %v", err)
	}

	if !called {
		t.Fatal("调用方自己的 unary 拦截器没执行")
	}
	entry := lastLogEntry(t, logs.String())
	if got := entry["caller"]; got != "service-a" {
		t.Fatalf("caller = %v, want service-a, 老的 WithUnaryInterceptor 把 caller 拦截器挤掉了", got)
	}
}

func TestCallerCannotBeOverriddenPerCall(t *testing.T) {
	t.Setenv("HOSTNAME", "service-a-6d8f7b9c4d-x2abc")
	t.Setenv("ENABLE_CALLER", "true")

	logs := captureMiddlewareLog(t)
	target, dialer := startHealthServer(t)

	conn, err := grpcClient.New(target, grpc.WithContextDialer(dialer))
	if err != nil {
		t.Fatalf("grpcClient.New: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ctx = metadata.AppendToOutgoingContext(ctx, tool.MetaCaller, "cms-gateway")

	if _, err := grpc_health_v1.NewHealthClient(conn).Check(ctx, &grpc_health_v1.HealthCheckRequest{}); err != nil {
		t.Fatalf("health check: %v", err)
	}

	entry := lastLogEntry(t, logs.String())
	if got := entry["caller"]; got != "service-a" {
		t.Fatalf("caller = %v, want service-a, caller 必须始终是直接发起方", got)
	}
}
