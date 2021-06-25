package hurlean


import (
	"net"
	"errors"
	"strconv"
	"fmt"
	"sync"
	"encoding/gob"
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
	
	stopSignalChannel := make(chan struct{})
	
	var serverUpdateWaitGroup = sync.WaitGroup{}
	serverUpdateWaitGroup.Add(1)
	
	go func(serverInstance *ServerInstance, serverUpdateWaitGroup *sync.WaitGroup) {
		
		for serverInstance.Running {
			serverUpdater.OnServerUpdate(serverInstance)
		}
		
		fmt.Println("ServerUpdate has stopped")
		
		//ln.Close()
		close(stopSignalChannel)
		
		serverUpdateWaitGroup.Done()
	}(&serverInstance, &serverUpdateWaitGroup)
	
	var clientConnectionsWaitGroup = sync.WaitGroup{}
	
	go func() {
		
		loop:
		for {
			select {
			case <- stopSignalChannel:
				fmt.Println("ServerListen has stopped")
				break loop
			
			default:
				conn, err := ln.Accept()
				if err != nil {
					//return errors.New("Failed to accept a client: " + err.Error())
					fmt.Println("Failed to accept a client", err)
				}
				
				newId := serverInstance.IDCounter
				serverInstance.IDCounter += 1
				
				clientHandler.OnClientConnect(newId)
				
				clientConnectionsWaitGroup.Add(1)
				go handleClient(stopSignalChannel, &clientConnectionsWaitGroup, newId, conn, clientHandler)
			}
		}
	}()
	
	clientConnectionsWaitGroup.Wait()
	
	serverUpdateWaitGroup.Wait()
	
	return nil
}

func handleClient(
	stopSignalChannel <-chan struct{},
	clientConnectionsWaitGroup *sync.WaitGroup,
	id uint32, conn net.Conn, clientHandler ClientHandler) {
	
	defer disconnectClient(id, conn, clientHandler)
	
	messageChannel := make(chan Message)
	doneChannel := make(chan struct{})
	
	var wg = sync.WaitGroup{}
	
	wg.Add(2)
	go listenToMessages(stopSignalChannel, messageChannel, doneChannel, &wg, conn)
	go handleMessage(messageChannel, doneChannel, &wg, id, conn, clientHandler)
	
	wg.Wait()
	
	fmt.Println(id, "HandleClient has stopped")
	
	clientConnectionsWaitGroup.Done()
}

// sender
func listenToMessages(
	stopSignalChannel <-chan struct{},
	messageChannel chan<- Message, doneChannel chan<- struct{},
	wg *sync.WaitGroup,
	conn net.Conn) {
	
	loop:
	for {
		select {
		case <- stopSignalChannel:
			break loop
		
		default:
			var message Message
			decoder := gob.NewDecoder(conn)
			if err := decoder.Decode(&message); err != nil {
				fmt.Println("Server Error (message decoding): ", err)
				//doneChannel <- struct{}{}
				close(doneChannel)
				break
			} else {
				messageChannel <- message
			}
		}
	}
	
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
	
	fmt.Println(id, "Disconnected")
}