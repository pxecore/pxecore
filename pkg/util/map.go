package util

import (
	"fmt"
	"github.com/pxecore/pxecore/pkg/errors"
)

// IntFromMap extract an integer from a map if the type is incorrect will return an error.
// If the key is missing will return the default value.
func IntFromMap(m map[string]interface{}, k string, d int) (int, error) {
	if e, ok := m[k]; ok {
		val, ok := e.(int)
		if !ok {
			return val, &errors.Error{
				Code: errors.EInvalidType,
				Msg:  fmt.Sprint("[util.IntFromMap] invalid type for ", k),
			}
		}
		return val, nil
	}
	return d, nil
}

// StringFromMap extract an string from a map if the type is incorrect will return an error.
// If the key is missing will return the default value.
func StringFromMap(m map[string]interface{}, k string, d string) (string, error) {
	if e, ok := m[k]; ok {
		val, ok := e.(string)
		if !ok {
			return val, &errors.Error{
				Code: errors.EInvalidType,
				Msg:  fmt.Sprint("[util.StringFromMap] invalid type for ", k),
			}
		}
		return val, nil
	}
	return d, nil
}
