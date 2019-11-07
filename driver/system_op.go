package driver

import "os"
import "os/exec"
import "log"
import "fmt"
import "strings"
import "bufio"

var location string
var connection []string

func Initfs(l string, auth_host []string) {
	mode := int(0755)

	connection = auth_host
	location = l
	_ = os.MkdirAll(location+"/base", os.FileMode(mode))
}

func Usb_info() [][]string {
	info := execute("lsusb")
	scanner := bufio.NewScanner(strings.NewReader(info))

	a := [][]string{}
	for scanner.Scan() {
		columns := strings.Fields(strings.Replace(scanner.Text(), ":", "", -1))
		bus := strings.TrimLeft(columns[1], "0")
		dev := strings.TrimLeft(columns[3], "0")
		info := strings.Join(columns[6:], " ")

		hdevice_info := []string{bus, dev, info}
		a = append(a, hdevice_info)
	}

	return a
}

func execute(parameters ...string) string {

	if connection != nil {
		parameters = append(connection, parameters...)
	}

	name, args := parameters[0], parameters[1:]

	out, err := exec.Command(name, args...).Output()
	if err != nil {
		msg := fmt.Sprintf("the next command has failed:\n   %s", strings.Join(parameters, " "))
		fmt.Println(msg)
		log.Fatal(err)
	}

	return string(out)
}

func folder(user uint) string {
	usrlocation := fmt.Sprintf("%s%s%d", location, "/user/", user)
	execute("mkdir", "-p", usrlocation)

	return usrlocation
}

func basefile(base string) string {
	return fmt.Sprintf("%s%s%s", location, "/base/", base)
}

func Vmfile(user uint, name string) string {
	return fmt.Sprintf("%s/%s.qcow2", folder(user), name)
}

func Copy_base(base string, user uint, name string) {

	file := basefile(base)

	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Fatalf(file)
		log.Fatalf("file not exists something bad happened!")
	}

	execute("cp", file, Vmfile(user, name))
}
