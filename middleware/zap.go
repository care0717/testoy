package middleware

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

func Zap(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Info("incoming request",
				zap.String("host", r.Host),
				zap.String("path", r.URL.Path),
				zap.String("ua", r.UserAgent()),
				zap.String("query", fmt.Sprintf("%v", r.URL.Query())),
			)
			next.ServeHTTP(w, r)
		})
	}
}
