package models

import (
	"net"
)

type RequrestInfo struct {
	Addr       *net.UDPAddr
	Data       [1024]byte
	CountTotal int
}

type ResponseInfo struct {
	Addr *net.UDPAddr
	Data []byte
}
