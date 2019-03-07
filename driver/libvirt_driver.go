package driver

import (
	"fmt"
	"log"

	libvirt "github.com/libvirt/libvirt-go"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
)

type Libvirt struct {
	Name string
	conn *libvirt.Connect
}

func (l *Libvirt) Init() {
	fmt.Println("init")
}

func (l *Libvirt) Close() {
	defer l.conn.Close()
}

func (l *Libvirt) List() {
	doms, err := l.conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
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

func (l *Libvirt) Create() {
	var drive uint
	drive = 0

	conn, err := libvirt.NewConnect("qemu:///system")

	if err != nil {
		log.Fatalf("failed to connect to qemu")
	}
	l.conn = conn

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
				},
			},
		},
	}

	xml, err := domcfg.Marshal()
	if err != nil {
		panic(err)
	}

	fmt.Println(xml)

	domain, err := l.conn.DomainDefineXML(xml)
	if err != nil {
		panic(err)
	}

	fmt.Println(domain.GetName())
}
