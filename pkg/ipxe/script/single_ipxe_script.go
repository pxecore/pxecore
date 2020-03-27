package script

import "io/ioutil"

// SingleIPXEScript returns the same script for all mac addresses.
type SingleIPXEScript struct {
	ipxeScript string
}

// NewSingleIPXEScript instantiates a new SingleIPXEScript with it's read-only attributes
func NewSingleIPXEScript(script string) *SingleIPXEScript {
	return &SingleIPXEScript{ipxeScript: script}
}

// NewSingleIPXEScriptFromFile instantiates a new SingleIPXEScript from a file.
func NewSingleIPXEScriptFromFile(path string) (*SingleIPXEScript, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return &SingleIPXEScript{ipxeScript: string(file)}, nil

}

// Lookup returns the script for the provided mac address.
// See gitlab.com/pliego/pxecore/pkg/tftp/IPXEScript
func (s *SingleIPXEScript) Lookup(mac string) string {
	return s.ipxeScript
}
