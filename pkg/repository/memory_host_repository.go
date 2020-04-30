package repository

import (
	"fmt"
	"github.com/pxecore/pxecore/pkg/entity"
	"github.com/pxecore/pxecore/pkg/errors"
)

// HostRepository defines the CRUD procedure for entity.Host
type memoryHostRepository struct {
	session           Session
	config            MemoryConfig
	hosts             map[string]*entity.Host
	hardwareAddrIndex map[string]*entity.Host
}

// NewHostRepository instantiates a new repository for entity.Host
func newMemoryHostRepository(
	s Session,
	config MemoryConfig,
	hosts map[string]*entity.Host,
	hardwareAddrIndex map[string]*entity.Host) *HostRepository {
	var hr HostRepository
	hr = &memoryHostRepository{
		s,
		config,
		hosts,
		hardwareAddrIndex,
	}
	return &hr
}

// Create add a new entity.Host to the repository
func (h *memoryHostRepository) Create(host entity.Host) error {
	if h.session.IsReadOnly() {
		return &errors.Error{Code: errors.ERepositoryReadOnly, Msg: "read-only mode"}
	}
	e := host
	if e.ID == "" {
		return &errors.Error{Code: errors.ERepositoryEmptyKey,
			Msg: "entity.Host key is empty"}
	}
	if _, ok := h.hosts[e.ID]; ok {
		return &errors.Error{Code: errors.ERepositoryKeyExist,
			Msg: fmt.Sprintf("entity.Host key %v already exists ", e.ID)}
	}
	for _, e := range e.HardwareAddr {
		if _, ok := h.hardwareAddrIndex[e]; ok {
			return &errors.Error{Code: errors.ERepositoryKeyExist,
				Msg: fmt.Sprintf("entity.Host HardwareAddr %v already exists ", e)}
		}
	}

	var err error
	if e.TemplateID != "" {
		if _, err = h.session.Template().Get(e.TemplateID); err != nil {
			return &errors.Error{Code: errors.ERepositoryKeyNotFound,
				Msg: fmt.Sprintf("entity.Host TemplateID %v not found.", e)}
		}
	}
	if e.GroupID != "" {
		g, err := h.session.Group().Get(e.GroupID)
		if g, err = h.session.Group().Get(e.GroupID); err != nil {
			return &errors.Error{Code: errors.ERepositoryKeyNotFound,
				Msg: fmt.Sprintf("entity.Host GroupID %v not found.", e),
				Err: err}
		}
		g.AddHost(e.ID)
		if err := h.session.Group().Update(g); err != nil {
			return &errors.Error{Code: errors.EUnknown,
				Msg: fmt.Sprintf("entity.Host GroupID %v can't be updated.", e),
				Err: err}
		}
	}

	h.hosts[e.ID] = &e
	for _, m := range e.HardwareAddr {
		h.hardwareAddrIndex[m] = &e
	}
	return nil
}

// Get implements repository.HostRepository interface
func (h *memoryHostRepository) Get(ID string) (entity.Host, error) {
	if val, ok := h.hosts[ID]; ok {
		return *val, nil
	}
	return entity.Host{}, &errors.Error{Code: errors.ERepositoryKeyNotFound,
		Msg: fmt.Sprintf("entity.Host key %v not found", ID)}
}

// FindByHardwareAddr implements repository.HostRepository interface
func (h *memoryHostRepository) FindByHardwareAddr(hardwareAddr string) (entity.Host, error) {
	if val, ok := h.hardwareAddrIndex[hardwareAddr]; ok {
		return *val, nil
	}
	return entity.Host{}, &errors.Error{Code: errors.ERepositoryKeyNotFound,
		Msg: fmt.Sprintf("entity.Host key %v not found", hardwareAddr)}
}

// Update implements repository.HostRepository interface
func (h *memoryHostRepository) Update(host entity.Host) error {
	if h.session.IsReadOnly() {
		return &errors.Error{Code: errors.ERepositoryReadOnly, Msg: "read-only mode"}
	}
	e := host
	if e.ID == "" {
		return &errors.Error{Code: errors.ERepositoryEmptyKey,
			Msg: "entity.Host key is empty"}
	}
	oe, ok := h.hosts[e.ID]
	if !ok {
		return &errors.Error{Code: errors.ERepositoryKeyNotFound,
			Msg: fmt.Sprintf("entity.Host key %v not found ", e.ID)}
	}
	for _, ee := range e.HardwareAddr {
		if val, ok := h.hardwareAddrIndex[ee]; ok {
			if val.ID != e.ID {
				return &errors.Error{Code: errors.ERepositoryKeyExist,
					Msg: fmt.Sprintf("entity.Host HardwareAddr %v already exists ", ee)}
			}
		}
	}
	if e.TemplateID != "" {
		if _, err := h.session.Template().Get(e.TemplateID); err != nil {
			return &errors.Error{Code: errors.ERepositoryKeyNotFound,
				Msg: fmt.Sprintf("entity.Host TemplateID %v not found.", e)}
		}
	}

	if e.GroupID != "" {
		g, err := h.session.Group().Get(e.GroupID)
		if err != nil {
			return &errors.Error{Code: errors.ERepositoryKeyNotFound,
				Msg: fmt.Sprintf("entity.Host GroupID %v not found.", e),
				Err: err}
		}
		g.AddHost(e.ID)
		if err := h.session.Group().Update(g); err != nil {
			return &errors.Error{Code: errors.EUnknown,
				Msg: fmt.Sprintf("entity.Host GroupID %v can't be updated.", e),
				Err: err}
		}
	}

	if oe.GroupID != "" {
		g, err := h.session.Group().Get(oe.GroupID)
		if err == nil {
			g.RemoveHost(e.ID)
			if err := h.session.Group().Update(g); err != nil {
				return &errors.Error{Code: errors.EUnknown,
					Msg: fmt.Sprintf("entity.Host GroupID %v can't be removed from the leaving group.", e),
					Err: err}
			}
		}
	}
	h.hosts[e.ID] = &e
	for _, val := range oe.HardwareAddr {
		delete(h.hardwareAddrIndex, val)
	}
	for _, m := range e.HardwareAddr {
		h.hardwareAddrIndex[m] = &e
	}
	return nil
}

// Delete implements repository.HostRepository interface
func (h *memoryHostRepository) Delete(host entity.Host) error {
	if h.session.IsReadOnly() {
		return &errors.Error{Code: errors.ERepositoryReadOnly, Msg: "read-only mode"}
	}
	e := host
	if e.ID == "" {
		return &errors.Error{Code: errors.ERepositoryEmptyKey,
			Msg: "entity.Host key is empty"}
	}
	oe, ok := h.hosts[e.ID]
	if !ok {
		return &errors.Error{Code: errors.ERepositoryKeyNotFound,
			Msg: fmt.Sprintf("entity.Host key %v not found ", e.ID)}
	}
	for _, e := range e.HardwareAddr {
		if _, ok := h.hardwareAddrIndex[e]; ok {
			return &errors.Error{Code: errors.ERepositoryKeyExist,
				Msg: fmt.Sprintf("entity.Host HardwareAddr %v already exists ", e)}
		}
	}
	for _, val := range oe.HardwareAddr {
		delete(h.hardwareAddrIndex, val)
	}
	if oe.GroupID != "" {
		if g, err := h.session.Group().Get(oe.GroupID); err == nil {
			g.RemoveHost(e.ID)
			if err := h.session.Group().Update(g); err != nil {
				return &errors.Error{Code: errors.EUnknown,
					Msg: fmt.Sprintf("entity.Host GroupID %v can't be removed from the leaving group.", e),
					Err: err}
			}
		}
	}
	delete(h.hosts, oe.ID)
	return nil
}
