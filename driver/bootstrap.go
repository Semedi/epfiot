package driver

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

/*
* STRUCTS & ENUMS
*
 */
const (
	SERVER   = "server"
	ENDPOINT = "endpoint"
)

type tendpoint struct {
	NAME   string
	DELETE string
	SERVER string
}

type tserver struct {
	ID        string
	URI       string
	BOOTSTRAP string
	LIFETIME  string
	SECURITY  string
}

var server = &tserver{
	ID:        "id",
	URI:       "uri",
	BOOTSTRAP: "bootrstap",
	LIFETIME:  "lifetime",
	SECURITY:  "security",
}

var endpoint = &tendpoint{
	NAME:   "Name",
	DELETE: "Delete",
	SERVER: "Server",
}

type petition struct {
	header   string
	commands []command
}

type command struct {
	key   string
	value string
}

// petition renders examples
// echo -n "server|id=2,uri=coap://leshan.eclipse.org:5683,bootstrap=no,lifetime=300,security=NoSec" | nc -4u -w1 localhost 5400
// echo -n "endpoint|Delete=/0,Delete=/1,Server=2" | nc -4u -w1 localhost 5400

func (c command) out() string {
	return fmt.Sprintf("%s=%s", c.key, c.value)
}

func (p petition) out() string {
	var list []string

	for _, command := range p.commands {
		list = append(list, command.out())
	}

	chain := fmt.Sprintf("%s|%s", p.header, strings.Join(list, ","))
	return strings.TrimSpace(chain)
}

func new_petition(header string) *petition {

	var commands []command

	return &petition{header: header, commands: commands}
}

func new_command(key string, value string) command {

	return command{key: key, value: value}
}

func (p *petition) add(c command) {
	p.commands = append(p.commands, c)
}

func bootstrap_register(request petition) error {
	p := make([]byte, 2048)
	conn, err := net.Dial("udp", "127.0.0.1:5400")
	if err != nil {
		return err
	}

	fmt.Fprintf(conn, request.out())
	_, err = bufio.NewReader(conn).Read(p)
	if err == nil {
		fmt.Printf("%s\n", p)
	} else {
		return err
	}
	conn.Close()

	return nil
}

func server_request() error {
	pet := new_petition(SERVER)
	pet.add(new_command(server.ID, "2"))
	pet.add(new_command(server.URI, "coap://localhost:5683"))
	pet.add(new_command(server.BOOTSTRAP, "no"))
	pet.add(new_command(server.LIFETIME, "300"))
	pet.add(new_command(server.SECURITY, "NoSec"))

	return bootstrap_register(*pet)
}

func thing_request() error {
	pet := new_petition(ENDPOINT)
	pet.add(new_command(endpoint.DELETE, "/0"))
	pet.add(new_command(endpoint.DELETE, "/1"))
	pet.add(new_command(endpoint.SERVER, "2"))

	return bootstrap_register(*pet)
}
