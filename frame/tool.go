package frame

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Request(ctx context.Context) *http.Request {
	return ctx.Value(gin.ContextRequestKey).(*http.Request)
}

func Header(ctx context.Context) http.Header {
	req := Request(ctx)
	return req.Header
}
