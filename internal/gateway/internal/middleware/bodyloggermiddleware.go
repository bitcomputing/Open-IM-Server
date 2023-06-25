package middleware

import (
	"Open_IM/pkg/logger"
	"net/http"
)

type BodyLoggerMiddleware struct {
}

func NewBodyLoggerMiddleware() *BodyLoggerMiddleware {
	return &BodyLoggerMiddleware{}
}

func (m *BodyLoggerMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.HandleRequest(r)
		next(w, r)
	}
}
