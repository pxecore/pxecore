package ipxe

import (
	"log"
	"testing"
)

func TestGetIPXEFile(t *testing.T) {
	SetIPXEBiosFile([]byte{72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100})
	s := string(GetIPXEBiosFile())
	if s != "Hello World" {
		log.Fatal("Wrong Byte to String")
	}

	SetIPXEBiosFile([]byte{})
	SetIPXEUEFIFile([]byte{72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100})
	s = string(GetIPXEUEFIFile())
	if s != "Hello World" {
		log.Fatal("Wrong Byte to String")
	}
}
