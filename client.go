package hurlean


import (
	"net"
	"errors"
	"strconv"
	"fmt"
)


func ConnectToServer(ip string, port int) error {
	
	conn, err := net.Dial("tcp", ip + ":" + strconv.Itoa(port))
	if err != nil {
		return errors.New("Failed to connect to the server: " + err.Error())
	}
	defer conn.Close()
	
	fmt.Println("Successfully connected to the server")
	
	conn.Write([]byte("hello server"))
	
	return nil
}