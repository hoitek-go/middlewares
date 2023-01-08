package bodylimit

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMiddlewareWhenBodyLimitHasCorrectValid(t *testing.T) {
	handler := Middleware(2 * 1024)
	handlerFunc := handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"statusCode":200,"data":"data"}`))
	}))
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{"statusCode":200,"data":"data"}`))
	handlerFunc.ServeHTTP(w, r)
	body := w.Result().Body
	defer body.Close()
	bytes, err := io.ReadAll(body)
	if err != nil {
		t.Error(err)
	}
	var bodyData map[string]interface{}
	err = json.Unmarshal(bytes, &bodyData)
	if err != nil {
		t.Error(err)
	}
	statusCodeField, ok := bodyData["statusCode"]
	if !ok {
		t.Error("Status code is not found in response body")
	}
	dataField, ok := bodyData["data"]
	if !ok {
		t.Error("Data is not present in response body")
	}
	if statusCodeField.(float64) != float64(http.StatusOK) {
		t.Error("Status code is invalid")
	}
	if dataField.(string) != "data" {
		t.Error("Data is invalid")
	}
}

func TestMiddlewareWhenBodyLimitHasNotCorrectValid(t *testing.T) {
	handler := Middleware(2 * 1024)
	handlerFunc := handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"statusCode":200,"data":"data"}`))
	}))
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(`{"statusCode":200,"data":"data"}`))
	r.ContentLength = 5 * 1024
	handlerFunc.ServeHTTP(w, r)
	body := w.Result().Body
	defer body.Close()
	bytes, err := io.ReadAll(body)
	if err != nil {
		t.Error(err)
	}
	var bodyData map[string]interface{}
	err = json.Unmarshal(bytes, &bodyData)
	if err != nil {
		t.Error(err)
	}
	statusCodeField, ok := bodyData["statusCode"]
	if !ok {
		t.Error("Status code is not found in response body")
	}
	dataField, ok := bodyData["data"]
	if !ok {
		t.Error("Data is not present in response body")
	}
	messageField, ok := bodyData["message"]
	if !ok {
		t.Error("Message is not present in response body")
	}
	if statusCodeField.(float64) != float64(http.StatusRequestEntityTooLarge) {
		t.Error("Status code is invalid")
	}
	if dataField.(string) != "" {
		t.Error("Data is invalid")
	}
	if messageField.(string) != "Request Entity Too Large" {
		t.Error("Data is invalid")
	}
}
