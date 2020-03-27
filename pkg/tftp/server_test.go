package tftp

import (
	"github.com/pin/tftp"
	"github.com/pxecore/pxecore/pkg/ipxe"
	"io"
	"strings"
	"testing"
)

func TestServer_tftpReadHandler(t *testing.T) {
	ipxeBiosFile := ipxe.GetIPXEBiosFile()
	ipxeUEFIFile := ipxe.GetIPXEUEFIFile()
	defer ipxe.SetIPXEBiosFile(ipxeBiosFile)
	defer ipxe.SetIPXEUEFIFile(ipxeUEFIFile)
	ipxe.SetIPXEBiosFile([]byte("IPXEBiosFile"))
	ipxe.SetIPXEUEFIFile([]byte("IPXEUEFIFile"))

	type fields struct {
		config            *ServerConfig
		server            *tftp.Server
		defaultIPXEScript []byte
		ipxeScript        IPXEScript
	}
	type args struct {
		filename string
		rf       io.ReaderFrom
	}
	tests := []struct {
		fields  fields
		args    args
		wantErr bool
	}{
		{fields{defaultIPXEScript: []byte("default"),
			ipxeScript: &IPXEScriptValidator{t, "", "IPXEScriptValidator"}},
			args{filename: "undionly.kpxe",
				rf: &ReaderFromValidator{t, "IPXEBiosFile"},},
			false},
		{fields{defaultIPXEScript: []byte("default"),
			ipxeScript: &IPXEScriptValidator{t, "", "IPXEScriptValidator"}},
			args{filename: "ipxe.efi",
				rf: &ReaderFromValidator{t, "IPXEUEFIFile"},},
			false},
		{fields{defaultIPXEScript: []byte("default"),
			ipxeScript: &IPXEScriptValidator{t, "", ""}},
			args{filename: "", rf: &ReaderFromValidator{t, "default"},},
			false},
		{fields{defaultIPXEScript: []byte("default"),
			ipxeScript: &IPXEScriptValidator{t, "36-c6-49-6f-72-d7", "IPXEScriptValidator"}},
			args{filename: "mac-36-c6-49-6f-72-d7.ipxe",
				rf: &ReaderFromValidator{t, "IPXEScriptValidator"},},
			false},
		{fields{defaultIPXEScript: []byte("default"),
			ipxeScript: &IPXEScriptValidator{t, "36-c6-49-6f-72-d7", "IPXEScriptValidator"}},
			args{filename: "mac-36-c6-49-6f--d7.ipxe",
				rf: &ReaderFromValidator{t, "default"},},
			false},
	}
	for _, tt := range tests {
		s := &Server{
			config:            tt.fields.config,
			server:            tt.fields.server,
			defaultIPXEScript: tt.fields.defaultIPXEScript,
			ipxeScript:        &tt.fields.ipxeScript,
		}
		if err := s.tftpReadHandler(tt.args.filename, tt.args.rf); (err != nil) != tt.wantErr {
			t.Errorf("tftpReadHandler() error = %v, wantErr %v", err, tt.wantErr)
		}
	}
}

type ReaderFromValidator struct {
	T          *testing.T
	Validation string
}

func (r *ReaderFromValidator) ReadFrom(rf io.Reader) (n int64, err error) {
	bucket := make([]byte, 100)
	if _, err = rf.Read(bucket); err != nil {
		r.T.Error(err)
		return 0, err
	}
	s := strings.Replace(string(bucket), "\x00", "", -1)
	if r.Validation != "" && s != r.Validation {
		r.T.Fatalf("ReaderFromValidator validation error Expected: %s - Received: %s", r.Validation, s)
	}
	return 0, nil
}

type IPXEScriptValidator struct {
	T      *testing.T
	Mac    string
	Script string
}

func (i *IPXEScriptValidator) Lookup(mac string) string {
	if i.Mac != "" && i.Mac != mac {
		i.T.Fatalf("IPXEScriptValidator validation error Expected: %s - Received: %s", i.Mac, mac)
	}
	return i.Script
}
