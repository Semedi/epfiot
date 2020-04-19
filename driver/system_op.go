package driver

import "os/exec"
import "log"
import "fmt"
import "strings"
import "bufio"

const METHOD = 0
const HOST = 1

var Location string
var Connection []string

func Initfs(l string, auth_host []string) {
	Location = l
	Connection = auth_host
	execute("mkdir", "-p", Location+"/base")
}

func Usb_info() [][]string {
	info, err := execute("lsusb")

	if err != nil {
		log.Fatal(err)
	}

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

func execute(parameters ...string) (string, error) {

	if Connection != nil {
		parameters = append(Connection, parameters...)
	}

	name, args := parameters[0], parameters[1:]

	out, err := exec.Command(name, args...).Output()
	if err != nil {
		msg := fmt.Sprintf("the next command has failed:\n   %s", strings.Join(parameters, " "))
		fmt.Println(msg)

		return "", err
	}

	return string(out), nil
}

func folder(user uint) string {
	usrlocation := fmt.Sprintf("%s%s%d", Location, "/user/", user)
	execute("mkdir", "-p", usrlocation)

	return usrlocation
}

func basefile(base string) string {
	return fmt.Sprintf("%s%s%s", Location, "/base/", base)
}

func Vmfile(user uint, name string) string {
	return fmt.Sprintf("%s/%s.qcow2", folder(user), name)
}

func Vmconfig(user uint, name string) string {
	return fmt.Sprintf("%s/%s.iso", folder(user), name)
}

func Erasefiles(user uint, name string) {
	execute("rm", Vmfile(user, name))
	execute("rm", Vmconfig(user, name))
}

func Copy_base(base string, user uint, name string) {

	file := basefile(base)

	_, err := execute("stat", file)

	if err != nil {
		log.Fatal(err)
	}

	execute("cp", file, Vmfile(user, name))
}
