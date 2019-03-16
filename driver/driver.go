package driver

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Controller interface {
	Init()
	Create()
	List()
	Close()
}

type Driver struct {
	Controller
}
type conf struct {
	Host   string `yaml:"host"`
	Auth   string `yaml:"auth"`
	Driver string `yaml:"driver"`
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

func (d *Driver) Init() {
	var config conf
	config.get()

	switch config.Driver {
	case "kvm":
		d.Controller = New_kvm(config.uri())
	default:
		log.Fatalf("Unrecognized Driver")
	}
}

func (d *Driver) Start() {
	d.Controller.Create()
	d.Controller.List()
	d.Controller.Close()
}
