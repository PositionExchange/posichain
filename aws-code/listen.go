package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

func listenMain() { // TODO: temp fix to unblock test
	ln, err := net.Listen("tcp", "localhost:9000")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		io.WriteString(conn, fmt.Sprint("Hello World\n", time.Now(), "\n"))
		conn.Close()
	}
}
