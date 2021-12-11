package main

import (
    "fmt"
    "io"
    "net"
)

func PortForward(input_host string, output_host string) {
    fmt.Printf("Forwarding %s -> %s", input_host, output_host)
    ln, err := net.Listen("tcp", input_host)
    if err != nil {
        panic(err)
    }

    for {
        conn, err := ln.Accept()
        if err != nil {
            panic(err)
        }

        go handleRequest(conn, output_host)
    }
	return
}

func handleRequest(conn net.Conn, output_host string) {
    fmt.Println("new client")

    proxy, err := net.Dial("tcp", output_host)
    if err != nil {
        panic(err)
    }

    fmt.Println("proxy connected")
    go copyIO(conn, proxy)
    go copyIO(proxy, conn)
}

func copyIO(src, dest net.Conn) {
    defer src.Close()
    defer dest.Close()
    io.Copy(src, dest)
}
