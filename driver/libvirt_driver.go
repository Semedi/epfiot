package driver

import (
	"fmt"
	"log"

	libvirt "github.com/libvirt/libvirt-go"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
)

type Libvirt_driver struct {
	Name string
}

func Create() {
	var drive uint
	drive = 0

	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		log.Fatalf("failed to connect to qemu")
	}
	defer conn.Close()

	domcfg := &libvirtxml.Domain{
		Name: "demo01",
		Type: "kvm",
		OS: &libvirtxml.DomainOS{
			Type: &libvirtxml.DomainOSType{
				Type: "hvm",
			},
		},
		Memory: &libvirtxml.DomainMemory{
			Value:    2048,
			Unit:     "MB",
			DumpCore: "on"},
		VCPU: &libvirtxml.DomainVCPU{Value: 1},
		CPU:  &libvirtxml.DomainCPU{Mode: "host-model"},
		Devices: &libvirtxml.DomainDeviceList{
			Disks: []libvirtxml.DomainDisk{
				{
					Source: &libvirtxml.DomainDiskSource{File: &libvirtxml.DomainDiskSourceFile{File: "/home/semedi/Downloads/alpine.qcow"}},
					Target: &libvirtxml.DomainDiskTarget{Dev: "hda", Bus: "ide"},
					Alias:  &libvirtxml.DomainAlias{Name: "ide0-0-0"},
					Address: &libvirtxml.DomainAddress{
						Drive: &libvirtxml.DomainAddressDrive{
							Controller: &drive, Bus: &drive, Target: &drive, Unit: &drive},
					},
				}}}}

	xml, err := domcfg.Marshal()
	if err != nil {
		panic(err)
	}

	fmt.Println(xml)

	domain, err := conn.DomainDefineXML(xml)
	if err != nil {
		panic(err)
	}

	fmt.Println(domain)

	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d running domains:\n", len(doms))
	for _, dom := range doms {
		name, err := dom.GetName()
		if err == nil {
			fmt.Printf("  %s\n", name)
		}
		dom.Free()
	}
}
