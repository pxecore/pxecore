package repository

import (
	"fmt"
	"github.com/pxecore/pxecore/pkg/entity"
	"github.com/pxecore/pxecore/pkg/errors"
)

// TemplateRepository defines the CRUD procedure for entity.Template
type memoryTemplateRepository struct {
	session   Session
	config    MemoryConfig
	templates map[string]*entity.Template
}

// NewTemplateRepository instantiates a new repository for entity.Template
func newMemoryTemplateRepository(s Session, config MemoryConfig, templates map[string]*entity.Template) *TemplateRepository {
	var hr TemplateRepository
	hr = &memoryTemplateRepository{
		s,
		config,
		templates,
	}
	return &hr
}

// Create add a new entity.Template to the repository
func (h *memoryTemplateRepository) Create(template entity.Template) error {
	if h.session.IsReadOnly() {
		return &errors.Error{Code: errors.ERepositoryReadOnly, Msg: "read-only mode"}
	}
	e := template
	if e.ID == "" {
		return &errors.Error{Code: errors.ERepositoryEmptyKey,
			Msg: "entity.Template key is empty"}
	}
	if _, ok := h.templates[e.ID]; ok {
		return &errors.Error{Code: errors.ERepositoryKeyExist,
			Msg: fmt.Sprintf("entity.Template key %v already exists ", e.ID)}
	}
	h.templates[e.ID] = &e
	return nil
}

// Get implements repository.TemplateRepository interface
func (h *memoryTemplateRepository) Get(ID string) (entity.Template, error) {
	if val, ok := h.templates[ID]; ok {
		return *val, nil
	}
	return entity.Template{}, &errors.Error{Code: errors.ERepositoryKeyNotFound,
		Msg: fmt.Sprintf("entity.Template key %v not found", ID)}
}

// Update implements repository.TemplateRepository interface
func (h *memoryTemplateRepository) Update(template entity.Template) error {
	if h.session.IsReadOnly() {
		return &errors.Error{Code: errors.ERepositoryReadOnly, Msg: "read-only mode"}
	}
	e := template
	if e.ID == "" {
		return &errors.Error{Code: errors.ERepositoryEmptyKey,
			Msg: "entity.Template key is empty"}
	}
	if _, ok := h.templates[e.ID]; !ok {
		return &errors.Error{Code: errors.ERepositoryKeyNotFound,
			Msg: fmt.Sprintf("entity.Template key %v not found ", e.ID)}
	}
	h.templates[e.ID] = &e
	return nil
}

// Delete implements repository.TemplateRepository interface
func (h *memoryTemplateRepository) Delete(template entity.Template) error {
	if h.session.IsReadOnly() {
		return &errors.Error{Code: errors.ERepositoryReadOnly, Msg: "read-only mode"}
	}
	e := template
	if e.ID == "" {
		return &errors.Error{Code: errors.ERepositoryEmptyKey,
			Msg: "entity.Template key is empty"}
	}
	oe, ok := h.templates[e.ID]
	if !ok {
		return &errors.Error{Code: errors.ERepositoryKeyNotFound,
			Msg: fmt.Sprintf("entity.Template key %v not found ", e.ID)}
	}
	delete(h.templates, oe.ID)
	return nil
}
