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

			var response string
			if len(strings) > 0 {
				if strings[0] == "PING" {
					response = "+PONG\r\n"
				} else if strings[0] == "ECHO" && len(strings) > 1 {
					response = "+" + strings[1:][0]
					for _, s := range strings[2:] {
						response += " " + s
					}
					response += "\r\n"
				} else {
					response = "+PONG\r\n"
				}
			}

			write, err := conn.Write([]byte(response))
			if err != nil {
				fmt.Println("Error write: ", err.Error())
				return
			}
			fmt.Println("it was written: ", write)
		}
	}
}
