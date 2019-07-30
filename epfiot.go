package main

import (
	d "github.com/semedi/epfiot/driver"
    //"github.com/semedi/epfiot/service"
)

func main() {
	driver := new(d.Controller)

    driver.Init()
    //server := service.New()

    //server.Run(driver)

    driver.Start()
}
