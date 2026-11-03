package frame

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func ctxWithRequest(uri string) context.Context {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RequestURI = uri

	return context.WithValue(context.Background(), gin.ContextRequestKey, req)
}

func TestUrlContainsCMS(t *testing.T) {
	cases := []struct {
		name string
		ctx  context.Context
		want bool
	}{
		{
			name: "纯 gRPC 调用, ctx 上没有 gin 请求, 判定为非 CMS 入口且不 panic",
			ctx:  context.Background(),
			want: false,
		},
		{
			name: "ctx 上挂的不是 http.Request",
			ctx:  context.WithValue(context.Background(), gin.ContextRequestKey, "not-a-request"),
			want: false,
		},
		{
			name: "CMS 后台入口",
			ctx:  ctxWithRequest("/CMS/asset/beanCtrl/exchange"),
			want: true,
		},
		{
			name: "用户端 HTTP 入口",
			ctx:  ctxWithRequest("/asset/beanCtrl/exchange"),
			want: false,
		},
		{
			name: "匹配大小写敏感, 小写 cms 不算 CMS 入口",
			ctx:  ctxWithRequest("/cms/asset/beanCtrl/exchange"),
			want: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := UrlContainsCMS(c.ctx); got != c.want {
				t.Errorf("UrlContainsCMS() = %t, want %t", got, c.want)
			}
		})
	}
}
