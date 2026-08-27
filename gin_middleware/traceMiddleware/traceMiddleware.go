// Package traceMiddleware 提供 gin 的 OpenTelemetry 链路追踪中间件。
//
// 它负责整条链路中最关键的入站一环: 从请求头中提取上游透传的 W3C traceparent
// (Extract), 并以此为父级创建当前服务的根 span。缺少这一步时, 上游 Inject 进
// 请求头的 traceId 会被丢弃, 每个服务都会各自开启一条新链路。
//
// 中间件通过 engine.Use 注册后对所有路由生效, 包括直接使用 gin 原生方式注册的
// 旧接口, 无需逐个 handler 手工埋点。
package traceMiddleware

import (
	"net/http"
	"strings"

	"github.com/afret0/wheel/tool"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "gin"

type Option struct {
	// Service 服务名, 仅用于 span 属性标注; 上报使用的服务名由 TracerProvider 的 resource 决定
	Service string `json:"service"`
	// WhiteList 命中(子串匹配)的请求路径不会产生 span, 适合排除 pprof、健康检查等噪声路由
	WhiteList []string `json:"whiteList"`
}

// TraceMiddleware 返回链路追踪中间件。
// 未开启 TRACE 环境变量时返回空实现, 不产生任何额外开销。
func TraceMiddleware(opts ...*Option) gin.HandlerFunc {
	opt := new(Option)
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}

	if !tool.EnvEnabled("TRACE") {
		return func(c *gin.Context) { c.Next() }
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		for _, w := range opt.WhiteList {
			if w != "" && strings.Contains(path, w) {
				c.Next()
				return
			}
		}

		// 关键: 先从入站请求头还原上游 span, 再以其为父级开启本服务的根 span
		ctx := otel.GetTextMapPropagator().Extract(
			c.Request.Context(),
			propagation.HeaderCarrier(c.Request.Header),
		)

		ctx, span := otel.Tracer(tracerName).Start(ctx, spanName(c),
			trace.WithSpanKind(trace.SpanKindServer),
		)
		defer span.End()

		span.SetAttributes(
			attribute.String("opId", opId(c)),
			attribute.String("http.request.method", c.Request.Method),
			attribute.String("url.path", path),
			attribute.String("server.address", c.Request.Host),
			attribute.String("client.address", c.ClientIP()),
		)
		if opt.Service != "" {
			span.SetAttributes(attribute.String("service", opt.Service))
		}

		// 写回 Request, 使后续 c.Request.Context() 与 gin.Context.Value 都能取到该 span
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		status := c.Writer.Status()
		span.SetAttributes(attribute.Int("http.response.status_code", status))
		if status >= http.StatusInternalServerError {
			span.SetStatus(codes.Error, http.StatusText(status))
		}
		if len(c.Errors) > 0 {
			span.SetStatus(codes.Error, c.Errors.String())
		}
	}
}

// spanName 优先使用路由模板(如 /room/v1/user/:id), 避免路径参数导致 span 名称基数爆炸
func spanName(c *gin.Context) string {
	if p := c.FullPath(); p != "" {
		return c.Request.Method + " " + p
	}
	return c.Request.Method + " " + c.Request.URL.Path
}

// opId 复用 loggerMiddleware 写入的业务链路 id, 便于日志与 span 互相定位
func opId(c *gin.Context) string {
	if v := tool.OpIdWithoutDefault(c); v != "" {
		return v
	}
	return c.GetHeader("opId")
}
