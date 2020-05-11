package locator

import (
	"bytes"
	"github.com/pxecore/pxecore/pkg/errors"
	"github.com/pxecore/pxecore/pkg/repository"
	"github.com/pxecore/pxecore/pkg/template"
	"regexp"
	"strings"
)

// RepositoryIPXEScript searches the for the mac address in the configured repository.
type RepositoryIPXEScript struct {
	repository          repository.Repository
	ipxePathPattern     *regexp.Regexp
	pxelinuxPathPattern *regexp.Regexp
}

// NewRepositoryIPXEScript construct RepositoryIPXEScript
func NewRepositoryIPXEScript(repository repository.Repository) *RepositoryIPXEScript {
	s := new(RepositoryIPXEScript)
	s.repository = repository
	s.ipxePathPattern, _ = regexp.Compile("^mac-(([0-9a-f]{2}[-]){5}([0-9a-f]{2}))\\.ipxe$")
	s.pxelinuxPathPattern, _ = regexp.Compile("^pxelinux\\.cfg/[0-9a-f]{2}-(([0-9a-f]{2}[-]){5}([0-9a-f]{2}))$")
	return s
}

// Lookup returns the locator for the provided mac address.
// See github.com/pxecore/pxecore/pkg/tftp/IPXEScript
func (s RepositoryIPXEScript) Lookup(path string) ([]byte, error) {
	fn := strings.ToLower(path)
	var ha string
	var ok bool
	if ha, ok = s.MatchIPXEPath(fn); !ok {
		if ha, ok = s.MatchPXELINUXPathPattern(fn); !ok {
			return nil, &errors.Error{Code: errors.ENotFound, Msg: "[tftp.locator] Path is not an IPXE script"}
		}
	}
	buf := new(bytes.Buffer)
	if err := template.CompileWithHardwareAddr(buf, s.repository, ha, ""); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MatchIPXEPath searches the hardware address in the IPXE defined path.
func (s RepositoryIPXEScript) MatchIPXEPath(path string) (string, bool) {
	g := s.ipxePathPattern.FindStringSubmatch(path)
	if len(g) > 1 {
		return g[1], true
	}
	return "", false
}

// MatchPXELINUXPathPattern searches the hardware address in the PXELINUX defined path.
func (s RepositoryIPXEScript) MatchPXELINUXPathPattern(path string) (string, bool) {
	g := s.pxelinuxPathPattern.FindStringSubmatch(path)
	if len(g) > 1 {
		return g[1], true
	}
	return "", false
}
