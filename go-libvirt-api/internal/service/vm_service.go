package service

/*
#cgo pkg-config: libvirt
#include <libvirt/libvirt.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"unsafe"

	"go-libvirt-api/internal/libvirtclient"
)

func getDomainByName(client *libvirtclient.Client, name string) (*C.virDomainPtr, error) {
	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	dom := C.virDomainLookupByName(client.Conn(), cName)
	if dom == nil {
		return nil, errors.New("failed to find domain")
	}
	return &dom, nil
}

func StartVM(client *libvirtclient.Client, name string) error {
	dom, err := getDomainByName(client, name)
	if err != nil {
		return err
	}
	defer C.virDomainFree(*dom)

	if C.virDomainCreate(*dom) < 0 {
		return errors.New("failed to start domain")
	}
	return nil
}

func StopVM(client *libvirtclient.Client, name string) error {
	dom, err := getDomainByName(client, name)
	if err != nil {
		return err
	}
	defer C.virDomainFree(*dom)

	if C.virDomainShutdown(*dom) < 0 {
		return errors.New("failed to shutdown domain")
	}
	return nil
}
