package service

import (
	"fmt"
	"go-libvirt-api/internal/libvirtclient"
)

func StartVM(client *libvirtclient.Client, name string) error {
	dom, err := client.Conn().LookupDomainByName(name)
	if err != nil {
		return fmt.Errorf("failed to lookup domain %s: %v", name, err)
	}
	defer dom.Free()
	return dom.Create()
}

func StopVM(client *libvirtclient.Client, name string) error {
	dom, err := client.Conn().LookupDomainByName(name)
	if err != nil {
		return fmt.Errorf("failed to lookup domain %s: %v", name, err)
	}
	defer dom.Free()
	return dom.Shutdown()
}

func RebootVM(client *libvirtclient.Client, name string) error {
	dom, err := client.Conn().LookupDomainByName(name)
	if err != nil {
		return fmt.Errorf("failed to lookup domain %s: %v", name, err)
	}
	defer dom.Free()
	return dom.Reboot(0)
}

func PauseVM(client *libvirtclient.Client, name string) error {
	dom, err := client.Conn().LookupDomainByName(name)
	if err != nil {
		return fmt.Errorf("failed to lookup domain %s: %v", name, err)
	}
	defer dom.Free()
	return dom.Suspend()
}

func ResumeVM(client *libvirtclient.Client, name string) error {
	dom, err := client.Conn().LookupDomainByName(name)
	if err != nil {
		return fmt.Errorf("failed to lookup domain %s: %v", name, err)
	}
	defer dom.Free()
	return dom.Resume()
}

func DeleteVM(client *libvirtclient.Client, name string) error {
	dom, err := client.Conn().LookupDomainByName(name)
	if err != nil {
		return fmt.Errorf("failed to lookup domain %s: %v", name, err)
	}
	defer dom.Free()
	if err := dom.Destroy(); err != nil {
		return fmt.Errorf("failed to destroy domain %s: %v", name, err)
	}
	return dom.Undefine()
}
