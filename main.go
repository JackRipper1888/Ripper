package main

import (
	"Ripper/udpkit"
)

func main() {
	udpkit.Peerclient("192.168.0.17", 8091)
}
