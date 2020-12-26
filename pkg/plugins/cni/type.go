package plugins

import "net"

type CNIIPResult struct {
	Addresses []CNIAddress
}

type CNIAddress struct {
	IP      net.IP
	Gateway net.IP
}
