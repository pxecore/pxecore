package errors

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	// EUnknown code for Unknown error.
	EUnknown string = "EUnknown"
	// EInvalidType code for invalid type in type conversion.
	EInvalidType string = "EInvalidType"
	// ERepositoryKeyExist code when a key exists where it shouldn't.
	ERepositoryKeyExist string = "ERepositoryKeyExist"
	// ERepositoryKeyNotFound code when a key was not found.
	ERepositoryKeyNotFound string = "ERepositoryKeyDontExist"
	// ERepositoryEmptyKey code when a key should have been provided.
	ERepositoryEmptyKey string = "ERepositoryEmptyKey"
)

// Error data structure
type Error struct {
	Code string `json:"code"`
	Msg  string `json:"message,omitempty"`
	Err  error  `json:"error,omitempty"`
}

// Error implements golang's error interface.
func (e *Error) Error() string {
	var s strings.Builder
	s.WriteString(e.Code)
	s.WriteString(":")
	if e.Msg != "" {
		s.WriteString(fmt.Sprintf(" %s", e.Msg))
	}
	if e.Err != nil {
		s.WriteString(fmt.Sprintf(" <%s>", e.Err.Error()))
	}
	return s.String()
}

// Is checks if an error is of a certain type.
func Is(err error, code string) bool {
	val, ok := err.(*Error)
	if !ok {
		return false
	}
	return val.Code == code
}

// MarshalJSON converts an error into JSON.
// If the error is a default golang's stringError it will wrap it into an EUnknown error.
// If the ``json.Marshal`` the whole error is wrapped in based64 for further processing.
func MarshalJSON(err error) []byte {
	val, ok := err.(*Error)
	if !ok {
		return MarshalJSON(&Error{
			Code: EUnknown,
			Msg:  fmt.Sprint(err.Error()),
		})
	}
	j, err := json.Marshal(val)
	if err != nil {
		return MarshalJSON(&Error{
			Code: val.Code,
			Msg: fmt.Sprint("can't convert error to JSON Base64:",
				base64.StdEncoding.EncodeToString([]byte(val.Error()))),
		})
	}
	return j
}
