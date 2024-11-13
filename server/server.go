package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", ":8087")

	if err != nil {
		fmt.Println("Error in setting up a listener: ", err)
	}

	defer listener.Close()
	fmt.Println("Server is listening on port 8087.....")

	// open the Pipe
	pipe, err := os.Open("/tmp/songfifo")
	if err != nil {
		fmt.Println("Error opening pipe:", err)
		return
	}
	defer pipe.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error in accepting the connection:", err)
			continue
		}
		fmt.Println("CLient is connected")

		// Continuously read data from the pipe and send it to the client
		buf := make([]byte, 32768)

		// Goroutine to handle reading from the pipe and sending data to the client
		go func(conn net.Conn) {
			defer conn.Close()

			for {
				n, err := pipe.Read(buf)
				if err != nil {

					if err.Error() == "EOF" {
						fmt.Println("EOF reached")
						break
					}
					fmt.Println("Error reading from pipe:", err)
					break
				}

				// sendinf data to client
				_, err = conn.Write(buf[:n])
				fmt.Printf("Read %d bytes: %s\n", n, string(buf[:5]))
				if err != nil {
					fmt.Println("Error in sending data to client:", err)
					break
				}
			}
			conn.Close()
			fmt.Println("Client is gone !")
		}(conn)
	}

}
