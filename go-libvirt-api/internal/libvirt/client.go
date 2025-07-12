package libvirt

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

// ConnectLibvirt establishes a new connection to the libvirt host
func ConnectLibvirt(uri string) (*libvirt.Connect, error) {
	conn, err := libvirt.NewConnect(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to libvirt at %s: %v", uri, err)
	}
	return conn, nil
}
