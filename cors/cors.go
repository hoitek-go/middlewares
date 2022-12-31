package cors

import (
	"net/http"

	"github.com/gorilla/handlers"
)

func Cors(next http.Handler) http.Handler {
	cors := handlers.CORS(
		handlers.AllowedOrigins(
			[]string{
				"http://localhost:3000",
			},
		),
	)
	return cors(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}))
}
