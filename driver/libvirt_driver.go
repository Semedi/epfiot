package driver

import (
	"log"
	"fmt"

	libvirt "github.com/libvirt/libvirt-go"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
    "github.com/semedi/epfiot/core/model"
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

func domain_def(vcpu int) libvirtxml.Domain {
	domcfg := libvirtxml.Domain{
		Type: "kvm",
		OS: &libvirtxml.DomainOS{
			Type: &libvirtxml.DomainOSType{
				Type: "hvm",
			},
		},
		VCPU: &libvirtxml.DomainVCPU{
			Placement: "static",
			Value:     vcpu,
		},
		CPU: &libvirtxml.DomainCPU{Mode: "host-model"},
		Devices: &libvirtxml.DomainDeviceList{
			Graphics: []libvirtxml.DomainGraphic{
				{
					Spice: &libvirtxml.DomainGraphicSpice{
						AutoPort: "yes",
					},
				},
			},
		},
	}

	return domcfg
}

func (l *Libvirt) Init() {
}

func (l *Libvirt) Close() {
	defer l.conn.Close()
}



func (l *Libvirt) get(query string) (bool, *libvirt.Domain){
    doms := l.get_all()

	for _, dom := range doms {
		name, err := dom.GetName()
		if err == nil {
			fmt.Printf("  %s\n", name)
		}

        if name == query{
            return true, &dom
        }

		dom.Free()
	}

    return false, nil
}

func (l *Libvirt) Destroy(query string) (error){
    r, dom := l.get(query)

    if  r == true{
        err := dom.Destroy()

        if err != nil {
            return err
        }
    }

    return nil
}


// take care of free the c pointer after calling this method
func (l *Libvirt) get_all()([]libvirt.Domain) {
	da, err := l.conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	if err != nil {
		panic(err)
	}

	di, err := l.conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_INACTIVE)
	if err != nil {
		panic(err)
	}

    return append(da, di...)
}

func (l *Libvirt) List() {
    doms := l.get_all()

	fmt.Printf("%d running domains:\n", len(doms))
	for _, dom := range doms {
		name, err := dom.GetName()
		if err == nil {
			fmt.Printf("  %s\n", name)
		}
		dom.Free()
	}
}

//1 : CONNECT_LIST_DOMAINS_ACTIVE
//2 : CONNECT_LIST_DOMAINS_INACTIVE
//4 : CONNECT_LIST_DOMAINS_PERSISTENT
//16: CONNECT_LIST_DOMAINS_RUNNING
//32: CONNECT_LIST_DOMAINS_PAUSED
//64: CONNECT_LIST_DOMAINS_SHUTOFF
func (l *Libvirt) Listt() {
    l.List()
}

//Source: &libvirtxml.DomainInterfaceSource{
//    Bridge: &libvirtxml.DomainInterfaceSourceBridge{
//        Bridge: "epfiot_net",
//    },
//},
func setDevices(d *libvirtxml.Domain, ilocation string) {
	d.Devices.Interfaces = []libvirtxml.DomainInterface{
		{
			Source: &libvirtxml.DomainInterfaceSource{
				Network: &libvirtxml.DomainInterfaceSourceNetwork{
					Network: "epfiot-vm",
				},
			},
		},
	}
	d.Devices.Disks = []libvirtxml.DomainDisk{
		{
			Source: &libvirtxml.DomainDiskSource{File: &libvirtxml.DomainDiskSourceFile{File: ilocation}},
			Driver: &libvirtxml.DomainDiskDriver{Name: "qemu", Type: "qcow2"},
			Target: &libvirtxml.DomainDiskTarget{Dev: "hda", Bus: "virtio"},
		},
	}
}
func setMemory(d *libvirtxml.Domain, m int) {
    d.Memory= &libvirtxml.DomainMemory{
                Value:    (uint)(m),
                Unit:     "MB",
                DumpCore: "on",
	        }
}

func (l *Libvirt) Create(vm model.Vm, uid uint) {
	domcfg := domain_def(vm.Vcpu)

	domcfg.Name = vm.Name
	setDevices(&domcfg, Vmfile(uid, vm.Name))
    setMemory(&domcfg, vm.Memory)

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

