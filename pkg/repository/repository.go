package repository

import (
	"errors"
	"github.com/pxecore/pxecore/pkg/entity"
	"strings"
)

// Repository defines how to retrieve all particular entity repositories.
type Repository interface {
	Host() *HostRepository
}

// HostRepository defines the CRUD procedure for entity.Host
type HostRepository interface {
	Create(host entity.Host) error
	Get(ID string) (entity.Host, error)
	FindByHardwareAddr(hardwareAddr string) (entity.Host, error)
	Update(host entity.Host) error
	Delete(host entity.Host) error
}

// NewRepository instantiates a new repository.
// Based on the "driver" key a different repository is created and
// passed the configuration.
func NewRepository(config map[string]interface{}) (*Repository, error) {
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
