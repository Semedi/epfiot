package driver

import (
	"fmt"
	"log"

	libvirt "github.com/libvirt/libvirt-go"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
)

type Libvirt struct {
	conn *libvirt.Connect
}

func New_kvm(c string) *Libvirt {
	l := new(Libvirt)
	conn, err := libvirt.NewConnect(c)

	if err != nil {
		log.Fatalf("failed to connect to qemu")
	}

	l.conn = conn

	return l
}

func domain_def() libvirtxml.Domain {
	domcfg := libvirtxml.Domain{
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
			Interfaces: []libvirtxml.DomainInterface{
				{
					Source: &libvirtxml.DomainInterfaceSource{
						Bridge: &libvirtxml.DomainInterfaceSourceBridge{
							Bridge: "epfiot_net",
						},
					},
				},
			},
			Graphics: []libvirtxml.DomainGraphic{
				{
					Spice: &libvirtxml.DomainGraphicSpice{
						AutoPort: "yes",
					},
				},
			},
			Disks: []libvirtxml.DomainDisk{
				{
					Source: &libvirtxml.DomainDiskSource{File: &libvirtxml.DomainDiskSourceFile{File: "/home/semedi/Downloads/vm.qcow2"}},
					Driver: &libvirtxml.DomainDiskDriver{Name: "qemu", Type: "qcow2"},
					Target: &libvirtxml.DomainDiskTarget{Dev: "hda", Bus: "virtio"},
				},
			},
		},
	}

	return domcfg
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
	domcfg := domain_def()
	domcfg.Name = "demo01"

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
