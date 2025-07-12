package libvirtclient

import "C"
import (
	"errors"
	"unsafe"
)

type Client struct {
	conn *C.virConnectPtr
}

func Connect(uri string) (*Client, error) {
	cUri := C.CString(uri)
	defer C.free(unsafe.Pointer(cUri))

	conn := C.virConnectOpen(cUri)
	if conn == nil {
		return nil, errors.New("failed to open connection to libvirt")
	}
	return &Client{conn: conn}, nil
}

func (c *Client) Close() {
	if c.conn != nil {
		C.virConnectClose(c.conn)
		c.conn = nil
	}
}
