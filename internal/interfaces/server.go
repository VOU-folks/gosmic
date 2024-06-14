package interfaces

import (
	"net"
)

type Server interface {
	Listen() error
	Close() error
	Listener() *net.Listener
}
