//go:generate go run generator.go

package ipxe

import "io/ioutil"

var ipxeBiosFile []byte
var ipxeUEFIFile []byte

// GetIPXEBiosFile retrieves the IPXE Bios file from memory
func GetIPXEBiosFile() []byte {
	return ipxeBiosFile
}

// SetIPXEBiosFile stores the IPXE Bios file in memory
func SetIPXEBiosFile(file []byte) {
	ipxeBiosFile = file
}

// LoadIPXEBiosFile reads the provided path and load it's binary content.
func LoadIPXEBiosFile(path string) error {
	read, err := ioutil.ReadFile(path)
	if err != nil{
		return err
	}
	SetIPXEBiosFile(read)
	return nil
}

// GetIPXEUEFIFile retrieves the IPXE UEFI file from memory
func GetIPXEUEFIFile() []byte {
	return ipxeUEFIFile
}

// SetIPXEUEFIFile stores the IPXE UEFI file in memory
func SetIPXEUEFIFile(file []byte) {
	ipxeUEFIFile = file
}

// LoadIPXEUEFIFile reads the provided path and load it's binary content.
func LoadIPXEUEFIFile(path string) error {
	read, err := ioutil.ReadFile(path)
	if err != nil{
		return err
	}
	SetIPXEUEFIFile(read)
	return nil
}
