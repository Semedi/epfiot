package driver

import "os"
import "os/exec"
import "log"
import "fmt"

var location string
var connection []string

func Initfs(l string, auth_host []string) {
	mode := int(0755)

	connection = auth_host
	location = l
	_ = os.MkdirAll(location+"/base", os.FileMode(mode))
}

func execute(parameters ...string) {

	if connection != nil {
		parameters = append(connection, parameters...)
	}

	name, parameters := parameters[0], parameters[1:]

	cmd := exec.Command(name, parameters...)
	err := cmd.Run()

	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
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
		log.Fatalf("file not exists something bad happened!")
	}

	execute("cp", file, Vmfile(user, name))
}
