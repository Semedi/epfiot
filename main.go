package main

import (
	d "github.com/semedi/epfiot/driver"
)

func main() {

	driver := d.Driver{Controller: d.New_kvm("qemu:///system")}

	driver.Start()
}
