package controller

import (
	"github.com/gorilla/mux"
	server "github.com/pxecore/pxecore/pkg/http"
	"net/http"
)

//~ STRUCT - Server -----------------------------------------------------------

// Static controller for the "/static" base path operations.
type Static struct {
	BaseDir string
}

// Register implements http.Controller interface.
func (t Static) Register(r *mux.Router, config server.Config) {
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(t.BaseDir))))
}
