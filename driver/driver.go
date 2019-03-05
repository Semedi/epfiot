package driver

type Driver interface {
	init()
	create()
}

type Controller struct {
	d Driver
}
