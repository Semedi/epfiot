package driver

import (
	"errors"
	"fmt"
	"log"
	"strconv"

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

//TODO:
//	make nvram uefi paths configurable
//	make Machine type (pc-q35) configurable
func domain_def(name string, vcpu int) libvirtxml.Domain {
	nvram := fmt.Sprintf("/var/lib/libvirt/qemu/nvram/%s_uefi_VARS.fd", name)

	domcfg := libvirtxml.Domain{
		Name: name,
		Type: "kvm",
		OS: &libvirtxml.DomainOS{
			Type: &libvirtxml.DomainOSType{
				Arch:    "x86_64",
				Type:    "hvm",
				Machine: "pc-q35-3.1",
			},
			Loader: &libvirtxml.DomainLoader{
				Path:     "/usr/share/OVMF/OVMF_CODE.fd",
				Type:     "pflash",
				Readonly: "yes",
			},
			NVRam: &libvirtxml.DomainNVRam{
				NVRam: nvram,
			},
			BootDevices: []libvirtxml.DomainBootDevice{
				{
					Dev: "hd",
				},
			},
		},
		Features: &libvirtxml.DomainFeatureList{
			ACPI: &libvirtxml.DomainFeature{},
			APIC: &libvirtxml.DomainFeatureAPIC{},
			VMPort: &libvirtxml.DomainFeatureState{
				State: "off",
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

/*
* Private
 */
func state(dom libvirt.Domain) (libvirt.DomainState, error) {

	info, err := dom.GetInfo()
	if err != nil {
		return -1, err
	}

	return info.State, nil
}

func (l *Libvirt) get(query string) (bool, *libvirt.Domain) {
	doms := l.get_all()

	for _, dom := range doms {
		name, err := dom.GetName()
		if err == nil {
			fmt.Printf("  %s\n", name)
		}

		if name == query {

			return true, &dom
		}

		dom.Free()
	}

	return false, nil
}

func (l *Libvirt) PowerOn(query string) error {
	r, dom := l.get(query)

	if r != true {
		return errors.New("vm does not exist")
	}

	Active, err := dom.IsActive()
	if err != nil {
		return err
	}

	if Active {
		return errors.New("vm active")
	}

	err = dom.Create()

	if err != nil {
		return err
	}

	return nil
}

func (l *Libvirt) Shutdown(query string) error {
	r, dom := l.get(query)

	if r != true {
		return errors.New("vm does not exist")
	}
	Active, err := dom.IsActive()
	if err != nil {
		return err
	}

	if !Active {
		return errors.New("vm is not active")
	}

	err = dom.Shutdown()

	if err != nil {
		return err
	}

	return nil
}

func (l *Libvirt) Update(vm *model.Vm) error {
	r, dom := l.get(vm.Name)
	if r != true {
		errors.New("domain not found!")
	}

	if s, err := state(*dom); s == libvirt.DOMAIN_RUNNING {
		addresses, err := dom.ListAllInterfaceAddresses(0)

		if err != nil {
			return err
		}

		if len(addresses) > 0 {
			vm.Ip = addresses[0].Addrs[0].Addr
		}

		vm.State = "RUNNING"

		return nil
	} else {
		return err
	}
}

func (l *Libvirt) Destroy(query string) error {

	r, dom := l.get(query)
	if r != true {
		errors.New("domain not found!")
	}

	info, err := dom.GetInfo()
	if err != nil {
		return err
	}

	//SHUTOFF
	if info.State == 5 {

		// delete nvram file
		err := dom.UndefineFlags(7)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Libvirt) ForceDestroy(query string) error {

	err := l.ForceOFF(query)
	if err != nil {
		return err
	}

	err = l.Destroy(query)
	if err != nil {
		return err
	}

	return nil
}

func (l *Libvirt) ForceOFF(query string) error {

	r, dom := l.get(query)
	if r != true {
		errors.New("domain not found!")
	}

	info, err := dom.GetInfo()
	if err != nil {
		return err
	}

	//RUNNING
	if info.State == 1 {
		err := dom.Destroy()
		if err != nil {
			return err
		}
	}

	return nil
}

// take care of free the c pointer after calling this method
func (l *Libvirt) get_all() []libvirt.Domain {
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
func xmldevice(hdev model.Hostdev) (string, error) {

	bus, _ := strconv.ParseUint(hdev.Bus, 10, 32)
	dev, _ := strconv.ParseUint(hdev.Device, 10, 32)

	ubus := uint(bus)
	udev := uint(dev)

	usb := &libvirtxml.DomainAddressUSB{
		Bus:    &ubus,
		Device: &udev,
	}

	devcfg := libvirtxml.DomainHostdev{
		Managed: "yes",
		SubsysUSB: &libvirtxml.DomainHostdevSubsysUSB{
			Source: &libvirtxml.DomainHostdevSubsysUSBSource{
				Address: usb},
		},
	}
	xml, err := devcfg.Marshal()
	if err != nil {
		return "", err
	}

	return xml, nil
}

func (l *Libvirt) DetachDevice(vm model.Vm, hdev model.Hostdev) error {
	r, dom := l.get(vm.Name)
	if r != true {
		return errors.New("domain not found!")
	}

	xml, err := xmldevice(hdev)
	if err != nil {
		return err
	}
	err = dom.DetachDevice(xml)
	if err != nil {
		return err
	}

	return nil
}

func (l *Libvirt) AttachDevice(vm model.Vm, hdev model.Hostdev) error {
	r, dom := l.get(vm.Name)
	if r != true {
		return errors.New("domain not found!")
	}

	xml, err := xmldevice(hdev)
	if err != nil {
		return err
	}

	err = dom.AttachDevice(xml)
	if err != nil {
		return err
	}

	return nil
}

func unmanaged_interface() []libvirtxml.DomainInterface {
	return []libvirtxml.DomainInterface{
		{
			Source: &libvirtxml.DomainInterfaceSource{
				Bridge: &libvirtxml.DomainInterfaceSourceBridge{
					Bridge: "epfiot-net",
				},
			},
		},
	}
}

func managed_interface() []libvirtxml.DomainInterface {
	return []libvirtxml.DomainInterface{
		{

			Source: &libvirtxml.DomainInterfaceSource{
				Network: &libvirtxml.DomainInterfaceSourceNetwork{
					Network: "epfiot_managed",
				},
			},
		},
	}
}

func setDevices(d *libvirtxml.Domain, ilocation string, vm model.Vm, config_path *string) {
	d.Devices.Interfaces = managed_interface()

	if vm.Dev != nil {

		d.Devices.Hostdevs = []libvirtxml.DomainHostdev{}
		for _, hdev := range vm.Dev {
			bus, _ := strconv.ParseUint(hdev.Bus, 10, 32)
			dev, _ := strconv.ParseUint(hdev.Device, 10, 32)

			ubus := uint(bus)
			udev := uint(dev)

			usb := &libvirtxml.DomainAddressUSB{
				Bus:    &ubus,
				Device: &udev,
			}

			d.Devices.Hostdevs = append(d.Devices.Hostdevs, libvirtxml.DomainHostdev{
				Managed: "yes",
				SubsysUSB: &libvirtxml.DomainHostdevSubsysUSB{
					Source: &libvirtxml.DomainHostdevSubsysUSBSource{
						Address: usb},
				},
			})

		}
	}

	d.Devices.Disks = []libvirtxml.DomainDisk{
		{
			Source: &libvirtxml.DomainDiskSource{File: &libvirtxml.DomainDiskSourceFile{File: ilocation}},
			Driver: &libvirtxml.DomainDiskDriver{Name: "qemu", Type: "qcow2"},
			Target: &libvirtxml.DomainDiskTarget{Dev: "hda", Bus: "virtio"},
		},
	}

	if config_path != nil {
		d.Devices.Disks = append(d.Devices.Disks, libvirtxml.DomainDisk{
			Device:   "cdrom",
			Source:   &libvirtxml.DomainDiskSource{File: &libvirtxml.DomainDiskSourceFile{File: *config_path}},
			Driver:   &libvirtxml.DomainDiskDriver{Name: "qemu", Type: "raw"},
			Target:   &libvirtxml.DomainDiskTarget{Dev: "sda", Bus: "sata"},
			ReadOnly: &libvirtxml.DomainDiskReadOnly{},
		})

	}

}

func setMemory(d *libvirtxml.Domain, m int) {
	d.Memory = &libvirtxml.DomainMemory{
		Value:    (uint)(m),
		Unit:     "MB",
		DumpCore: "on",
	}
}
func (l *Libvirt) Create(vm model.Vm, uid uint, config_path *string) {
	domcfg := domain_def(vm.Name, vm.Vcpu)

	setDevices(&domcfg, Vmfile(uid, vm.Name), vm, config_path)
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
