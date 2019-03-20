package service

import (
	"fmt"
)

func main() {
	driver := new(d.Driver)

	driver.Init()
	driver.Start()
}
