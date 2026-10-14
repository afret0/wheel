package tool

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func newTestRequest(remoteAddr string, headers map[string]string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.RemoteAddr = remoteAddr
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return req
}

func newTestGinCtx(req *http.Request) *gin.Context {
	gin.SetMode(gin.TestMode)

	c, e := gin.CreateTestContext(httptest.NewRecorder())
	if err := e.SetTrustedProxies([]string{"127.0.0.1", "10.0.0.0/8"}); err != nil {
		panic(err)
	}
	c.Request = req

	return c
}

func TestClientIP(t *testing.T) {
	cases := []struct {
		name      string
		buildCtx  func(*http.Request) context.Context
		req       *http.Request
		want      string
		wantGinIP bool // 期望值直接取 gin.Context.ClientIP(), 不硬编码 gin 的代理规则
	}{
		{
			name: "网关注入的 X-Original-Forwarded-For 优先",
			req: newTestRequest("10.0.0.1:52400", map[string]string{
				"X-Original-Forwarded-For": "203.0.113.7",
				"X-Forwarded-For":          "198.51.100.9",
				"X-Real-Ip":                "198.51.100.10",
			}),
			buildCtx: func(req *http.Request) context.Context { return newTestGinCtx(req) },
			want:     "203.0.113.7",
		},
		{
			name:     "无 X-Original-Forwarded-For 时走 gin 的解析",
			req:      newTestRequest("10.0.0.1:52400", map[string]string{"X-Forwarded-For": "198.51.100.9"}),
			buildCtx: func(req *http.Request) context.Context { return newTestGinCtx(req) },
			want:     "198.51.100.9",
		},
		{
			name:      "无代理头时返回 RemoteAddr",
			req:       newTestRequest("203.0.113.7:52400", nil),
			buildCtx:  func(req *http.Request) context.Context { return newTestGinCtx(req) },
			wantGinIP: true,
		},
		{
			name:     "gin.Context 派生的 context",
			req:      newTestRequest("10.0.0.1:52400", map[string]string{"X-Forwarded-For": "198.51.100.9"}),
			buildCtx: func(req *http.Request) context.Context { return context.WithValue(newTestGinCtx(req), "opId", "op") },
			want:     "198.51.100.9",
		},
		{
			name:     "只有 *http.Request 时按 XFF 最左侧解析",
			req:      newTestRequest("10.0.0.1:52400", map[string]string{"X-Forwarded-For": "198.51.100.9, 10.0.0.1"}),
			buildCtx: requestCtx,
			want:     "198.51.100.9",
		},
		{
			name:     "只有 *http.Request 且无 XFF 时取 X-Real-Ip",
			req:      newTestRequest("10.0.0.1:52400", map[string]string{"X-Real-Ip": "198.51.100.10"}),
			buildCtx: requestCtx,
			want:     "198.51.100.10",
		},
		{
			name:     "只有 *http.Request 且无代理头时取 RemoteAddr",
			req:      newTestRequest("203.0.113.7:52400", nil),
			buildCtx: requestCtx,
			want:     "203.0.113.7",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := tc.buildCtx(tc.req)

			want := tc.want
			if tc.wantGinIP {
				want = newTestGinCtx(tc.req).ClientIP()
			}

			if got := ClientIP(ctx); got != want {
				t.Fatalf("ClientIP() = %q, want %q", got, want)
			}
		})
	}
}

func TestClientIP_WithoutRequest(t *testing.T) {
	if got := ClientIP(context.Background()); got != "" {
		t.Fatalf("ClientIP(context.Background()) = %q, want empty", got)
	}

	// c.Request.Context() 不携带 gin 写入的 key, 取不到 IP
	c := newTestGinCtx(newTestRequest("203.0.113.7:52400", nil))
	if got := ClientIP(c.Request.Context()); got != "" {
		t.Fatalf("ClientIP(request context) = %q, want empty", got)
	}
}

func requestCtx(req *http.Request) context.Context {
	return context.WithValue(context.Background(), gin.ContextRequestKey, req)
}
