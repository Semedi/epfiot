package main

import (
	"github.com/semedi/epfiot/driver"
	"github.com/semedi/epfiot/service"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type conf struct {
	Host    string `yaml:"host"`
	Auth    string `yaml:"auth"`
	Driver  string `yaml:"driver"`
	Storage string `yaml:"storage"`
}

func read_config() *conf {
	config := new(conf)

	yamlFile, err := ioutil.ReadFile("epfiot.conf")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return config
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
	config := read_config()
	driver.Initfs(config.Storage)

	server := service.New()
	controller := driver.New(config.Driver, config.uri())

	//controller.Handler.Destroy("polla")

	server.Run(controller)

	//controller.Start()
}
