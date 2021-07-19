package hurlean


import (
	"net"
	"errors"
	"fmt"
	"encoding/gob"
	"sync"
	"time"
	"reflect"
	"io"
)


// When the server connects to the server with hurlean.ConnectToServer function call,
// ClientInstance object will be created on success.
// Client Instance is used to control the client's state and send messages to the server
type ClientInstance struct {
	Connected bool
	Conn net.Conn
}

// Sends a message to the server
func (ci ClientInstance) Send(message Message) {
	
	encoder := gob.NewEncoder(ci.Conn)
	if err := encoder.Encode(message); err != nil {
		fmt.Printf(
			"Client Error (message encoding): encoding message = [%v] for sending to the server, error = [%v]\n",
			message, err)
	}
}

// Disconnects from the server
func (ci *ClientInstance) Disconnect() {
	
	ci.Connected = false
	ci.Conn.Close()
}


type ClientFunctionalityProvider interface {
	
	// you can specify client state in your ClientHandler implementation
	
	// Is called when the client receives a message from the server,
	// 'message' is the received message
	OnServerMessage(clientInstance *ClientInstance, message Message)
	
	// Is called once, when the client application starts
	OnClientInit(clientInstance *ClientInstance)
	
	// Is called on each client update, used as a 'main' logic function,
	// e.g. getting an input from the user of the client application
	OnClientUpdate(clientInstance *ClientInstance)
}

// Attempts to connect to the server on ip:port
// returns error on failure
func ConnectToServer(ip string, port string, clientFunctionalityProvider ClientFunctionalityProvider) error {
	
	conn, err := net.Dial("tcp", ip + ":" + port)
	if err != nil {
		return errors.New("__hurlean__  Failed to connect to the server: " + err.Error())
	}
	defer conn.Close()
	
	fmt.Printf("__hurlean__  Successfully connected to the server on %v:%v\n", ip, port)
	
	clientInstance := ClientInstance{
		Connected: true,
		Conn: conn,
	}
	
	clientFunctionalityProvider.OnClientInit(&clientInstance)
	
	var clientUpdateWaitGroup = sync.WaitGroup{}
	clientUpdateWaitGroup.Add(1)
	
	go func(clientInstance *ClientInstance, clientUpdateWaitGroup *sync.WaitGroup) {
		
		for clientInstance.Connected {
			clientFunctionalityProvider.OnClientUpdate(clientInstance)
		}
		
		// DEBUG MESSAGE
		if (debug) { fmt.Println("__hurlean__  ClientUpdate has stopped") }
		
		clientUpdateWaitGroup.Done()
	}(&clientInstance, &clientUpdateWaitGroup)
	
	clientInstance.Conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
	
	// used to check if err in decoder.Decode is of type net.Error, because err may be EOF,
	// which is not of type net.Error, so the program panics, the additional checking prevents that
	netErrorType := reflect.TypeOf((*net.Error)(nil)).Elem()
	
	var message Message
	
	for clientInstance.Connected {
		decoder := gob.NewDecoder(clientInstance.Conn)
		if err := decoder.Decode(&message); err != nil {
			if reflect.TypeOf(err).Implements(netErrorType) && err.(net.Error).Timeout() {
				clientInstance.Conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
				continue
			} else if err == io.EOF {
				fmt.Printf("__hurlean__  Client: connection %v has been closed\n", err)
				break
			} else {
				fmt.Println("__hurlean__  Client Error (message decoding): ", err)
				break
			}
		} else {
			clientFunctionalityProvider.OnServerMessage(&clientInstance, message)
		}
	}
	
	clientInstance.Connected = false
	
	// DEBUG MESSAGE
	if (debug) { fmt.Println("__hurlean__  ClientRead has stopped") }
	
	clientUpdateWaitGroup.Wait()
	
	return nil
}