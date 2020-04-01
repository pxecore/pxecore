package memory

import "github.com/pxecore/pxecore/pkg/repository"

// Repository defines how to retrieve all particular entity repositories.
type Repository struct {
	hostRepository *repository.HostRepository
}

// NewRepository creates a new repository for the driver memory.
func NewRepository(config map[string]interface{}) (*repository.Repository, error) {
	r := new(Repository)
	c, err := NewConfig(r, config)
	if err != nil {
		return nil, err
	}

	hr, err := NewHostRepository(r, c)
	if err != nil {
		return nil, err
	}

	r.hostRepository = hr
	var ri repository.Repository = r
	return &ri, nil
}

// Host return the Host repository instance.
func (r *Repository) Host() *repository.HostRepository {
	return r.hostRepository
}
