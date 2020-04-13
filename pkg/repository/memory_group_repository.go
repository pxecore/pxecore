package repository

import (
	"fmt"
	"github.com/pxecore/pxecore/pkg/entity"
	"github.com/pxecore/pxecore/pkg/errors"
)

// GroupRepository defines the CRUD procedure for entity.Group
type memoryGroupRepository struct {
	session Session
	config  MemoryConfig
	groups  map[string]*entity.Group
}

// NewGroupRepository instantiates a new repository for entity.Group
func newMemoryGroupRepository(s Session, config MemoryConfig, groups map[string]*entity.Group) *GroupRepository {
	var hr GroupRepository
	hr = &memoryGroupRepository{
		s,
		config,
		groups,
	}
	return &hr
}

// Create add a new entity.Group to the repository
func (h *memoryGroupRepository) Create(Group entity.Group) error {
	if h.session.IsReadOnly() {
		return &errors.Error{Code: errors.ERepositoryReadOnly, Msg: "read-only mode"}
	}
	e := Group
	if e.ID == "" {
		return &errors.Error{Code: errors.ERepositoryEmptyKey,
			Msg: "entity.Group key is empty"}
	}
	if _, ok := h.groups[e.ID]; ok {
		return &errors.Error{Code: errors.ERepositoryKeyExist,
			Msg: fmt.Sprintf("entity.Group key %v already exists ", e.ID)}
	}
	h.groups[e.ID] = &e
	return nil
}

// Get implements repository.GroupRepository interface
func (h *memoryGroupRepository) Get(ID string) (entity.Group, error) {
	if val, ok := h.groups[ID]; ok {
		return *val, nil
	}
	return entity.Group{}, &errors.Error{Code: errors.ERepositoryKeyNotFound,
		Msg: fmt.Sprintf("entity.Group key %v not found", ID)}
}

// Update implements repository.GroupRepository interface
func (h *memoryGroupRepository) Update(Group entity.Group) error {
	if h.session.IsReadOnly() {
		return &errors.Error{Code: errors.ERepositoryReadOnly, Msg: "read-only mode"}
	}
	e := Group
	if e.ID == "" {
		return &errors.Error{Code: errors.ERepositoryEmptyKey,
			Msg: "entity.Group key is empty"}
	}

	if _, ok := h.groups[e.ID]; !ok {
		return &errors.Error{Code: errors.ERepositoryKeyNotFound,
			Msg: fmt.Sprintf("entity.Group key %v not found ", e.ID)}
	}
	h.groups[e.ID] = &e
	return nil
}

// Delete implements repository.GroupRepository interface
func (h *memoryGroupRepository) Delete(Group entity.Group) error {
	if h.session.IsReadOnly() {
		return &errors.Error{Code: errors.ERepositoryReadOnly, Msg: "read-only mode"}
	}
	e := Group
	if e.ID == "" {
		return &errors.Error{Code: errors.ERepositoryEmptyKey,
			Msg: "entity.Group key is empty"}
	}
	oe, ok := h.groups[e.ID]
	if !ok {
		return &errors.Error{Code: errors.ERepositoryKeyNotFound,
			Msg: fmt.Sprintf("entity.Group key %v not found ", e.ID)}
	}
	delete(h.groups, oe.ID)
	return nil
}
