package middleware

import (
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func RequestLogger(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()
			defer func() {
				l.Info("request",
					zap.String("proto", r.Proto),
					zap.String("path", r.URL.Path),
					zap.Duration("latency", time.Since(start)),
					zap.String("method", r.Method),
					zap.Int("status", ww.Status()),
					zap.Int("size", ww.BytesWritten()),
					zap.String("ip", r.RemoteAddr),
					zap.String("reqId", middleware.GetReqID(r.Context())))
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
