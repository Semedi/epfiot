package driver

import (
	"github.com/semedi/epfiot/core/model"
	"log"
)

type Driver interface {
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

type Controller struct {
	Handler Driver
}

func New(d string, uri string) *Controller {
	var drv Driver
	switch d {
	case "kvm":
		drv = New_kvm(uri)
	default:
		log.Fatalf("Unrecognized Driver")
	}

	c := &Controller{Handler: drv}

	c.Handler.Init()

	return c
}

func (c *Controller) Start() {
	c.Handler.Close()
}
