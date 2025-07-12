package libvirtclient

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

type Client struct {
	conn *libvirt.Connect
}

func Connect(uri string) (*Client, error) {
	conn, err := libvirt.NewConnect(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to libvirt at %s: %v", uri, err)
	}
	return &Client{conn: conn}, nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
		c.conn = nil
	}
}

func (c *Client) Conn() *libvirt.Connect {
	return c.conn
}
