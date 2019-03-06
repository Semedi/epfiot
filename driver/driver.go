package driver

type Controller interface {
	Init()
	Create()
	List()
	Close()
}

type Driver struct {
	Controller
}

func (d *Driver) Start() {
	d.Controller.Init()
	d.Controller.Create()
	d.Controller.Close()
}
