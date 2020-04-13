package repository

import (
	"github.com/pxecore/pxecore/pkg/entity"
	"github.com/pxecore/pxecore/pkg/errors"
	"sync"
)

//~ STRUCT - memoryRepository -------------------------------------------------

// Repository defines how to retrieve all particular entity repositories.
type memoryRepository struct {
	lock              *sync.RWMutex
	config            MemoryConfig
	hosts             map[string]*entity.Host
	hardwareAddrIndex map[string]*entity.Host
	groups            map[string]*entity.Group
	templates         map[string]*entity.Template
}

func (m *memoryRepository) Open(write bool) (Session, error) {
	if write {
		m.lock.Lock()
	} else {
		m.lock.RLock()
	}
	return newMemorySession(m, m.config, !write), nil
}

func (m *memoryRepository) Read(f func(session Session) error) error {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return f(newMemorySession(m, m.config, true))
}

func (m *memoryRepository) Write(f func(session Session) error) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	return f(newMemorySession(m, m.config, false))
}

// NewRepository creates a new repository for the driver memory.
func newMemoryRepository(config map[string]interface{}) (Repository, error) {
	r := new(memoryRepository)
	var ri Repository = r

	c, err := NewConfig(r, config)
	if err != nil {
		return nil, err
	}
	r.config = c
	r.lock = new(sync.RWMutex)
	r.hosts = make(map[string]*entity.Host)
	r.hardwareAddrIndex = make(map[string]*entity.Host)
	r.groups = make(map[string]*entity.Group)
	r.templates = make(map[string]*entity.Template)
	return ri, nil
}

//~ STRUCT - memorySession ------------------------------------------------

// MemorySession holds a single use of the repository.
type MemorySession struct {
	lock               *sync.RWMutex
	open               bool
	readOnly           bool
	repository         *memoryRepository
	config             MemoryConfig
	hostRepository     *HostRepository
	templateRepository *TemplateRepository
	groupRepository    *GroupRepository
}

// Close terminates the session
func (m *MemorySession) Close() error {
	if m.open {
		if m.readOnly {
			m.repository.lock.RUnlock()
		} else {
			m.repository.lock.Unlock()
		}
		m.open = false
		return nil
	}
	return &errors.Error{Code: errors.EUnknown, Msg: "repository already closed"}
}

// Host returns HostRepository
func (m *MemorySession) Host() HostRepository {
	if m.hostRepository == nil {
		m.hostRepository = newMemoryHostRepository(
			m,
			m.config,
			m.repository.hosts,
			m.repository.hardwareAddrIndex)
	}
	return *m.hostRepository
}

// Template returns TemplateRepository
func (m *MemorySession) Template() TemplateRepository {
	if m.templateRepository == nil {
		m.templateRepository = newMemoryTemplateRepository(m, m.config, m.repository.templates)
	}
	return *m.templateRepository
}

// Group returns GroupRepository
func (m *MemorySession) Group() GroupRepository {
	if m.groupRepository == nil {
		m.groupRepository = newMemoryGroupRepository(
			m,
			m.config,
			m.repository.groups)
	}
	return *m.groupRepository
}

// IsReadOnly returns true is the session is for read only.
func (m *MemorySession) IsReadOnly() bool {
	return m.readOnly
}

// IsOpen returns true is the session is open for transactions.
func (m *MemorySession) IsOpen() bool {
	return m.readOnly
}

func newMemorySession(r *memoryRepository, config MemoryConfig, readOnly bool) Session {
	m := MemorySession{
		repository: r,
		open:       true,
		readOnly:   readOnly,
		config:     config,
	}
	return &m
}

//~ STRUCT - MemoryConfig -------------------------------------------------

// MemoryConfig stores memory driver config for all repositories.
type MemoryConfig struct {
	allowReset bool
}

// NewConfig creates a new Config extracting and checking type of the required fields.
func NewConfig(r *memoryRepository, config map[string]interface{}) (MemoryConfig, error) {
	c := MemoryConfig{}

	c.allowReset = false
	if e, ok := config["allow-reset"]; ok {
		val, ok := e.(bool)
		if !ok {
			return c, &errors.Error{Code: errors.EInvalidType, Msg: "config invalid type for key allow-reset"}
		}
		c.allowReset = val
	}

	return c, nil
}
