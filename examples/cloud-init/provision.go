package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
	"os"
	"os/exec"
)

//----------------------------
// network-config
//----------------------------
type Net struct {
	Network struct {
		Version int
		Config  []Nconfig
	}
}
type Subnet struct {
	Type            string
	Address         string
	Netmask         string
	Gateway         string
	Dns_nameservers []string
}
type Nconfig struct {
	Type        string
	Name        string
	Mac_address string
	Subnets     []Subnet
}

//----------------------------
// user-data
//----------------------------

type Udata struct {
	Package_upgrade bool
	Users           []Uconfig
}

type Uconfig struct {
	Name        string
	Groups      string
	Lock_passwd bool
	Passwd      string
	Shell       string
	Sudo        []string
	Keys        []string `yaml:"ssh-authorized-keys,omitempty"`
}

//----------------------------
// meta-data
//----------------------------
type Metadata struct {
	Id     string `yaml:"instance-id"`
	Dsmode string
}

func write_config(t interface{}, filename string) {

	d, err := yaml.Marshal(&t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- yaml dump:\n%s\n\n", string(d))

	if filename == "user-data" {
		d = append([]byte("#cloud-config\n"), d...)
	}

	ioutil.WriteFile(filename, d, 0644)
}

func main() {
	t3 := Net{}
	t3.Network.Version = 1
	t3.Network.Config = []Nconfig{
		{
			Type:        "physical",
			Name:        "enp2s1",
			Mac_address: "52:54:00:27:ae:db",
			Subnets: []Subnet{
				{
					Type:            "static",
					Address:         "10.0.0.2",
					Netmask:         "255.255.255.0",
					Gateway:         "10.0.0.1",
					Dns_nameservers: []string{"10.0.0.1", "8.8.8.8"},
				},
			},
		},
	}
	key := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCbA3KjpnI6gJLGoKru/iq1qhw+3y3B7Bqu5+MRVv3DTcc8wUocpePR8VH4MYomwBtEOki/13ZBRsl4zEkRorrRaITUlC/atUiUhI8u/8nFHGRkTFSMD3aFriysonfm5Ipg2arhpQMvbtDcd/oVcCpHnc1ifEOXHfm1Eyslhg8A91rLj2frFB5+Cqx1Gi+sfZ+L8ysA+Psrzf00Xn9EkDhLuomizzGSVc06dRxPb/Y2V+qHd7R2D/DxQXaaBGuPFCHS/bzLh4Y4Md5LKTVpZ3mTDD8ywdnTb1CEjGUyg1RAWXqqx+fbzVPGeAmkPgW0ZZpc1J3VycQfYQXvbzd53JynuljXZPsRT27+KXYnabCHPGsm4On6OUUgzkWZB/GpVVUw/xtTUPBgD5VUW3N850Z+sBRfqW7+uEqBybxwznp8klT+GSQ2vJC2R6bXOS7EJmU1iPTp7fRPC5zJiIGR+7ChSLNtTabWdO2FPGeGnZ9Mt1IJpyYvoknbTsWBBXwXu+hjxhT64XX9LBD+pebejIaWckOg51zX5kVgf+bNPvX1XSK9W2dOUTfcRkeWHwo7WqpAhbmXAkGju09Icjmk66drOhyTMmuPlEWmeWogcYGMizXtQK2GBQgnplEFH6/Hr1nmtKu1WLuwoiiVvluUg/bkr8DRwLaUT7KXr41WwLAT5Q=="

	t4 := Udata{}
	t4.Package_upgrade = false
	t4.Users = []Uconfig{
		{
			Name:        "semedi",
			Groups:      "wheel",
			Lock_passwd: false,
			Passwd:      "$1$SaltSalt$pdqtneIkXPJIjowbO8gt7/",
			Shell:       "/bin/bash",
			Sudo:        []string{"ALL=(ALL) NOPASSWD:ALL"},
			Keys:        []string{key},
		},
	}

	t5 := Metadata{Id: "myVM", Dsmode: "local"}

	os.Remove("network-config")
	write_config(t3, "network-config")
	if _, err := os.Stat("network-config"); os.IsNotExist(err) {
		log.Fatalf("file not exists something bad happened!")
	}
	os.Remove("user-data")
	write_config(t4, "user-data")
	if _, err := os.Stat("user-data"); os.IsNotExist(err) {
		log.Fatalf("file not exists something bad happened!")
	}
	os.Remove("meta-data")
	write_config(t5, "meta-data")
	if _, err := os.Stat("meta-data"); os.IsNotExist(err) {
		log.Fatalf("file not exists something bad happened!")
	}

	cd := "config.iso"
	cmd := exec.Command("genisoimage", "-output", cd, "-volid", "cidata", "-joliet", "-rock", "user-data", "meta-data", "network-config")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}
