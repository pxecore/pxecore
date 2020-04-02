package repository

// Repository defines how to retrieve all particular entity repositories.
type memoryRepository struct {
	hostRepository *HostRepository
}

// NewRepository creates a new repository for the driver memory.
func newMemoryRepository(config map[string]interface{}) (*Repository, error) {
	r := new(memoryRepository)
	var ri Repository = r

	c, err := NewConfig(r, config)
	if err != nil {
		return nil, err
	}

	hr, err := newMemoryHostRepository(&ri, c)
	if err != nil {
		return nil, err
	}
	r.hostRepository = hr
	return &ri, nil
}

// Host return the Host repository instance.
func (r *memoryRepository) Host() *HostRepository {
	return r.hostRepository
}
