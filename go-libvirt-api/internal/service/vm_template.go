package service

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

func GenerateVMXML(name string, memory int, vcpu int, disk string, network string) string {
	return fmt.Sprintf(`
<domain type='kvm'>
  <name>%s</name>
  <memory unit='MiB'>%d</memory>
  <vcpu>%d</vcpu>
  <os>
    <type arch='x86_64'>hvm</type>
  </os>
  <devices>
    <disk type='file' device='disk'>
      <source file='%s'/>
      <target dev='vda' bus='virtio'/>
    </disk>
    <interface type='network'>
      <source network='%s'/>
      <model type='virtio'/>
    </interface>
  </devices>
</domain>`, name, memory, vcpu, disk, network)
}

func CreateVM(conn *libvirt.Connect, xml string) error {
	dom, err := conn.DomainDefineXML(xml)
	if err != nil {
		return fmt.Errorf("failed to define domain: %v", err)
	}
	defer dom.Free()

	if err := dom.Create(); err != nil {
		return fmt.Errorf("failed to create domain: %v", err)
	}
	return nil
}
