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
	encoder.Encode(helloMessage)
	//conn.Write([]byte("hello server"))
	
	//buffer := make([]byte, 1024)
	
	for {
		//_, err := conn.Read(buffer)
		var message Message
		decoder := gob.NewDecoder(conn)
		err := decoder.Decode(&message)
		
		if err != nil {
			fmt.Println("Client: ", err)
			
			break
		} else {
			messageHandler.OnServerMessage(message)
		}
	}
	
	return nil
}