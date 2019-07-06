package main

import (
	d "github.com/semedi/epfiot/driver"
    "github.com/semedi/epfiot/service"
)

func main() {
	driver := new(d.Driver)

    server := service.New()
    server.Run()

	driver.Init()
	driver.Start()
}
