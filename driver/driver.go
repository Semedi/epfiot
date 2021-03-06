package driver

import (
	"github.com/semedi/epfiot/core/model"
	"log"
)

type Provider interface {
	Init()
	Create(vm model.Vm, uid uint, config_path *string)
	AttachDevice(vm model.Vm, dev model.Hostdev) error
	DetachDevice(vm model.Vm, dev model.Hostdev) error
	List()
	Listt()
	Shutdown(query string) error
	Destroy(query string) error
	ForceOFF(query string) error
	PowerOn(query string) error
	Close()
}

func New(d string, uri string) Provider {
	var drv Provider
	switch d {
	case "kvm":
		drv = New_kvm(uri)
	default:
		log.Fatalf("Unrecognized Driver")
	}

	drv.Init()

	return drv
}
