package middleware

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"CBCTF/internal/resp"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const multipartFormMemory = 32 << 20

// LimitUploadSize limits the total upload request body size.
func LimitUploadSize(maxBytes int64) gin.HandlerFunc {
	formatByteSize := func(size int64) string {
		const unit = 1024
		if size < unit {
			return fmt.Sprintf("%d B", size)
		}
		div, exp := int64(unit), 0
		for n := size / unit; n >= unit; n /= unit {
			div *= unit
			exp++
		}
		return fmt.Sprintf("%.1f %ciB", float64(size)/float64(div), "KMGTPE"[exp])
	}
	return func(ctx *gin.Context) {
		if maxBytes <= 0 || ctx.Request.Body == nil {
			ctx.Next()
			return
		}

		if ctx.Request.ContentLength > maxBytes {
			resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.RequestTooLarge, Attr: map[string]any{"Limit": formatByteSize(maxBytes)}})
			return
		}

		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxBytes)
		if strings.HasPrefix(strings.ToLower(ctx.GetHeader("Content-Type")), "multipart/form-data") {
			if err := ctx.Request.ParseMultipartForm(multipartFormMemory); err != nil {
				if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
					resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.RequestTooLarge, Attr: map[string]any{"Limit": formatByteSize(maxBytes)}})
					return
				}
				resp.AbortJSON(ctx, model.RetVal{Msg: i18n.Response.BadRequest, Attr: map[string]any{"Error": err.Error()}})
				return
			}
		}
		ctx.Next()
	}
}
