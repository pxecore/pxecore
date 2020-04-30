package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pxecore/pxecore/pkg/entity"
	"github.com/pxecore/pxecore/pkg/errors"
	server "github.com/pxecore/pxecore/pkg/http"
	"github.com/pxecore/pxecore/pkg/repository"
	"net/http"
	"regexp"
)

var (
	hostIDRegex, _ = regexp.Compile("^[a-zA-Z0-9]+(?:[-_][a-zA-Z0-9]+)*$")
)

//~ STRUCT - Server -----------------------------------------------------------

// Host controller for the "/host" base path operations.
type Host struct {
	Repository repository.Repository // Repository dependency injection.
}

// Register implements http.Controller interface.
func (t Host) Register(r *mux.Router, config server.Config) {
	r.HandleFunc("/host/{id:[a-zA-Z0-9]+}", t.Get).Methods(http.MethodGet)
	r.HandleFunc("/host", t.Put).Methods(http.MethodPut)
}

// Get returns a template by ID.
func (t Host) Get(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	s, _ := v["id"]

	hb := new(HostBody)
	if err := t.Repository.Read(func(session repository.Session) error {
		t, err := session.Host().Get(s)
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

	w.WriteHeader(203)
}

// Put stores a new host.
func (t Host) Put(w http.ResponseWriter, r *http.Request) {
	tp := NewHostBody()
	if err := json.NewDecoder(r.Body).Decode(&tp); err != nil {
		server.WriteJSON(w, errors.MarshalJSON(err), http.StatusBadRequest)
		return
	}
	if err := tp.Validate(); err != nil {
		server.WriteJSON(w, errors.MarshalJSON(err), http.StatusBadRequest)
		return
	}

	if err := t.Repository.Write(func(session repository.Session) error {
		var err error
		if err = session.Host().Create(tp.ToEntity()); err != nil {
			if errors.Is(err, errors.ERepositoryKeyExist) {
				err = session.Host().Update(tp.ToEntity())
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
	server.WriteJSON(w, tp.JSON(), http.StatusCreated)
}

//~ STRUCT - JSON -----------------------------------------------------------

// HostBody stores host request and response data as well
// hold transformations and validations.
type HostBody struct {
	ID           string            `json:"id"`
	HardwareAddr []string          `json:"hardware-addr"`
	TrapMode     bool              `json:"trap-mode"`
	Vars         map[string]string `json:"vars"`
	GroupID      string            `json:"group-id"`
	TemplateID   string            `json:"template-id"`
}

// NewHostBody construct a new HostBody with default vars.
func NewHostBody() HostBody {
	return HostBody{
		ID:           "",
		HardwareAddr: make([]string, 0),
		TrapMode:     false,
		Vars:         make(map[string]string),
		GroupID:      "",
		TemplateID:   "",
	}
}

// Validate checks if the data hold in the instance follows the desired schema.
func (t HostBody) Validate() error {

	if !hostIDRegex.MatchString(t.ID) {
		return &errors.Error{
			Code: errors.EInvalidType,
			Msg:  fmt.Sprint("[controller.Host] id should follow pattern: ", hostIDRegex.String()),
		}
	}

	if len(t.HardwareAddr) == 0 {
		return &errors.Error{
			Code: errors.EInvalidType,
			Msg:  "[controller.Host] HardwareAddr should not be empty. ",
		}
	}

	return nil
}

// ToEntity returns an entity from the provided request.
func (t HostBody) ToEntity() entity.Host {
	return entity.Host{
		ID:            t.ID,
		HardwareAddr:  t.HardwareAddr,
		TrapMode:      t.TrapMode,
		TrapTriggered: false,
		Vars:          t.Vars,
		GroupID:       t.GroupID,
		TemplateID:    t.TemplateID,
	}
}

// JSON returns a json representation of the structure.
func (t HostBody) JSON() []byte {
	j, _ := json.Marshal(t)
	return j
}

// LoadEntity adds Entity vars to the HostBody
func (t *HostBody) LoadEntity(h entity.Host) {
	t.ID = h.ID
	t.GroupID = h.GroupID
	t.TemplateID = h.TemplateID
	t.Vars = h.Vars
	t.TrapMode = h.TrapMode
	t.HardwareAddr = h.HardwareAddr
}
