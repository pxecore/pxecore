package http

import (
	"github.com/gorilla/mux"
	"github.com/pxecore/pxecore/pkg/errors"
	"github.com/pxecore/pxecore/pkg/util"
	"net/http"
	"time"
)

//~ INTERFACE - Controller ----------------------------------------------------

// Controller interface.
type Controller interface {
	Register(r *mux.Router, config Config)
}

//~ STRUCT - Server -----------------------------------------------------------

// Server manages all http interaction.
type Server struct {
	Controllers []Controller
	server      *http.Server
	router      *mux.Router
}

// Start initiates the server.
func (s *Server) Start(config Config) error {
	if s.server != nil {
		return &errors.Error{Code: errors.EAlreadyRunning, Msg: "server already running"}
	}

	s.router = mux.NewRouter()
	for _, c := range s.Controllers {
		c.Register(s.router, config)
	}

	s.server = &http.Server{
		Handler:      s.router,
		Addr:         config.Address,
		WriteTimeout: config.ReadTimeout,
		ReadTimeout:  config.WriteTimeout,
	}
	return s.server.ListenAndServe()
}

//~ STRUCT - Config -----------------------------------------------------------

// Config converts and stores the server config.
type Config struct {
	Address      string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

// NewConfig populates the Config with values.
func NewConfig(config map[string]interface{}) (Config, error) {
	c := Config{}

	s, err := util.StringFromMap(config, "address", ":80")
	if err != nil {
		return c, &errors.Error{Code: errors.Code(err), Msg: "[http.Config] error reading config."}
	}
	c.Address = s

	i, err := util.IntFromMap(config, "write-timeout", 10)
	if err != nil {
		return c, &errors.Error{Code: errors.Code(err), Msg: "[http.Config] error reading config."}
	}
	c.WriteTimeout = time.Duration(i) * time.Second

	i, err = util.IntFromMap(config, "read-timeout", 10)
	if err != nil {
		return c, &errors.Error{Code: errors.Code(err), Msg: "[http.Config] error reading config."}
	}
	c.ReadTimeout = time.Duration(i) * time.Second

	return c, nil
}

// WriteJSON writes a json in byte array as response.
func WriteJSON(w http.ResponseWriter, j []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	w.Write(j)
}

// WriteText writes a json in byte array as response.
func WriteText(w http.ResponseWriter, j string, code int) {
	w.Header().Set("Content-Type", "application/text; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	w.Write([]byte(j))
}
