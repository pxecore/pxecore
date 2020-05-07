# PXECORE - Preboot eXecution Environment Core

**Important: This project is under heavy development so breaking changes are expected as well as
incomplete documentation**

PXECORE facilitates the installation of the physical and virtual machine by
providing mechanisms of provisioning over the IPXE Firmware. It aims to assist
both home server enthusiasts as bare metal cloud administrators alike.

![Test](https://github.com/pxecore/pxecore/workflows/Test/badge.svg?branch=master)

## Installation

### GO

* Install [golang](https://golang.org/).
* Get Binary:

```shell
go get github.com/pxecore/pxecore@latest
```

* Play:

```shell
pxecore
```

## Documentation

Documentation and samples are located at https://pxecore.org/.

## Dependencies

PXECORE wouldn't be possible without the help of the following projects:

* [iPXE](https://ipxe.org/) 
[[License](https://github.com/ipxe/ipxe/blob/master/COPYING)] - The leading open source network boot firmware.
* [pin/tftp](https://github.com/pin/tftp) 
[[License](https://github.com/pin/tftp/blob/master/LICENSE)] - TFTP server and client library for Golang.  
* [gorilla/mux](https://github.com/gorilla/mux) 
[[License](https://github.com/gorilla/mux/blob/master/LICENSE)] - Request router and dispatcher
* [spf13/pflag](https://github.com/spf13/pflag) 
[[License](https://github.com/spf13/pflag/blob/master/LICENSE)] - Replacement of Go's native flag package
* [spf13/viper](https://github.com/spf13/viper) 
[[License](https://github.com/spf13/viper/blob/master/LICENSE)] - Complete configuration solution for Go.  

## Contact 

* Issues: https://github.com/pxecore/pxecore/issues
* Telegram Group: https://t.me/joinchat/HBYQWho5u5dhgbD6hcIFRA
* Email: contact(at)pxecore.org

### Author

* Martin Pliego. [Website](https://github.com/mpliego) -
[Github](https://github.com/mpliego) - [Linkdin](https://github.com/mpliego)

