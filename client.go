package hurlean


import (
	"net"
	"errors"
	"strconv"
	"fmt"
)


type ServerMessageHandler interface {
	
	OnServerMessage(message []byte)
}


func ConnectToServer(ip string, port int, messageHandler ServerMessageHandler) error {
	
	conn, err := net.Dial("tcp", ip + ":" + strconv.Itoa(port))
	if err != nil {
		return errors.New("Failed to connect to the server: " + err.Error())
	}
	defer conn.Close()
	
	fmt.Println("Successfully connected to the server")
	
	conn.Write([]byte("hello server"))
	
	buffer := make([]byte, 1024)
	
	for {
		_, err := conn.Read(buffer)
		
		if err != nil {
			fmt.Println("Client: ", err)
			
			break
		} else {
			messageHandler.OnServerMessage(buffer)
		}
	}
	
	return nil
}