package hurlean


import (
	"net"
	"strconv"
	"fmt"
)


func StartServer(port int) {
	
	socket, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		fmt.Println("Failed to set up server application", err)
	}
	
	for {
		conn, err := socket.Accept()
		if err != nil {
			fmt.Println("Failed to connect a client", err)
		} else {
			fmt.Println("A new client has been connected to the server")
			fmt.Println(conn)
		}
	}
}