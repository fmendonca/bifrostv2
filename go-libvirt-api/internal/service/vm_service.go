package service

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

func StartVM(conn *libvirt.Connect, vmName string) error {
	dom, err := conn.LookupDomainByName(vmName)
	if err != nil {
		return fmt.Errorf("failed to lookup domain %s: %v", vmName, err)
	}
	return dom.Create()
}

func StopVM(conn *libvirt.Connect, vmName string) error {
	dom, err := conn.LookupDomainByName(vmName)
	if err != nil {
		return fmt.Errorf("failed to lookup domain %s: %v", vmName, err)
	}
	return dom.Shutdown()
}

func RebootVM(conn *libvirt.Connect, vmName string) error {
	dom, err := conn.LookupDomainByName(vmName)
	if err != nil {
		return fmt.Errorf("failed to lookup domain %s: %v", vmName, err)
	}
	// no flags are defined, pass 0
	return dom.Reboot(0)
}

func PauseVM(conn *libvirt.Connect, vmName string) error {
	dom, err := conn.LookupDomainByName(vmName)
	if err != nil {
		return fmt.Errorf("failed to lookup domain %s: %v", vmName, err)
	}
	return dom.Suspend()
}

func ResumeVM(conn *libvirt.Connect, vmName string) error {
	dom, err := conn.LookupDomainByName(vmName)
	if err != nil {
		return fmt.Errorf("failed to lookup domain %s: %v", vmName, err)
	}
	return dom.Resume()
}

func DeleteVM(conn *libvirt.Connect, vmName string) error {
	dom, err := conn.LookupDomainByName(vmName)
	if err != nil {
		return fmt.Errorf("failed to lookup domain %s: %v", vmName, err)
	}
	if err := dom.Destroy(); err != nil {
		return fmt.Errorf("failed to destroy domain %s: %v", vmName, err)
	}
	return dom.Undefine()
}
