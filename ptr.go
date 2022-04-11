package main

import "net"

type PTRRecord struct {
	IP    net.IP
	Names []string
}
