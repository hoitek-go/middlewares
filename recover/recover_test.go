package recover

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware(t *testing.T) {
	t.Run("When panic is occurred", func(t *testing.T) {
		handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			panic("some error")
		}))
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		handler.ServeHTTP(w, r)
		if w.Code != http.StatusInternalServerError {
			t.Error("Status code is not 500")
		}
		defer w.Result().Body.Close()
		bytes, err := io.ReadAll(w.Result().Body)
		if err != nil {
			t.Error(err)
		}
		var bodyData map[string]interface{}
		err = json.Unmarshal(bytes, &bodyData)
		if err != nil {
			t.Error(err)
		}
		if bodyData["message"].(string) != "Something Went Wrong" {
			t.Error("Message is not Something Went Wrong")
		}
		if bodyData["statusCode"].(float64) != http.StatusInternalServerError {
			t.Error("Status code is not 500")
		}
		if bodyData["data"].(string) != "" {
			t.Error("Data is not empty")
		}
	})
}
