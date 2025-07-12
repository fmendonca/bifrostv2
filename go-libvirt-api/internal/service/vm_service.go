package service

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

// StartVM starts a virtual machine by name
func StartVM(conn *libvirt.Connect, vmName string) error {
	dom, err := conn.LookupDomainByName(vmName)
	if err != nil {
		return fmt.Errorf("failed to find VM %s: %v", vmName, err)
	}
	if err := dom.Create(); err != nil {
		return fmt.Errorf("failed to start VM %s: %v", vmName, err)
	}
	return nil
}

// StopVM stops (shuts down) a virtual machine by name
func StopVM(conn *libvirt.Connect, vmName string) error {
	dom, err := conn.LookupDomainByName(vmName)
	if err != nil {
		return fmt.Errorf("failed to find VM %s: %v", vmName, err)
	}
	if err := dom.Shutdown(); err != nil {
		return fmt.Errorf("failed to stop VM %s: %v", vmName, err)
	}
	return nil
}

// RebootVM reboots a virtual machine by name
func RebootVM(conn *libvirt.Connect, vmName string) error {
	dom, err := conn.LookupDomainByName(vmName)
	if err != nil {
		return fmt.Errorf("failed to find VM %s: %v", vmName, err)
	}
	if err := dom.Reboot(libvirt.DOMAIN_REBOOT_DEFAULT); err != nil {
		return fmt.Errorf("failed to reboot VM %s: %v", vmName, err)
	}
	return nil
}

// PauseVM suspends a virtual machine by name
func PauseVM(conn *libvirt.Connect, vmName string) error {
	dom, err := conn.LookupDomainByName(vmName)
	if err != nil {
		return fmt.Errorf("failed to find VM %s: %v", vmName, err)
	}
	if err := dom.Suspend(); err != nil {
		return fmt.Errorf("failed to pause VM %s: %v", vmName, err)
	}
	return nil
}

// ResumeVM resumes a paused virtual machine by name
func ResumeVM(conn *libvirt.Connect, vmName string) error {
	dom, err := conn.LookupDomainByName(vmName)
	if err != nil {
		return fmt.Errorf("failed to find VM %s: %v", vmName, err)
	}
	if err := dom.Resume(); err != nil {
		return fmt.Errorf("failed to resume VM %s: %v", vmName, err)
	}
	return nil
}

// DeleteVM destroys and undefines a virtual machine by name
func DeleteVM(conn *libvirt.Connect, vmName string) error {
	dom, err := conn.LookupDomainByName(vmName)
	if err != nil {
		return fmt.Errorf("failed to find VM %s: %v", vmName, err)
	}
	if err := dom.Destroy(); err != nil {
		return fmt.Errorf("failed to destroy VM %s: %v", vmName, err)
	}
	if err := dom.Undefine(); err != nil {
		return fmt.Errorf("failed to undefine VM %s: %v", vmName, err)
	}
	return nil
}
