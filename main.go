package main

import (
	d "github.com/semedi/epfiot/driver"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type conf struct {
	Host string `yaml:"host"`
	Auth string `yaml:"auth"`
}

func (c *conf) get() *conf {

	yamlFile, err := ioutil.ReadFile("epfiot.conf")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}

func (c conf) uri() string {
	var u string

	if c.Host == "localhost" {
		u = "qemu:///system"
	} else {
		u = "qemu+" + c.Auth + "://" + c.Host + "/system"
	}

	return u
}

func main() {
	var config conf
	config.get()

	driver := d.Driver{Controller: d.New_kvm(config.uri())}

	driver.Start()
}
