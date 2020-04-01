package errors

import (
	"errors"
	"reflect"
	"testing"
)

func TestError_Error(t *testing.T) {
	type fields struct {
		Code string
		Msg  string
		Err  error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Case#1", fields{"parent", "parent msg",
			&Error{"child", "child msg", nil}},
			"parent: parent msg <child: child msg>"},
		{"Case#2", fields{"parent", "parent msg",
			&Error{"child", "child msg",
				&Error{"child2", "child2 msg", nil}}},
			"parent: parent msg <child: child msg <child2: child2 msg>>"},
		{"Case#3", fields{"parent", "parent msg",
			&Error{"child", "",
				&Error{"child2", "child2 msg", nil}}},
			"parent: parent msg <child: <child2: child2 msg>>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				Code: tt.fields.Code,
				Msg:  tt.fields.Msg,
				Err:  tt.fields.Err,
			}
			f := func() error {
				return e
			}
			if err := f(); err == nil {
				t.Error("Error() error missing")
			} else {
				if !reflect.DeepEqual(err.Error(), tt.want) {
					t.Errorf("Error() got = %v, want %v", err, tt.want)
				}
			}
		})
	}
}

func TestError_JSON(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"Case#1", &Error{"parent", "parent msg",
			&Error{"child", "child msg", nil}},
			"{\"code\":\"parent\",\"message\":\"parent msg\"," +
				"\"error\":{\"code\":\"child\",\"message\":\"child msg\"}}"},
		{"Case#2", errors.New("string error"),
			"{\"code\":\"EUnknown\",\"message\":\"string error\"}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MarshalJSON(tt.err)
			if !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("JSON() got = %v, want %v", string(got), tt.want)
			}
		})
	}
}
