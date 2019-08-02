package driver

import "log"

type Driver interface {
	Init()
	Create()
	List()
	Close()
}

type Controller struct {
    Handler Driver
}

func New(d string, uri string) *Controller{
    var drv Driver;
	switch d {
        case "kvm":
            drv = New_kvm(uri)
        default:
            log.Fatalf("Unrecognized Driver")
	}

    c := &Controller{Handler: drv}

    c.Handler.Init()

    return  c
}


func (c *Controller) Start() {
	c.Handler.Create()
	c.Handler.List()
	c.Handler.Close()
}
