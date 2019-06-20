package main

import (
	d "github.com/semedi/epfiot/driver"
    "github.com/semedi/epfiot/service"
)

func main() {
    service.Run()
	driver := new(d.Driver)

	driver.Init()
	driver.Start()
}
