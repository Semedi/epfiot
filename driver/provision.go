package driver

import (
	"fmt"
	"github.com/kless/osutil/user/crypt/sha512_crypt"
	"github.com/semedi/epfiot/core/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	Address         string   `yaml:",omitempty"`
	Netmask         string   `yaml:",omitempty"`
	Gateway         string   `yaml:",omitempty"`
	Dns_nameservers []string `yaml:",omitempty"`
}
type Nconfig struct {
	Type    string
	Name    string
	Subnets []Subnet
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
	Passwd      string `yaml:",omitempty"`
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

	if filepath.Base(filename) == "user-data" {
		d = append([]byte("#cloud-config\n"), d...)
	}

	ioutil.WriteFile(filename, d, 0644)
}

func Create_config(vmname string, c *model.ConfigInput) (*string, error) {

	if c == nil {
		return nil, nil
	}

	// mandatory network:
	t3 := Net{}
	t3.Network.Version = 1
	t3.Network.Config = []Nconfig{
		{
			Type: "physical",
			Name: "enp2s1",
			Subnets: []Subnet{
				//{Type: "dhcp"},
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

	t4 := Udata{}
	t4.Package_upgrade = false
	user_config := Uconfig{}

	if c.User != nil {

		user_config.Name = *c.User
		user_config.Groups = "wheel"
		user_config.Lock_passwd = false
		user_config.Shell = "/bin/bash"
		user_config.Sudo = []string{"ALL=(ALL) NOPASSWD:ALL"}

		if c.Password != nil {
			encrypt := sha512_crypt.New()
			hash, err := encrypt.Generate([]byte(*c.Password), []byte("$6$usesomesillystringforsalt"))
			if err != nil {
				return nil, err
			}
			user_config.Passwd = hash
		}

		user_config.Keys = []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCbA3KjpnI6gJLGoKru/iq1qhw+3y3B7Bqu5+MRVv3DTcc8wUocpePR8VH4MYomwBtEOki/13ZBRsl4zEkRorrRaITUlC/atUiUhI8u/8nFHGRkTFSMD3aFriysonfm5Ipg2arhpQMvbtDcd/oVcCpHnc1ifEOXHfm1Eyslhg8A91rLj2frFB5+Cqx1Gi+sfZ+L8ysA+Psrzf00Xn9EkDhLuomizzGSVc06dRxPb/Y2V+qHd7R2D/DxQXaaBGuPFCHS/bzLh4Y4Md5LKTVpZ3mTDD8ywdnTb1CEjGUyg1RAWXqqx+fbzVPGeAmkPgW0ZZpc1J3VycQfYQXvbzd53JynuljXZPsRT27+KXYnabCHPGsm4On6OUUgzkWZB/GpVVUw/xtTUPBgD5VUW3N850Z+sBRfqW7+uEqBybxwznp8klT+GSQ2vJC2R6bXOS7EJmU1iPTp7fRPC5zJiIGR+7ChSLNtTabWdO2FPGeGnZ9Mt1IJpyYvoknbTsWBBXwXu+hjxhT64XX9LBD+pebejIaWckOg51zX5kVgf+bNPvX1XSK9W2dOUTfcRkeWHwo7WqpAhbmXAkGju09Icjmk66drOhyTMmuPlEWmeWogcYGMizXtQK2GBQgnplEFH6/Hr1nmtKu1WLuwoiiVvluUg/bkr8DRwLaUT7KXr41WwLAT5Q=="}
	}

	network_config := "/tmp/network-config"
	user_data := "/tmp/user-data"
	meta_data := "/tmp/meta-data"
	cd := fmt.Sprintf("/tmp/%s.iso", vmname)

	t4.Users = []Uconfig{
		user_config,
	}

	t5 := Metadata{Id: vmname, Dsmode: "local"}

	write_config(t3, network_config)
	if _, err := os.Stat(network_config); os.IsNotExist(err) {
		return nil, err
	}
	write_config(t4, user_data)
	if _, err := os.Stat(user_data); os.IsNotExist(err) {
		return nil, err
	}
	write_config(t5, meta_data)
	if _, err := os.Stat(meta_data); os.IsNotExist(err) {
		return nil, err
	}

	cmd := exec.Command("genisoimage", "-output", cd, "-volid", "cidata", "-joliet", "-rock", user_data, meta_data, network_config)

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	// copy to remote
	if Connection[METHOD] == "ssh" {
		dest_cmd := fmt.Sprintf("%s %s:%s", cd, Connection[HOST], Location)

		cmd := exec.Command("scp", strings.Fields(dest_cmd)...)
		err := cmd.Run()
		if err != nil {
			return nil, err
		}
	}

	final_path := fmt.Sprintf("%s/%s", Location, filepath.Base(cd))

	return &final_path, nil
}
