package json

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMiddleware(t *testing.T) {
	jsonData := `{"statusCode":200,"data":"data"}`
	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonData))
	}))
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(jsonData))
	handler.ServeHTTP(w, r)
	if w.Header().Get("Content-Type") != "application/json" {
		t.Error("Content-Type is invalid")
	}
}
