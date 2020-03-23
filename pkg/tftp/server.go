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
	IPXEBiosFilename string = "undionly.kpxe"
	IPXEEFIFilename string = "ipxe.efi"
)


type ServerConfig struct {
	// Address tftp server will listen to. Example: ":69"
	Address string
	// Timeout duration when a connection will be closed. Example: "5 * time.Second"
	Timeout time.Duration
}

type Server struct {
	config *ServerConfig
	server *tftp.Server
}

func (s *Server) StartInBackground(config ServerConfig) {
	go s.Start(config)
}

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

func (s *Server) Shutdown() error {
	if s.config != nil {
		return errors.New("server not started")
	}
	s.server.Shutdown()
	return nil
}

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
