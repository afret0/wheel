package frame

import (
	"context"
	"net/http"

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
