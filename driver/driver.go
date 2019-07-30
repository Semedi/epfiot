package driver

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Driver interface {
	Init()
	Create()
	List()
	Close()
}

type Controller struct {
    driver Driver
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

func (c *Controller) Init() {
	var config conf
	config.get()

	switch config.Driver {
	case "kvm":
		c.driver = New_kvm(config.uri())
	default:
		log.Fatalf("Unrecognized Driver")
	}
}

func (c *Controller) Start() {
	c.driver.Create()
	c.driver.List()
	c.driver.Close()
}
