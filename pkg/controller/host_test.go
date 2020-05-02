package controller

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/pxecore/pxecore/pkg/entity"
	server "github.com/pxecore/pxecore/pkg/http"
	"github.com/pxecore/pxecore/pkg/repository"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHost(t *testing.T) {
	r, _ := repository.NewRepository(map[string]interface{}{"driver": "memory"})
	_ = r.Write(func(session repository.Session) error {
		_ = session.Group().Create(entity.Group{ID: "group1"})
		_ = session.Template().Create(entity.Template{ID: "template1"})
		return nil
	})
	ro := mux.NewRouter()
	ss := Host{Repository: r}
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
		{"OK_CREATE", http.MethodPut, "/host",
			"application/json",
			"{\"id\": \"host1\",\"hardware-addr\":[\"00-14-22-04-25-37\",\"00-14-22-04-25-38\"]," +
				"\"trap-mode\":true,\"vars\":{\"foo\":\"bar\"},\"group-id\":\"group1\",\"template-id\":\"template1\"}",
			http.StatusCreated, ""},
		{"KO_NOT_FOUND", http.MethodGet, "/host/id2",
			"application/json", "",
			http.StatusNotFound, ""},
		{"OK_FOUND", http.MethodGet, "/host/host1",
			"application/json", "",
			http.StatusOK, "{\"id\":\"host1\",\"hardware-addr\":[\"00-14-22-04-25-37\",\"00-14-22-04-25-38\"]," +
			"\"trap-mode\":true,\"vars\":{\"foo\":\"bar\"},\"group-id\":\"group1\",\"template-id\":\"template1\"}"},
		{"OK_CHANGE_ONLY_VAR", http.MethodPut, "/host",
			"application/json",
			"{\"id\": \"host1\",\"hardware-addr\":[\"00-14-22-04-25-37\",\"00-14-22-04-25-38\"]," +
				"\"trap-mode\":true,\"vars\":{\"foo\":\"bar1\"},\"group-id\":\"group1\",\"template-id\":\"template1\"}",
			http.StatusCreated, ""},
		{"OK_FOUND", http.MethodGet, "/host/host1",
			"application/json", "",
			http.StatusOK, "{\"id\":\"host1\",\"hardware-addr\":[\"00-14-22-04-25-37\",\"00-14-22-04-25-38\"]," +
			"\"trap-mode\":true,\"vars\":{\"foo\":\"bar1\"},\"group-id\":\"group1\",\"template-id\":\"template1\"}"},
		{"KO_MISSING_GROUP", http.MethodPut, "/host",
			"application/json",
			"{\"id\": \"host1\",\"hardware-addr\":[\"00-14-22-04-25-37\",\"00-14-22-04-25-38\"]," +
				"\"trap-mode\":true,\"vars\":{\"foo\":\"bar1\"},\"group-id\":\"group2\",\"template-id\":\"template1\"}",
			http.StatusFailedDependency, ""},
		{"KO_MISSING_TEMPLATE", http.MethodPut, "/host",
			"application/json",
			"{\"id\": \"host1\",\"hardware-addr\":[\"00-14-22-04-25-37\",\"00-14-22-04-25-38\"]," +
				"\"trap-mode\":true,\"vars\":{\"foo\":\"bar1\"},\"group-id\":\"group1\",\"template-id\":\"template2\"}",
			http.StatusFailedDependency, ""},
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
