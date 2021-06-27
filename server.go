package hurlean


import (
	"net"
	"errors"
	"strconv"
	"fmt"
	"sync"
	"encoding/gob"
	"time"
	"reflect"
	"io"
)


type ServerInstance struct {
	IDCounter uint32
	Running bool
	Clients map[uint32]net.Conn
	clientsMutex sync.RWMutex
	State interface{}
}

func (si ServerInstance) Send(id uint32, message Message) {
	
	si.clientsMutex.Lock()
	
	if conn, ok := si.Clients[id]; ok {
		encoder := gob.NewEncoder(conn)
		if err := encoder.Encode(message); err != nil {
			fmt.Printf(
				"Server Error (message encoding ): encoding message = [%v] for sending to client with id = [%v], error = [%v]\n",
				message, id, err)
		}
	}
	
	si.clientsMutex.Unlock()
}

func (si ServerInstance) SendAll(message Message) {
	
	si.clientsMutex.Lock()
	
	for id, conn := range si.Clients {
		encoder := gob.NewEncoder(conn)
		if err := encoder.Encode(message); err != nil {
			fmt.Printf(
				"Server Error (message encoding ): encoding message = [%v] for sending to client with id = [%v], error = [%v]\n",
				message, id, err)
		}
	}
	
	si.clientsMutex.Unlock()
}

func (si *ServerInstance) Stop() {
	
	si.Running = false
}

type Message struct {
	Type string
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


// controls the debug prints
var debug bool = false

func StartServer(port int, clientHandler ClientHandler, serverUpdater ServerUpdater, serverState interface{}) error {
	
	ln, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		return errors.New("__hurlean__  Failed to set up server application: " + err.Error())
	}
	defer ln.Close()
	
	serverInstance := ServerInstance{
		IDCounter: 0,
		Running: true,
		Clients: make(map[uint32]net.Conn),
		clientsMutex: sync.RWMutex{},
		State: serverState,
	}
	
	var serverUpdateWaitGroup = sync.WaitGroup{}
	serverUpdateWaitGroup.Add(1)
	
	go func(serverInstance *ServerInstance, serverUpdateWaitGroup *sync.WaitGroup) {
		
		for serverInstance.Running {
			serverUpdater.OnServerUpdate(serverInstance)
		}
		
		// DEBUG MESSAGE
		if (debug) {
			fmt.Println("__hurlean__  ServerUpdate has stopped")
		}
		
		ln.Close()
		
		serverUpdateWaitGroup.Done()
	}(&serverInstance, &serverUpdateWaitGroup)
	
	var clientConnectionsWaitGroup = sync.WaitGroup{}
	
	fmt.Printf("__hurlean__  The Server has been started on port %v\n", port)
	
	for serverInstance.Running {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("__hurlean__  Failed to accept a client: ", err)
		} else {
			newId := serverInstance.IDCounter
			serverInstance.IDCounter += 1
			
			clientConnectionsWaitGroup.Add(1)
			go handleClient(&serverInstance, &clientConnectionsWaitGroup, newId, conn, clientHandler)
			
			clientHandler.OnClientConnect(&serverInstance, newId)
		}
	}
	
	// DEBUG MESSAGE
	if (debug) {
		fmt.Println("__hurlean__  ServerListen has stopped")
	}
	
	clientConnectionsWaitGroup.Wait()
	
	serverUpdateWaitGroup.Wait()
	
	return nil
}

func handleClient(
	serverInstance *ServerInstance,
	clientConnectionsWaitGroup *sync.WaitGroup,
	id uint32, conn net.Conn, clientHandler ClientHandler) {
	
	serverInstance.Clients[id] = conn
	
	messageChannel := make(chan Message)
	doneChannel := make(chan struct{})
	
	var wg = sync.WaitGroup{}
	
	wg.Add(2)
	go listenToMessages(serverInstance, messageChannel, doneChannel, &wg, conn)
	go handleMessage(serverInstance, messageChannel, doneChannel, &wg, id, conn, clientHandler)
	
	wg.Wait()
	
	disconnectClient(serverInstance, id, conn, clientHandler)
	
	// DEBUG MESSAGE
	if (debug) {
		fmt.Println("__hurlean__  HandleClient has stopped")
	}
	
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
	
	// used to chekc if err in decoder.Decode is of type net.Error, because err may be EOF,
	// which is not of type net.Error, so the program panics, the additional checking prevents that
	netErrorType := reflect.TypeOf((*net.Error)(nil)).Elem()
	
	for serverInstance.Running {
		// TODO : move message var outside for
		var message Message
		decoder := gob.NewDecoder(conn)
		if err := decoder.Decode(&message); err != nil {
			if reflect.TypeOf(err).Implements(netErrorType) && err.(net.Error).Timeout() {
				conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100))
				continue
			} else if err == io.EOF {
				fmt.Printf("__hurlean__  Server: connection %v has been closed\n", conn)
				
				serverInstance.clientsMutex.Lock()
				// remove the connection from the map
				for k, v := range serverInstance.Clients {
					if v == conn {
						delete(serverInstance.Clients, k)
						break
					}
				}
				serverInstance.clientsMutex.Unlock()
				
				break
			} else {
				fmt.Println("__hurlean__  Server Error (message decoding): ", err)
				break
			}
		} else {
			messageChannel <- message
		}
	}
	
	close(doneChannel)
	
	// DEBUG MESSAGE
	if (debug) {
		fmt.Println("__hurlean__  Sender has stopped")
	}
	
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
	
	// DEBUG MESSAGE
	if (debug) {
		fmt.Println("__hurlean__  Receiver has stopped")
	}
	
	wg.Done()
}

func disconnectClient(serverInstance *ServerInstance, id uint32, conn net.Conn, clientHandler ClientHandler) {
	
	clientHandler.OnClientDisconnect(serverInstance, id)
	conn.Close()
}