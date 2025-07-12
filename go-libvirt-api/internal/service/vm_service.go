package service

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

func StartVM(conn *libvirt.Connect, vmName string) error {
	dom, err := conn.LookupDomainByName(vmName)
	if err != nil {
		return fmt.Errorf("failed to find VM %s: %v", vmName, err)
	}
	return dom.Create()
}

func StopVM(conn *libvirt.Connect, vmName string) error {
	dom, err := conn.LookupDomainByName(vmName)
	if err != nil {
		return fmt.Errorf("failed to find VM %s: %v", vmName, err)
	}
	return dom.Shutdown()
}

func RebootVM(conn *libvirt.Connect, vmName string) error {
	dom, err := conn.LookupDomainByName(vmName)
	if err != nil {
		return fmt.Errorf("failed to find VM %s: %v", vmName, err)
	}
	return dom.Reboot(libvirt.DOMAIN_REBOOT_DEFAULT)
}

func PauseVM(conn *libvirt.Connect, vmName string) error {
	dom, err := conn.LookupDomainByName(vmName)
	if err != nil {
		return fmt.Errorf("failed to find VM %s: %v", vmName, err)
	}
	return dom.Suspend()
}

func ResumeVM(conn *libvirt.Connect, vmName string) error {
	dom, err := conn.LookupDomainByName(vmName)
	if err != nil {
		return fmt.Errorf("failed to find VM %s: %v", vmName, err)
	}
	return dom.Resume()
}

func DeleteVM(conn *libvirt.Connect, vmName string) error {
	dom, err := conn.LookupDomainByName(vmName)
	if err != nil {
		return fmt.Errorf("failed to find VM %s: %v", vmName, err)
	}
	if err := dom.Destroy(); err != nil {
		return fmt.Errorf("failed to destroy VM %s: %v", vmName, err)
	}
	return dom.Undefine()
}
