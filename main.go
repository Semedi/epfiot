package main

import (
	d "github.com/semedi/epfiot/driver"
)

func main() {
	driver := new(d.Driver)

	driver.Init()
	driver.Start()
}
