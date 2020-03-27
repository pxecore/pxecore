// Package tftp exposes a read-only TFTP server.
// See: github.com/pin/tftp
//
// If the Default PXE filenames are requested both IPXE are delivered from memory.
// Once the IPXE loads a new file is requested following the format "mac-${mac-address-hyp}.ipxe"
// Example: mac-1d-af-02-34-ef-77.ipxe
package tftp

import (
	"bytes"
	"errors"
	"github.com/pin/tftp"
	"github.com/pxecore/pxecore/pkg/ipxe"
	log "github.com/sirupsen/logrus"
	"io"
	"regexp"
	"strings"
	"time"
)

const (
	// IPXEBiosFilename holds the name of a legacy BIOS IPXE filename.
	IPXEBiosFilename string = "undionly.kpxe"
	// IPXEEFIFilename holds the name of a UEFI IPXE filename.
	IPXEEFIFilename string = "ipxe.efi"
	// DefaultIPXEScript holds the default static used with no other option is available.
	DefaultIPXEScript string = "#!ipxe\n\necho No script defined exiting.\nsleep 5"
)

var macAddressRegex, _ = regexp.Compile("^mac-(([0-9a-f]{2}[-]){5}([0-9a-f]{2}))\\.ipxe$")

// ServerConfig holds the information that will be used to configure the TFTP Server.
type ServerConfig struct {
	// Address tftp server will listen to. Example: ":69".
	Address string
	// Timeout duration when a connection will be closed. Example: "5 * time.Second".
	Timeout time.Duration
	// IPXEScript is used to retrieve the IPXE particular to a host mac address.
	// If SingleModeFile is present it will be ignored and if not present the DefaultIPXEScript will be used.
	IPXEScript *IPXEScript
}

// IPXEScript implements the IPXE static lookup procedure.
type IPXEScript interface {
	// Lookup finds and returns the IPXE static suitable for the mac address provided.
	Lookup(mac string) string
}

// Server is the representation of the TFTP server for this domain.
type Server struct {
	config            *ServerConfig
	server            *tftp.Server
	defaultIPXEScript []byte
	ipxeScript        *IPXEScript
}

// StartInBackground starts the TFTP server in a different goroutine.
func (s *Server) StartInBackground(config ServerConfig) {
	go s.Start(config)
}

// Start initiates the TFTP server blocking the current goroutine.
func (s *Server) Start(config ServerConfig) error {
	if s.config != nil {
		return errors.New("server already started")
	}
	s.config = &config
	s.defaultIPXEScript = []byte(DefaultIPXEScript)
	s.ipxeScript = config.IPXEScript
	s.server = tftp.NewServer(s.tftpReadHandler, nil)
	s.server.SetTimeout(config.Timeout)
	log.WithField("address", config.Address).Info("TFTP Listening...")
	err := s.server.ListenAndServe(config.Address)
	return err
}

// Shutdown stops the current server.
func (s *Server) Shutdown() error {
	if s.config != nil {
		return errors.New("server not started")
	}
	s.server.Shutdown()
	s.config = nil
	return nil
}

// tftpReadHandler handles a read event in the TFTP server.
// If the Default PXE filenames are requested both IPXE are delivered from memory.
// Once the IPXE loads a new file is requested following the format "mac-${mac-address-hyp}.ipxe"
// Example: mac-1d-af-02-34-ef-77.ipxe
func (s *Server) tftpReadHandler(filename string, rf io.ReaderFrom) error {
	switch filename {
	case IPXEBiosFilename:
		s.sendResponse(rf, ipxe.GetIPXEBiosFile())
		break
	case IPXEEFIFilename:
		s.sendResponse(rf, ipxe.GetIPXEUEFIFile())
		break
	default:
		var response []byte
		if *s.ipxeScript != nil {
			fn := strings.ToLower(filename)
			g := macAddressRegex.FindStringSubmatch(fn)
			if len(g) > 1 {
				response = []byte(strings.TrimSpace((*s.ipxeScript).Lookup(g[1])))
			}
		}
		if len(response) == 0 {
			response = s.defaultIPXEScript
		}
		s.sendResponse(rf, response)
		log.WithFields(log.Fields{"filename": filename, "response": string(response)}).
			Debug("IPXE Script Sent.")
	}
	log.WithFields(log.Fields{"filename": filename}).Info("IPXE Script Sent.")
	return nil
}

// sendResponse writes the TFTP response and reports errors.
func (s *Server) sendResponse(rf io.ReaderFrom, response []byte) {
	if _, err := rf.ReadFrom(bytes.NewReader(response)); err != nil {
		log.WithError(err).Error("Error sending TFTP response")
	}
}
