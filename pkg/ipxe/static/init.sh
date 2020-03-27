#!/bin/bash

rm -rf ../ipxe
git clone https://github.com/ipxe/ipxe.git ../ipxe

# Change required flags
# See rikka0w0: https://gist.github.com/rikka0w0/50895b82cbec8a3a1e8c7707479824c1
sed -i 's/#undef\tDOWNLOAD_PROTO_NFS/#define\tDOWNLOAD_PROTO_NFS/' ../ipxe/src/config/general.h
sed -i 's/#undef\tDOWNLOAD_PROTO_HTTPS/#define\tDOWNLOAD_PROTO_HTTPS/' ../ipxe/src/config/general.h