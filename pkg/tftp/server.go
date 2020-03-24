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
	log "github.com/sirupsen/logrus"
	"gitlab.com/pliego/pxe-injector/pkg/ipxe"
	"io"
	"time"
)

const (
	// IPXEBiosFilename holds the name of a legacy BIOS IPXE filename.
	IPXEBiosFilename string = "undionly.kpxe"
	// IPXEEFIFilename holds the name of a UEFI IPXE filename.
	IPXEEFIFilename string = "ipxe.efi"
)

// ServerConfig holds the information that will be used to configure the TFTP Server.
type ServerConfig struct {
	// Address tftp server will listen to. Example: ":69".
	Address string
	// Timeout duration when a connection will be closed. Example: "5 * time.Second".
	Timeout time.Duration
}

// Server is the representation of the TFTP server for this domain.
type Server struct {
	config *ServerConfig
	server *tftp.Server
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
	s.server = tftp.NewServer(s.tftpReadHandler, nil)
	s.server.SetTimeout(config.Timeout)
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
	log.Info(filename)
	switch filename {
	case IPXEBiosFilename:
		_, _ = rf.ReadFrom(bytes.NewReader(ipxe.GetIPXEBiosFile()))
		break
	case IPXEEFIFilename:
		_, _ = rf.ReadFrom(bytes.NewReader(ipxe.GetIPXEUEFIFile()))
		break
	default:
		i := "#!ipxe\n\ndhcp\nchain --autofree https://boot.netboot.xyz"
		_, _ = rf.ReadFrom(bytes.NewReader([]byte(i)))
	}
	return nil
}
