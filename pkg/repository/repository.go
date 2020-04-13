package repository

import (
	"errors"
	"github.com/pxecore/pxecore/pkg/entity"
	"strings"
)

// Repository defines how to retrieve all particular entity repositories.
//
// Open(false) starts an read-only operation while Open(true) starts
// a read-write operation. Both operations need to be closed to prevent
// deadlocks.
//
// Read() runs a function without locking the repository access.
// Post execution automatically closes the session.
//
// Read() runs a function locking the repository access.
// Post execution automatically closes the session.
type Repository interface {
	Open(write bool) (Session, error)
	Read(func(session Session) error) error
	Write(func(session Session) error) error
}

// Session represents an repository read write operation.
//
// Close() terminates the session unlocking the repository access.
//
// IsReadOnly() returns true if the session only allows read-only operations.
//
// IsOpen() returns true if the session has not been closed.
//
// Host() returns entity.Host repository.
//
// Template() returns entity.Template repository.
type Session interface {
	Close() error
	IsReadOnly() bool
	IsOpen() bool
	Host() HostRepository
	Template() TemplateRepository
	Group() GroupRepository
}

// HostRepository defines the CRUD procedure for entity.Host
//
// Create() adds a new entity.Host into the repository or returns error
// errors.ERepositoryEmptyKey if the key is not provided,
// errors.ERepositoryKeyExist if the key or HardwareAddr already exists in the repository.
//
// Get() searches a entity.Host into by id or returns error
// errors.ERepositoryKeyNotFound if the key is not found.
//
// FindByHardwareAddr() searches a entity.Host by HardwareAddr or returns error
// errors.ERepositoryKeyNotFound if the HardwareAddr is not found.
//
// Update() update an existing entity.Host or returns error
// errors.ERepositoryEmptyKey if the key is not provided,
// errors.ERepositoryKeyNotFound if the key is not found,
// errors.ERepositoryKeyExist if the HardwareAddr already exists in the repository.
//
// Delete() deletes an entry of entity.Host or returns error
// errors.ERepositoryKeyNotFound if the HardwareAddr is not found.
type HostRepository interface {
	Create(host entity.Host) error
	Get(ID string) (entity.Host, error)
	FindByHardwareAddr(hardwareAddr string) (entity.Host, error)
	Update(host entity.Host) error
	Delete(host entity.Host) error
}

// GroupRepository defines the CRUD procedure for entity.Group
//
// Create() adds a new entity.Group into the repository or returns error
// errors.ERepositoryEmptyKey if the key is not provided,
// errors.ERepositoryKeyExist if the key already exists in the repository.
//
// Get() searches a entity.Group into by id or returns error
// errors.ERepositoryKeyNotFound if the key is not found.
//
// Update() update an existing entity.Group or returns error
// errors.ERepositoryEmptyKey if the key is not provided,
// errors.ERepositoryKeyNotFound if the key is not found
//
// Delete() deletes an entry of entity.Group or returns error
// errors.ERepositoryKeyNotFound if the HardwareAddr is not found.
type GroupRepository interface {
	Create(host entity.Group) error
	Get(ID string) (entity.Group, error)
	Update(host entity.Group) error
	Delete(host entity.Group) error
}

// TemplateRepository defines the CRUD procedure for entity.Template
//
// Create() adds a new entity.Template into the repository or returns error
// errors.ERepositoryEmptyKey if the key is not provided,
// errors.ERepositoryKeyExist if the key or HardwareAddr already exists in the repository.
//
// Get() searches a entity.Template into by id or returns error
// errors.ERepositoryKeyNotFound if the key is not found.
//
// Update() update an existing entity.Template or returns error
// errors.ERepositoryEmptyKey if the key is not provided,
// errors.ERepositoryKeyNotFound if the key is not found,
// errors.ERepositoryKeyExist if the HardwareAddr already exists in the repository.
//
// Delete() deletes an entry of entity.Template or returns error
// errors.ERepositoryKeyNotFound if the HardwareAddr is not found.
type TemplateRepository interface {
	Create(host entity.Template) error
	Get(ID string) (entity.Template, error)
	Update(host entity.Template) error
	Delete(host entity.Template) error
}

// NewRepository instantiates a new repository.
// Based on the "driver" key a different repository is created and
// passed the configuration.
func NewRepository(config map[string]interface{}) (Repository, error) {
	if val, ok := config["driver"]; ok {
		if driver, ok := val.(string); ok {
			switch strings.ToLower(driver) {
			case "memory":
				return newMemoryRepository(config)
			}
		}
		return nil, errors.New("invalid type in repository type")
	}
	return nil, errors.New("missing repository type")
}
