package driver

import "log"

type Driver interface {
	Init()
	Create()
	List()
	Close()
}

type Controller struct {
    driver Driver
}

func New(d string, uri string) *Controller{
    var drv Driver;
	switch d {
        case "kvm":
            drv = New_kvm(uri)
        default:
            log.Fatalf("Unrecognized Driver")
	}

    c := &Controller{driver: drv}

    c.driver.Init()

    return  c
}


func (c *Controller) Start() {
	c.driver.Create()
	c.driver.List()
	c.driver.Close()
}
