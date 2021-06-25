package hurlean


import (
	"net"
	"errors"
	"strconv"
	"fmt"
	"sync"
	"encoding/gob"
	"time"
)


type ServerInstance struct {
	IDCounter uint32
	Running bool
}

type Message struct {
	Type string
	Size uint32
	Body string
}

type ClientHandler interface {
	
	OnClientConnect(id uint32)
	OnClientDisconnect(id uint32)
	OnClientMessage(id uint32, message Message) (Message, bool)	// returns (responseMessage, sendResponse)
}

type ServerUpdater interface {
	
	OnServerUpdate(serverInstance *ServerInstance)
}


func StartServer(port int, clientHandler ClientHandler, serverUpdater ServerUpdater) error {
	
	ln, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		return errors.New("Failed to set up server application: " + err.Error())
	}
	defer ln.Close()
	
	serverInstance := ServerInstance{
		IDCounter: 0,
		Running: true,
	}
	
	var serverUpdateWaitGroup = sync.WaitGroup{}
	serverUpdateWaitGroup.Add(1)
	
	go func(serverInstance *ServerInstance, serverUpdateWaitGroup *sync.WaitGroup) {
		
		for serverInstance.Running {
			serverUpdater.OnServerUpdate(serverInstance)
		}
		
		fmt.Println("ServerUpdate has stopped")
		
		ln.Close()
		
		serverUpdateWaitGroup.Done()
	}(&serverInstance, &serverUpdateWaitGroup)
	
	var clientConnectionsWaitGroup = sync.WaitGroup{}
	
	for serverInstance.Running {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Failed to accept a client: ", err)
		} else {
			newId := serverInstance.IDCounter
			serverInstance.IDCounter += 1
			
			clientHandler.OnClientConnect(newId)
			
			clientConnectionsWaitGroup.Add(1)
			go handleClient(&serverInstance, &clientConnectionsWaitGroup, newId, conn, clientHandler)
		}
	}
	
	fmt.Println("ServerListen has stopped")
	
	clientConnectionsWaitGroup.Wait()
	
	serverUpdateWaitGroup.Wait()
	
	return nil
}

func handleClient(
	serverInstance *ServerInstance,
	clientConnectionsWaitGroup *sync.WaitGroup,
	id uint32, conn net.Conn, clientHandler ClientHandler) {
	
	defer disconnectClient(id, conn, clientHandler)
	
	messageChannel := make(chan Message)
	doneChannel := make(chan struct{})
	
	var wg = sync.WaitGroup{}
	
	wg.Add(2)
	go listenToMessages(serverInstance, messageChannel, doneChannel, &wg, conn)
	go handleMessage(messageChannel, doneChannel, &wg, id, conn, clientHandler)
	
	wg.Wait()
	
	fmt.Println("HandleClient has stopped")
	
	clientConnectionsWaitGroup.Done()
}

// sender
func listenToMessages(
	serverInstance *ServerInstance,
	messageChannel chan<- Message, doneChannel chan<- struct{},
	wg *sync.WaitGroup,
	conn net.Conn) {
	
	conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
	
	decoder := gob.NewDecoder(conn)
	
	for serverInstance.Running {
		fmt.Println("Sender update")
		var message Message
		if err := decoder.Decode(&message); err != nil {
			if err.(net.Error).Timeout() {
				continue
			} else {
				fmt.Println("Server Error (message decoding): ", err)
				break
			}
		} else {
			messageChannel <- message
		}
	}
	
	close(doneChannel)
	
	fmt.Println("Sender has stopped")
	
	wg.Done()
}

// receiver
func handleMessage(
	messageChannel <-chan Message, doneChannel <-chan struct{},
	wg *sync.WaitGroup,
	id uint32, conn net.Conn, clientHandler ClientHandler) {
	
	loop:
	for {
		fmt.Println("Sender update")
		select {
			case message := <- messageChannel:
				if responseMessage, sendResponse := clientHandler.OnClientMessage(id, message); sendResponse {
					encoder := gob.NewEncoder(conn)
					if err := encoder.Encode(responseMessage); err != nil {
						fmt.Println("Server Error (message encoding): ", err)
					}
				}
			
			case <- doneChannel:
				break loop
		}
	}
	
	fmt.Println("Receiver has stopped")
	
	wg.Done()
}

func disconnectClient(id uint32, conn net.Conn, clientHandler ClientHandler) {
	
	clientHandler.OnClientDisconnect(id)
	conn.Close()
}