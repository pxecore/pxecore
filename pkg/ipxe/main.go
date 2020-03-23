//go:generate go run generator.go

package ipxe

var ipxeBiosFile []byte
var ipxeUEFIFile []byte

func GetIPXEBiosFile() []byte {
	return ipxeBiosFile
}

func GetIPXEUEFIFile() []byte {
	return ipxeUEFIFile
}
