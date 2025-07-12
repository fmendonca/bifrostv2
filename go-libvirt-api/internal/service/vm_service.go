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

	if err := dom.Create(); err != nil {
		return fmt.Errorf("failed to start domain %s: %v", name, err)
	}
	return nil
}

func StopVM(client *libvirtclient.Client, name string) error {
	dom, err := client.Conn().LookupDomainByName(name)
	if err != nil {
		return fmt.Errorf("failed to lookup domain %s: %v", name, err)
	}
	defer dom.Free()

	if err := dom.Shutdown(); err != nil {
		return fmt.Errorf("failed to shutdown domain %s: %v", name, err)
	}
	return nil
}

func RebootVM(client *libvirtclient.Client, name string) error {
	dom, err := client.Conn().LookupDomainByName(name)
	if err != nil {
		return fmt.Errorf("failed to lookup domain %s: %v", name, err)
	}
	defer dom.Free()

	// The official libvirt binding doesn't expose DOMAIN_REBOOT_DEFAULT, pass 0
	if err := dom.Reboot(0); err != nil {
		return fmt.Errorf("failed to reboot domain %s: %v", name, err)
	}
	return nil
}

func PauseVM(client *libvirtclient.Client, name string) error {
	dom, err := client.Conn().LookupDomainByName(name)
	if err != nil {
		return fmt.Errorf("failed to lookup domain %s: %v", name, err)
	}
	defer dom.Free()

	if err := dom.Suspend(); err != nil {
		return fmt.Errorf("failed to pause domain %s: %v", name, err)
	}
	return nil
}

func ResumeVM(client *libvirtclient.Client, name string) error {
	dom, err := client.Conn().LookupDomainByName(name)
	if err != nil {
		return fmt.Errorf("failed to lookup domain %s: %v", name, err)
	}
	defer dom.Free()

	if err := dom.Resume(); err != nil {
		return fmt.Errorf("failed to resume domain %s: %v", name, err)
	}
	return nil
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
	if err := dom.Undefine(); err != nil {
		return fmt.Errorf("failed to undefine domain %s: %v", name, err)
	}
	return nil
}
