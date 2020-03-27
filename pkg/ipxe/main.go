//go:generate go run generator.go

package ipxe

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

// GetIPXEUEFIFile retrieves the IPXE UEFI file from memory
func GetIPXEUEFIFile() []byte {
	return ipxeUEFIFile
}

// SetIPXEUEFIFile stores the IPXE UEFI file in memory
func SetIPXEUEFIFile(file []byte) {
	ipxeUEFIFile = file
}
