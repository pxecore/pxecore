package locator

import (
	"github.com/pxecore/pxecore/pkg/entity"
	"github.com/pxecore/pxecore/pkg/repository"
	"reflect"
	"testing"
)

func TestRepositoryIPXEScript_Lookup(t *testing.T) {
	r, _ := repository.NewRepository(map[string]interface{}{"driver": "memory"})
	s, _ := r.Open(true)
	_ = s.Template().Create(entity.Template{ID: "template", Template: "A"})
	_ = s.Host().Create(entity.Host{ID: "host", HardwareAddr: []string{"88-99-aa-bb-cc-dd"},
		Vars: map[string]string{}, TemplateID: "template"})
	_ = s.Close()
	tests := []struct {
		name    string
		path    string
		want    []byte
		wantErr bool
	}{
		{"OK_1", "mac-88-99-aa-bb-cc-dd.ipxe", []byte{65}, false},
		{"OK_2", "mac-88-99-AA-BB-CC-DD.ipxe", []byte{65}, false},
		{"OK_3", "pxelinux.cfg/01-88-99-AA-BB-CC-DD", []byte{65}, false},
		{"OK_4", "pxelinux.cfg/01-88-99-AA-BB-CC-DD", []byte{65}, false},
		{"KO_1", "pxelinux.cfg/01-88-99-AA-BB-CC-EE", nil, true},
		{"KO_2", "none", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewRepositoryIPXEScript(r)
			g, err := s.Lookup(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lookup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				got := make([]byte, 1)
				g.Read(got)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Lookup() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestRepositoryIPXEScript_MatchIPXEPath(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		want  string
		want1 bool
	}{
		{"OK", "mac-12-d7-d5-cb-5a-43.ipxe", "12-d7-d5-cb-5a-43", true},
		{"KO_1", "mac-12-d7-d5-cb-5a-.ipxe", "", false},
		{"KO_2", "mac-12-d7-d5-cb-5a-43", "", false},
		{"KO_3", "mac-12-d7-d5-.ipxe", "", false},
		{"KO_4", "mac-12-d7-d.ipxe", "", false},
		{"KO_5", "mac-12-.ipxe", "", false},
		{"KO_6", "mac-1.ipxe", "", false},
		{"KO_7", "mac-.ipxe", "", false},
		{"KO_8", "mac.ipxe", "", false},
		{"KO_9", "mac", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewRepositoryIPXEScript(nil)
			got, got1 := s.MatchIPXEPath(tt.path)
			if got != tt.want {
				t.Errorf("MatchIPXEPath() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("MatchIPXEPath() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestRepositoryIPXEScript_MatchPXELINUXPathPattern(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		want  string
		want1 bool
	}{
		{"OK", "pxelinux.cfg/01-88-99-aa-bb-cc-dd", "88-99-aa-bb-cc-dd", true},
		{"OK", "pxelinux.cfg/01-12-d7-d5-cb-5a-43", "12-d7-d5-cb-5a-43", true},
		{"KO_1", "pxelinux.cfg/01-12-d7-d5-cb-5a", "", false},
		{"KO_2", "pxelinux.cfg/01-12-d7-d5-cb", "", false},
		{"KO_3", "pxelinux.cfg/01-12-d7-d5-", "", false},
		{"KO_4", "pxelinux.cfg/01-12-d7-d", "", false},
		{"KO_5", "pxelinux.cfg/01-12-", "", false},
		{"KO_6", "pxelinux.cfg/01-1", "", false},
		{"KO_7", "pxelinux.cfg/01-", "", false},
		{"KO_8", "pxelinux.cfg/01", "", false},
		{"KO_9", "pxelinux.cfg/", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewRepositoryIPXEScript(nil)
			got, got1 := s.MatchPXELINUXPathPattern(tt.path)
			if got != tt.want {
				t.Errorf("MatchIPXEPath() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("MatchIPXEPath() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
