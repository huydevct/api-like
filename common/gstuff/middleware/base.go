package middleware

import (
	"app/common/config"
	"io"
	"net/http"
)

var cfg = config.GetConfig()

type bodyDumpResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *bodyDumpResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
