package main

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection: ", err.Error())
		}
	}(conn)

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}
		strings, err := parseString(string(buf[:n]))
		if err != nil {
			fmt.Println("Error parsing string: ", err)
			return
		}

		fmt.Printf("We got: %s\n", strings)

		if n > 0 {

			write, err := conn.Write([]byte("+PONG\r\n"))
			if err != nil {
				fmt.Println("Error write: ", err.Error())
				return
			}
			fmt.Println("it was written: ", write)
		}
	}
}
