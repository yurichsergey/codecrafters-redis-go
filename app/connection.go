package main

import (
	"fmt"
	"net"
)

func handleConnection(processor *Processor, conn net.Conn) {
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
		inputStrings, err := parseString(string(buf[:n]))
		if err != nil {
			fmt.Println("Error parsing string: ", err)
			return
		}

		fmt.Printf("We got: %s\n", inputStrings)

		if n > 0 {
			response := processor.ProcessCommand(inputStrings)

			write, err := conn.Write([]byte(response))
			if err != nil {
				fmt.Println("Error write: ", err.Error())
				return
			}
			fmt.Println("it was written (bytes): ", write, response)
		}
	}
}
