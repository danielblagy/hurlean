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
	Clients map[uint32]net.Conn
}

func (si ServerInstance) Send(id uint32, message Message) {
	
	if conn, ok := si.Clients[id]; ok {
		encoder := gob.NewEncoder(conn)
		if err := encoder.Encode(message); err != nil {
			fmt.Printf(
				"Server Error (message encoding ): encoding message = [%v] for sending to client with id = [%v], error = p%v[",
				message, id, err)
		}
	}
}

type Message struct {
	Type string
	Size uint32
	Body string
}

type ClientHandler interface {
	
	OnClientConnect(si *ServerInstance, id uint32)
	OnClientDisconnect(si *ServerInstance, id uint32)
	OnClientMessage(si *ServerInstance, id uint32, message Message)
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
		Clients: make(map[uint32]net.Conn),
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
			
			clientConnectionsWaitGroup.Add(1)
			go handleClient(&serverInstance, &clientConnectionsWaitGroup, newId, conn, clientHandler)
			
			clientHandler.OnClientConnect(&serverInstance, newId)
		}
	}
	
	fmt.Println("ServerListen has stopped")
	
	clientConnectionsWaitGroup.Wait()
	
	serverUpdateWaitGroup.Wait()
	
	return nil
}

func Stop(serverInstance *ServerInstance) {
	
	serverInstance.Running = false
}

func handleClient(
	serverInstance *ServerInstance,
	clientConnectionsWaitGroup *sync.WaitGroup,
	id uint32, conn net.Conn, clientHandler ClientHandler) {
	
	defer disconnectClient(serverInstance, id, conn, clientHandler)
	
	serverInstance.Clients[id] = conn
	
	messageChannel := make(chan Message)
	doneChannel := make(chan struct{})
	
	var wg = sync.WaitGroup{}
	
	wg.Add(2)
	go listenToMessages(serverInstance, messageChannel, doneChannel, &wg, conn)
	go handleMessage(serverInstance, messageChannel, doneChannel, &wg, id, conn, clientHandler)
	
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
	
	// set a read deadline to force the loop to update (for now it's each 100 ms),
	// because otherwise the Decode function blocks (because conn.Read is a blocking operation)
	// and if the server has been stopped, this function won't know about it and the go routine won't be
	// closed properly
	conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
	
	decoder := gob.NewDecoder(conn)
	
	for serverInstance.Running {
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
	serverInstance *ServerInstance,
	messageChannel <-chan Message, doneChannel <-chan struct{},
	wg *sync.WaitGroup,
	id uint32, conn net.Conn, clientHandler ClientHandler) {
	
	loop:
	for {
		select {
			case message := <- messageChannel:
				clientHandler.OnClientMessage(serverInstance, id, message)
			
			case <- doneChannel:
				break loop
		}
	}
	
	fmt.Println("Receiver has stopped")
	
	wg.Done()
}

func disconnectClient(serverInstance *ServerInstance, id uint32, conn net.Conn, clientHandler ClientHandler) {
	
	clientHandler.OnClientDisconnect(serverInstance, id)
	conn.Close()
}