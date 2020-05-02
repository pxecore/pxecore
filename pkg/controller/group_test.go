package controller

import (
	"bytes"
	"github.com/gorilla/mux"
	server "github.com/pxecore/pxecore/pkg/http"
	"github.com/pxecore/pxecore/pkg/repository"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGroup(t *testing.T) {
	r, _ := repository.NewRepository(map[string]interface{}{"driver": "memory"})
	ro := mux.NewRouter()
	ss := Group{Repository: r}
	ss.Register(ro, server.Config{})
	tests := []struct {
		name           string
		method         string
		path           string
		contentType    string
		body           string
		wantStatusCode int
		wantResponse   string
	}{
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, bytes.NewBuffer([]byte(tt.body)))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add("Content-Type", tt.contentType)
			rr := httptest.NewRecorder()
			ro.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatusCode)
			}

			if body := rr.Body.String(); tt.wantResponse != "" && body != tt.wantResponse {
				t.Errorf("handler returned wrong body: got %v want %v",
					body, tt.wantResponse)
			}
		})
	}
}
