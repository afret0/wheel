package frame

import (
	"context"
	"net/http"

	"github.com/afret0/wheel/constant"
	"github.com/gin-gonic/gin"
)

func Request(ctx context.Context) *http.Request {
	req := ctx.Value(gin.ContextRequestKey)
	if req == nil {
		return nil
	}

	req1, ok := req.(*http.Request)
	if !ok {
		return nil
	}

	return req1
}

func Header(ctx context.Context) http.Header {
	req := Request(ctx)
	if req == nil {
		return nil
	}

	return req.Header
}

func App(ctx context.Context) string {
	if app := ctx.Value("app"); app != nil {
		if appStr, ok := app.(string); ok {
			return appStr
		}
	}

	header := Header(ctx)
	if header == nil {
		header = make(http.Header)
	}

	app := header.Get("app")
	if app == "" {
		app = constant.Keke
	}

	return app
}

func IsHayo(ctx context.Context) bool {
	app := App(ctx)
	return app == constant.Hayo
}

func IsKeke(ctx context.Context) bool {
	app := App(ctx)
	return app == constant.Keke
}
