package recover

import (
	"log"
	"net/http"

	"github.com/hoitek-go/kit/response"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				log.Println(err)
				serverError := response.ErrorInternalServerError("")
				response.ErrorWithWriter(w, serverError, serverError.StatusCode)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
