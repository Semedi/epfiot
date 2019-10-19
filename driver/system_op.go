package driver

import "os"
import "os/exec"
import "log"
import "fmt"

var location string

func Initfs(l string) {
	location = l
}

func folder(user uint) string {
	usrlocation := fmt.Sprintf("%s%s%d", location, "/user/", user)
	cmd := exec.Command("mkdir", "-p", usrlocation)

	err := cmd.Run()

	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	return usrlocation
}

func basefile(base string) string {
	return fmt.Sprintf("%s%s%s", location, "/base/", base)
}

func Vmfile(user uint, name string) string {
	return fmt.Sprintf("%s/%s.qcow2", folder(user), name)
}

func Copy_base(base string, user uint, name string) error {

	file := basefile(base)

	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Fatalf("file not exists something bad happened!")
	}

	cmd := exec.Command("cp", file, Vmfile(user, name))
	err := cmd.Run()

	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return err
	}

	return nil

}
