package ipxe

import (
	"log"
	"testing"
)

func TestGetIPXEFile(t *testing.T) {
	ipxeBiosFile = []byte{72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100}
	s := string(GetIPXEBiosFile())
	t.Log(s)
	if s != "Hello World" {
		log.Fatal("Wrong Byte to String")
	}

	ipxeBiosFile = []byte{}
	ipxeUEFIFile = []byte{72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100}
	s = string(GetIPXEUEFIFile())
	t.Log(s)
	if s != "Hello World" {
		log.Fatal("Wrong Byte to String")
	}
}
