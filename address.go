package main

import "net"

type Address struct {
	IP    net.IP
	Names []string
}
