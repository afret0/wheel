package tool

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Uid Deprecated  use frame.Uid instead
func Uid(ctx context.Context) string {

	req := ctx.Value(gin.ContextRequestKey)
	if req == nil {
		return ""
	}

	req1, ok := req.(*http.Request)
	if !ok {
		return ""
	}

	header := req1.Header

	//header := frame.Header(ctx)
	uid := header.Get("_uid")
	if uid == "" {
		panic("uid no exist")
	}

	//uid, ok := ctx.Value("_uid").(string)
	//if !ok {
	//	panic("uid not found in context")
	//}
	//if uid == "" {
	//	panic("uid is empty")
	//}

	return uid
}
