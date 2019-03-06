package main

import (
	d "github.com/semedi/epfiot/driver"
)

func main() {
	driver := d.Driver{Controller: &d.Libvirt{Name: "test"}}

	driver.Start()
}
