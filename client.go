package hurlean


import (
	"net"
	"errors"
	"strconv"
	"fmt"
	"encoding/gob"
	"sync"
	"time"
	"reflect"
	"io"
)


type ClientInstance struct {
	Connected bool
	Conn net.Conn
	State interface{}
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


func ConnectToServer(ip string, port int, messageHandler ServerMessageHandler, clientUpdater ClientUpdater, clientState interface{}) error {
	
	conn, err := net.Dial("tcp", ip + ":" + strconv.Itoa(port))
	if err != nil {
		return errors.New("__hurlean__  Failed to connect to the server: " + err.Error())
	}
	defer conn.Close()
	
	fmt.Printf("__hurlean__  Successfully connected to the server on %v:%v\n", ip, port)
	
	clientInstance := ClientInstance{
		Connected: true,
		Conn: conn,
		State: clientState,
	}
	
	var clientUpdateWaitGroup = sync.WaitGroup{}
	clientUpdateWaitGroup.Add(1)
	
	go func(clientInstance *ClientInstance, clientUpdateWaitGroup *sync.WaitGroup) {
		
		for clientInstance.Connected {
			clientUpdater.OnClientUpdate(clientInstance)
		}
		
		// DEBUG MESSAGE
		if (debug) {
			fmt.Println("__hurlean__  ClientUpdate has stopped")
		}
		
		clientUpdateWaitGroup.Done()
	}(&clientInstance, &clientUpdateWaitGroup)
	
	/*helloMessage := Message{
		Type: "hello",
		Body: "hello server",
	}
	clientInstance.Send(helloMessage)*/
	
	clientInstance.Conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
	
	// used to chekc if err in decoder.Decode is of type net.Error, because err may be EOF,
	// which is not of type net.Error, so the program panics, the additional checking prevents that
	netErrorType := reflect.TypeOf((*net.Error)(nil)).Elem()
	
	for clientInstance.Connected {
		// TODO : move message var outside for
		var message Message
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
			messageHandler.OnServerMessage(message)
		}
	}
	
	clientInstance.Connected = false
	
	// DEBUG MESSAGE
	if (debug) {
		fmt.Println("__hurlean__  ClientRead has stopped")
	}
	
	clientUpdateWaitGroup.Wait()
	
	return nil
}