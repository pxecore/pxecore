package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pxecore/pxecore/pkg/entity"
	"github.com/pxecore/pxecore/pkg/errors"
	server "github.com/pxecore/pxecore/pkg/http"
	"github.com/pxecore/pxecore/pkg/repository"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var (
	idRegex, _ = regexp.Compile("^[a-zA-Z0-9]+(?:[-_][a-zA-Z0-9]+)*$")
)

//~ STRUCT - Server -----------------------------------------------------------

// Template controller for the "/template" base path operations.
type Template struct {
	Repository repository.Repository // Repository dependency injection.
}

// Register implements http.Controller interface.
func (t Template) Register(r *mux.Router, config server.Config) {
	r.HandleFunc("/template/{id:[a-zA-Z0-9]+}", t.Get).Methods(http.MethodGet)
	r.HandleFunc("/template/{id:[a-zA-Z0-9]+}/template", t.GetTemplate).Methods(http.MethodGet)
	r.HandleFunc("/template/{id:[a-zA-Z0-9]+}/template", t.PostFile).Methods(http.MethodPost)
	r.HandleFunc("/template", t.Post).Methods(http.MethodPost)
}

// Get returns a template by ID.
func (t Template) Get(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	s, ok := v["id"]
	if !ok {
		server.WriteJSON(w, errors.MarshalJSON(
			&errors.Error{Code: errors.EUnknown, Msg: "missing url param id"}),
			http.StatusBadRequest)
	}
	var tb TemplateBody
	if err := t.Repository.Read(func(session repository.Session) error {
		t, err := session.Template().Get(s)
		if err != nil {
			return err
		}
		tb.LoadTemplate(t)
		return nil
	}); err != nil {
		if errors.Is(err, errors.ERepositoryKeyNotFound) {
			server.WriteJSON(w, errors.MarshalJSON(err), http.StatusNoContent)
		} else {
			server.WriteJSON(w, errors.MarshalJSON(err), http.StatusInternalServerError)
		}
	} else {
		server.WriteJSON(w, tb.JSON(), http.StatusOK)
	}
}

// GetTemplate returns a template by ID.
func (t Template) GetTemplate(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	s, ok := v["id"]
	if !ok {
		err := errors.Error{Code: errors.EUnknown, Msg: "missing url param id"}
		server.WriteText(w,
			err.Error(),
			http.StatusBadRequest)
	}
	var tb TemplateBody
	if err := t.Repository.Read(func(session repository.Session) error {
		t, err := session.Template().Get(s)
		if err != nil {
			return err
		}
		tb.LoadTemplate(t)
		return nil
	}); err != nil {
		if errors.Is(err, errors.ERepositoryKeyNotFound) {
			server.WriteText(w, err.Error(), http.StatusNoContent)
		} else {
			server.WriteText(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		server.WriteText(w, tb.Template, http.StatusOK)
	}
}

// Post saves a template
func (t Template) Post(w http.ResponseWriter, r *http.Request) {
	var tp TemplateBody
	err := json.NewDecoder(r.Body).Decode(&tp)
	if err != nil {
		server.WriteJSON(w, errors.MarshalJSON(err), http.StatusBadRequest)
		return
	}
	if err := tp.Validate(); err != nil {
		server.WriteJSON(w, errors.MarshalJSON(err), http.StatusBadRequest)
		return
	}
	if err := t.Repository.Write(func(session repository.Session) error {
		if err := session.Template().Create(tp.ToEntity()); err != nil {
			if errors.Is(err, errors.ERepositoryKeyExist) {
				if err := session.Template().Update(tp.ToEntity()); err != nil {
					return err
				}
				return nil
			}
			return err
		}
		return nil
	}); err != nil {
		server.WriteJSON(w, errors.MarshalJSON(err), http.StatusInternalServerError)
	}
}

// PostFile saves a template by reading the ID from the URL and the file from the body.
func (t Template) PostFile(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	s, ok := v["id"]
	if !ok {
		server.WriteJSON(w, errors.MarshalJSON(
			&errors.Error{Code: errors.EUnknown, Msg: "missing url param id"}),
			http.StatusBadRequest)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		server.WriteJSON(w, errors.MarshalJSON(err), http.StatusBadRequest)
		return
	}
	tp := TemplateBody{
		ID:       s,
		Template: string(body),
	}
	if err := tp.Validate(); err != nil {
		server.WriteJSON(w, errors.MarshalJSON(err), http.StatusBadRequest)
		return
	}
	if err := t.Repository.Write(func(session repository.Session) error {
		if err := session.Template().Create(tp.ToEntity()); err != nil {
			if errors.Is(err, errors.ERepositoryKeyExist) {
				if err := session.Template().Update(tp.ToEntity()); err != nil {
					return err
				}
				return nil
			}
			return err
		}
		return nil
	}); err != nil {
		server.WriteJSON(w, errors.MarshalJSON(err), http.StatusInternalServerError)
	}
}

//~ STRUCT - JSON -----------------------------------------------------------

// TemplateBody stores template request and response data as well
// hold transformations and validations.
type TemplateBody struct {
	ID       string `json:"id"`
	Template string `json:"template"`
}

// LoadTemplate fills the template with an entity values.
func (t *TemplateBody) LoadTemplate(e entity.Template) {
	t.ID = e.ID
	t.Template = e.Template
}

// Validate checks if the data hold in the instance follows the desired schema.
func (t TemplateBody) Validate() error {
	if !idRegex.MatchString(t.ID) {
		return &errors.Error{
			Code: errors.EInvalidType,
			Msg:  fmt.Sprint("[controller.Template] id should follow pattern: ", idRegex.String()),
		}
	}
	if strings.TrimSpace(t.Template) == "" {
		return &errors.Error{
			Code: errors.EInvalidType,
			Msg:  "[controller.Template] template should not be empty.",
		}
	}
	return nil
}

// ToEntity returns an entity from the provided request.
func (t TemplateBody) ToEntity() entity.Template {
	return entity.Template{
		ID:       t.ID,
		Template: t.Template,
	}
}

// JSON returns a json representation of the structure.
func (t TemplateBody) JSON() []byte {
	j, _ := json.Marshal(t)
	return j
}
