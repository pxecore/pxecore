package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pxecore/pxecore/pkg/entity"
	"github.com/pxecore/pxecore/pkg/errors"
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
	r.HandleFunc("/group", t.Put).Methods(http.MethodPut)
}

// Get returns a template by ID.
func (t Group) Get(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	s, _ := v["id"]

	hb := NewGroupBody()
	if err := t.Repository.Read(func(session repository.Session) error {
		t, err := session.Group().Get(s)
		if err != nil {
			return err
		}
		hb.LoadEntity(t)
		return nil
	}); err != nil {
		if errors.Is(err, errors.ERepositoryKeyNotFound) {
			server.WriteText(w, err.Error(), http.StatusNotFound)
		} else {
			server.WriteText(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		server.WriteJSON(w, hb.JSON(), http.StatusOK)
	}
}

// Put stores a new host.
func (t Group) Put(w http.ResponseWriter, r *http.Request) {
	body := NewGroupBody()
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		server.WriteJSON(w, errors.MarshalJSON(err), http.StatusBadRequest)
		return
	}
	if err := body.Validate(); err != nil {
		server.WriteJSON(w, errors.MarshalJSON(err), http.StatusBadRequest)
		return
	}

	if err := t.Repository.Write(func(session repository.Session) error {
		var err error
		if err = session.Group().Create(body.ToEntity()); err != nil {
			if errors.Is(err, errors.ERepositoryKeyExist) {
				err = session.Group().Update(body.ToEntity())
			}
		}
		return err
	}); err != nil {
		if errors.Is(err, errors.ERepositoryKeyNotFound) {
			server.WriteJSON(w, errors.MarshalJSON(err), http.StatusFailedDependency)
		} else {
			server.WriteJSON(w, errors.MarshalJSON(err), http.StatusInternalServerError)
		}
	}
	server.WriteJSON(w, []byte{}, http.StatusCreated)
}

//~ STRUCT - JSON -----------------------------------------------------------

// GroupBody stores group request and response data as well
// hold transformations and validations.
type GroupBody struct {
	ID         string            `json:"id"`
	Vars       map[string]string `json:"vars"`
	ParentID   string            `json:"parent-id"`
	TemplateID string            `json:"template-id"`
	HostsIDs   []string          `json:"hosts"`
	GroupIDs   []string          `json:"groups"`
}

// NewGroupBody constructs a new GroupBody
func NewGroupBody() GroupBody {
	return GroupBody{
		ID:         "",
		Vars:       make(map[string]string),
		ParentID:   "",
		TemplateID: "",
		HostsIDs:   make([]string, 0),
		GroupIDs:   make([]string, 0),
	}
}

// LoadEntity fills the template with an entity values.
func (t *GroupBody) LoadEntity(e entity.Group) {
	t.ID = e.ID
	t.Vars = e.Vars
	t.TemplateID = e.TemplateID
	t.ParentID = e.ParentID
	t.GroupIDs = e.GroupIDs
	t.HostsIDs = e.HostsIDs

}

// Validate checks if the data hold in the instance follows the desired schema.
func (t GroupBody) Validate() error {
	return nil
}

// ToEntity returns an entity from the provided request.
func (t GroupBody) ToEntity() entity.Group {
	return entity.Group{
		ID:         t.ID,
		Vars:       t.Vars,
		ParentID:   t.ParentID,
		TemplateID: t.TemplateID,
	}
}

// JSON returns a json representation of the structure.
func (t GroupBody) JSON() []byte {
	j, _ := json.Marshal(t)
	return j
}
