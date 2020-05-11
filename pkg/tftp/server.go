// Package tftp exposes a read-only TFTP server.
// See: github.com/pin/tftp
package tftp

import (
	"bytes"
	"github.com/pin/tftp"
	"github.com/pxecore/pxecore/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"time"
)

// ServerConfig holds the information that will be used to configure the TFTP Server.
type ServerConfig struct {
	// Address tftp server will listen to. Example: ":69".
	Address string
	// Timeout duration when a connection will be closed. Example: "5 * time.Second".
	Timeout time.Duration
	// FileLocators is used to retrieve the file of a particular file.
	FileLocators []FileLocator
}

// FileLocator implements the IPXE static lookup procedure.
type FileLocator interface {
	// Lookup finds and returns the IPXE static suitable for the mac address provided.
	Lookup(path string) ([]byte, error)
}

// Server is the representation of the TFTP server for this domain.
type Server struct {
	config       *ServerConfig
	server       *tftp.Server
	fileLocators []FileLocator
}

// StartInBackground starts the TFTP server in a different goroutine.
func (s *Server) StartInBackground(config ServerConfig) error {
	go s.Start(config)
	return nil
}

// Start initiates the TFTP server blocking the current goroutine.
func (s *Server) Start(config ServerConfig) error {
	if s.config != nil {
		return &errors.Error{Code: errors.EUnknown, Msg: "[tftp] server not started"}
	}
	s.config = &config
	s.fileLocators = config.FileLocators
	s.server = tftp.NewServer(s.tftpReadHandler, nil)
	s.server.SetTimeout(config.Timeout)
	log.WithField("address", config.Address).Info("TFTP Listening...")
	err := s.server.ListenAndServe(config.Address)
	return err
}

// Shutdown stops the current server.
func (s *Server) Shutdown() error {
	if s.config != nil {
		return &errors.Error{Code: errors.EUnknown, Msg: "[tftp] server not started"}
	}
	s.server.Shutdown()
	s.config = nil
	return nil
}

// tftpReadHandler handles a read event in the TFTP server.
func (s Server) tftpReadHandler(path string, rf io.ReaderFrom) error {
	for _, v := range s.fileLocators {
		r, err := v.Lookup(path)
		if err != nil {
			if !errors.Is(err, errors.ENotFound) {
				log.WithError(err).Error("[tftp] Error locating file.")
			}
			continue
		}
		if _, err := rf.ReadFrom(bytes.NewReader(r)); err != nil {
			log.WithError(err).Error("Error sending TFTP response")
			return nil
		}
		log.WithFields(log.Fields{"filename": path}).Info("IPXE Script Sent.")
		return nil
	}
	return &errors.Error{Code: errors.ENotFound, Msg: "[tftp] file not found."}
}

// sendResponse writes the TFTP response and reports errors.
func (s *Server) sendResponse(rf io.ReaderFrom, response []byte) {
	if _, err := rf.ReadFrom(bytes.NewReader(response)); err != nil {
		log.WithError(err).Error("Error sending TFTP response")
	}
}
