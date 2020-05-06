package script

import (
	"bytes"
	"github.com/pxecore/pxecore/pkg/repository"
	"github.com/pxecore/pxecore/pkg/template"
)

// RepositoryIPXEScript searches the for the mac address in the configured repository.
type RepositoryIPXEScript struct {
	Repository repository.Repository
}

// Lookup returns the script for the provided mac address.
// See github.com/pxecore/pxecore/pkg/tftp/IPXEScript
func (s RepositoryIPXEScript) Lookup(mac string) string {
	buf := new(bytes.Buffer)
	_ = template.CompileWithHardwareAddr(buf, s.Repository, mac, "")
	return buf.String()
}
