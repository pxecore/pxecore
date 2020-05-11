package locator

import (
	"github.com/pxecore/pxecore/pkg/errors"
	"github.com/pxecore/pxecore/pkg/ipxe"
)

const (
	// IPXEBiosFilename holds the name of a legacy BIOS IPXE filename.
	IPXEBiosFilename string = "undionly.kpxe"
	// IPXEEFIFilename holds the name of a UEFI IPXE filename.
	IPXEEFIFilename string = "ipxe.efi"
)

// IPXEFirmware returns the same locator for all mac addresses.
type IPXEFirmware struct {
}

// NewIPXEFirmware instantiates a new SingleIPXEScript with it's read-only attributes
func NewIPXEFirmware() *IPXEFirmware {
	return new(IPXEFirmware)
}

// Lookup returns the locator for the provided path.
// See gitlab.com/pliego/pxecore/pkg/tftp/IPXEScript
func (s IPXEFirmware) Lookup(path string) ([]byte, error) {
	switch path {
	case IPXEBiosFilename:
		return ipxe.GetIPXEBiosFile(), nil
	case IPXEEFIFilename:
		return ipxe.GetIPXEUEFIFile(), nil
	}
	return nil, &errors.Error{Code: errors.ENotFound, Msg: "[tftp.locator] firmware not found."}
}
