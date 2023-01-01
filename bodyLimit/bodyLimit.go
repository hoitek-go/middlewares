package bodylimit

import (
	"net/http"

	"github.com/hoitek-go/kit/response"
)

func Middleware(limit int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > limit {
				serverError := response.ErrorRequestEntityTooLarge("")
				response.ErrorWithWriter(w, serverError, serverError.StatusCode)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
