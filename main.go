package main

import (
    libvirt "github.com/libvirt/libvirt-go"
    "fmt"
)

func main() {
    conn, err := libvirt.NewConnect("qemu:///system")
    if err != nil {
    }
    defer conn.Close()

    doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
    if err != nil {
    }

    fmt.Printf("%d running domains:\n", len(doms))
    for _, dom := range doms {
        name, err := dom.GetName()
        if err == nil {
            fmt.Printf("  %s\n", name)
        }
        dom.Free()
    }
}

