package main

import (
	"fmt"
	"log"

	libvirt "github.com/libvirt/libvirt-go"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
)

func main() {
	var drive uint
	drive = 0
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		log.Fatalf("failed to connect to qemu")
	}
	defer conn.Close()

	domcfg := &libvirtxml.Domain{
		Type: "kvm", Name: "demo", Memory: &libvirtxml.DomainMemory{Value: 4096, Unit: "MB", DumpCore: "on"},
		VCPU: &libvirtxml.DomainVCPU{Value: 1},
		CPU:  &libvirtxml.DomainCPU{Mode: "host-model"},
		Devices: &libvirtxml.DomainDeviceList{
			Disks: []libvirtxml.DomainDisk{
				{
					Source:  &libvirtxml.DomainDiskSource{File: &libvirtxml.DomainDiskSourceFile{File: "/home/semedi/Downloads/alpine.qcow"}},
					Target:  &libvirtxml.DomainDiskTarget{Dev: "hda", Bus: "ide"},
					Alias:   &libvirtxml.DomainAlias{Name: "ide0-0-0"},
					Address: &libvirtxml.DomainAddress{Drive: &libvirtxml.DomainAddressDrive{Controller: &drive, Bus: &drive, Target: &drive, Unit: &drive}}}}}}

	xml, err := domcfg.Marshal()
	if err != nil {
		panic(err)
	}

	domain, err := conn.DomainDefineXML(xml)
	if err != nil {
		panic(err)
	}

	createDomain, err := conn.DomainCreateXML(xml, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println(xml)
	fmt.Println(domain)
	fmt.Println(createDomain)
}
