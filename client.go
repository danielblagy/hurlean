package hurlean


import (
	"net"
	"strconv"
	"fmt"
)


func ConnectToServer(ip string, port int) {
	conn, err := net.Dial("tcp", ip + ":" + strconv.Itoa(port))
	if err != nil {
		fmt.Println("Failed to connect to the server", err)
	}
	
	fmt.Println("Successfully connected to the server", conn)
}