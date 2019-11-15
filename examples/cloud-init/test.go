package main

import "os"
import "os/exec"
import "log"

func main() {
	file := "user-data"

	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Fatalf("file not exists something bad happened!")
	}

	cd := "config.iso"
	cmd := exec.Command("genisoimage", "-output", cd, "-volid", "cidata", "-joliet", "-rock", "user-data", "meta-data", "network-config")
	err := cmd.Run()

	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

}
