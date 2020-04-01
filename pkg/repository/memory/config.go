package memory

import "github.com/pxecore/pxecore/pkg/errors"

// Config stores memory driver config for all repositories.
type Config struct {
	allowReset bool
}

// NewConfig creates a new Config extracting and checking type of the required fields.
func NewConfig(r *Repository, config map[string]interface{}) (*Config, error) {
	c := new(Config)

	c.allowReset = false
	if e, ok := config["allow-reset"]; ok {
		val, ok := e.(bool)
		if !ok {
			return nil, &errors.Error{Code: errors.EInvalidType, Msg: "config invalid type for key allow-reset"}
		}
		c.allowReset = val
	}

	return c, nil
}
