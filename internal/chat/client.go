package chat

import (
	"net"
)

type Client struct {
	Name   string
	Conn   net.Conn
	Outbox chan string
}
