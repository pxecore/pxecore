package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pxecore/pxecore/pkg/entity"
	server "github.com/pxecore/pxecore/pkg/http"
	"github.com/pxecore/pxecore/pkg/repository"
	"net/http"
	"regexp"
)

var (
	groupIDRegex, _ = regexp.Compile("^[a-zA-Z0-9]+(?:[-_][a-zA-Z0-9]+)*$")
)

//~ STRUCT - Server -----------------------------------------------------------

// Group controller for the "/group" base path operations.
type Group struct {
	Repository repository.Repository // Repository dependency injection.
}

// Register implements http.Controller interface.
func (t Group) Register(r *mux.Router, config server.Config) {
	r.HandleFunc("/group/{id:[a-zA-Z0-9]+}", t.Get).Methods(http.MethodGet)
	r.HandleFunc("/group", t.Post).Methods(http.MethodPost)
}

// Get returns a template by ID.
func (t Group) Get(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(203)
}

// Post stores a new host.
func (t Group) Post(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)
}

//~ STRUCT - JSON -----------------------------------------------------------

// GroupBody stores group request and response data as well
// hold transformations and validations.
type GroupBody struct {
	ID                string `json:"id"`
	HardwareAddr      []string
	TrapMode          bool
	Vars              map[string]string
	GroupID           string
	DefaultTemplateID string
}

// LoadTemplate fills the template with an entity values.
func (t *GroupBody) LoadTemplate(e entity.Template) {
	t.ID = e.ID
}

// Validate checks if the data hold in the instance follows the desired schema.
func (t GroupBody) Validate() error {
	return nil
}

// ToEntity returns an entity from the provided request.
func (t GroupBody) ToEntity() entity.Host {
	return entity.Host{
		ID: t.ID,
	}
}

// JSON returns a json representation of the structure.
func (t GroupBody) JSON() []byte {
	j, _ := json.Marshal(t)
	return j
}
