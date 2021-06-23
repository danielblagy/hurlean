package hurlean


import (
	"net"
	"errors"
	"strconv"
	"fmt"
	"encoding/gob"
)


type ServerMessageHandler interface {
	
	OnServerMessage(message Message)
}


func ConnectToServer(ip string, port int, messageHandler ServerMessageHandler) error {
	
	conn, err := net.Dial("tcp", ip + ":" + strconv.Itoa(port))
	if err != nil {
		return errors.New("Failed to connect to the server: " + err.Error())
	}
	defer conn.Close()
	
	fmt.Println("Successfully connected to the server")
	
	// TODO : check for errors in Write
	helloMessage := Message{
		Type: "hello",
		Size: 0,
		Body: "hello server",
	}
	encoder := gob.NewEncoder(conn)
	if err := encoder.Encode(helloMessage); err != nil {
		fmt.Println("Client Error (message encoding): ", err)
	}
	
	for {
		var message Message
		decoder := gob.NewDecoder(conn)
		if err := decoder.Decode(&message); err != nil {
			fmt.Println("Client Error (message decoding): ", err)
			
			break
		} else {
			messageHandler.OnServerMessage(message)
		}
	}
	
	return nil
}