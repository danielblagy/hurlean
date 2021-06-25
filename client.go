package hurlean


import (
	"net"
	"errors"
	"strconv"
	"fmt"
	"encoding/gob"
)


type ClientInstance struct {
	Conn net.Conn
}

func (ci ClientInstance) Send(message Message) {
	
	encoder := gob.NewEncoder(ci.Conn)
	if err := encoder.Encode(message); err != nil {
		fmt.Printf(
			"Client Error (message encoding): encoding message = [%v] for sending to the server, error = [%v]",
			message, err)
	}
}

type ServerMessageHandler interface {
	
	OnServerMessage(message Message)
}

type ClientUpdater interface {
	
	OnClientUpdate(clientInstance *ClientInstance)
}


func ConnectToServer(ip string, port int, messageHandler ServerMessageHandler) error {
	
	conn, err := net.Dial("tcp", ip + ":" + strconv.Itoa(port))
	if err != nil {
		return errors.New("Failed to connect to the server: " + err.Error())
	}
	defer conn.Close()
	
	fmt.Println("Successfully connected to the server")
	
	clientInstance := ClientInstance{
		Conn: conn,
	}
	
	helloMessage := Message{
		Type: "hello",
		Body: "hello server",
	}
	clientInstance.Send(helloMessage)
	
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