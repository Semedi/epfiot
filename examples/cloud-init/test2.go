package main

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v2"
)

var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`

// Note: struct fields must be public in order for unmarshal to
// correctly populate the data.
type T struct {
	A string
	B struct {
		RenamedC int
		D        []int
	}
}

type M struct {
	A string
	C string
	D string
}

func write_config(t interface{}) {

	d, err := yaml.Marshal(&t)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- t dump:\n%s\n\n", string(d))
}

func main() {
	t := T{}
	t2 := M{}

	t.A = "prolla"
	t.B.D = []int{1, 2}

	t2.A = "colega"
	t2.C = "eyy"

	write_config(t)
	write_config(t2)

}
