package hurlean


import (
	"net"
	"errors"
	"strconv"
	"fmt"
	"encoding/gob"
	"sync"
	"time"
)


type ClientInstance struct {
	Connected bool
	Conn net.Conn
}

func (ci ClientInstance) Send(message Message) {
	
	encoder := gob.NewEncoder(ci.Conn)
	if err := encoder.Encode(message); err != nil {
		fmt.Printf(
			"Client Error (message encoding): encoding message = [%v] for sending to the server, error = [%v]\n",
			message, err)
	}
}

func (ci *ClientInstance) Disconnect() {
	
	ci.Connected = false
	ci.Conn.Close()
}

type ServerMessageHandler interface {
	
	OnServerMessage(message Message)
}

type ClientUpdater interface {
	
	OnClientUpdate(clientInstance *ClientInstance)
}


func ConnectToServer(ip string, port int, messageHandler ServerMessageHandler, clientUpdater ClientUpdater) error {
	
	conn, err := net.Dial("tcp", ip + ":" + strconv.Itoa(port))
	if err != nil {
		return errors.New("Failed to connect to the server: " + err.Error())
	}
	defer conn.Close()
	
	fmt.Println("Successfully connected to the server")
	
	clientInstance := ClientInstance{
		Connected: true,
		Conn: conn,
	}
	
	var clientUpdateWaitGroup = sync.WaitGroup{}
	clientUpdateWaitGroup.Add(1)
	
	go func(clientInstance *ClientInstance, clientUpdateWaitGroup *sync.WaitGroup) {
		
		for clientInstance.Connected {
			clientUpdater.OnClientUpdate(clientInstance)
		}
		
		// DEBUG MESSAGE
		fmt.Println("ClientUpdate has stopped")
		
		clientUpdateWaitGroup.Done()
	}(&clientInstance, &clientUpdateWaitGroup)
	
	helloMessage := Message{
		Type: "hello",
		Body: "hello server",
	}
	clientInstance.Send(helloMessage)
	
	clientInstance.Conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
	
	decoder := gob.NewDecoder(clientInstance.Conn)
	
	for clientInstance.Connected {
		// TODO : move message var outside for
		var message Message
		if err := decoder.Decode(&message); err != nil {
			if err.(net.Error).Timeout() {
				continue
			} else {
				fmt.Println("Client Error (message decoding): ", err)
				break
			}
		} else {
			messageHandler.OnServerMessage(message)
		}
	}
	
	// DEBUG MESSAGE
	fmt.Println("ClientRead has stopped")
	
	clientUpdateWaitGroup.Wait()
	
	return nil
}