package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	// connect to the server

	conn, err := net.Dial("tcp", "localhost:8087")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}

	defer conn.Close()
	fmt.Println("Connected to server")

	// read data sent by server
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("EOF reached. Server closed the connection.")
				break
			}
			fmt.Println("Error reading data from server:", err)
			break
		}
		fmt.Printf("We got : %s\n", string(buf[:n]))
	}

}
