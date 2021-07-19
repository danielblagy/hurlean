package hurlean


import (
	"net"
	"errors"
	"fmt"
	"sync"
	"encoding/gob"
	"time"
	"reflect"
	"io"
)


// When the server is started with hurlean.StartServer function call,
// ServerInstance object will be created on success.
// Server Instance is used to control the server's state and send messages to clients
type ServerInstance struct {
	IDCounter uint32
	Running bool
	Clients map[uint32]net.Conn
	clientsMutex sync.RWMutex
	State interface{}
}

// Sends a message to the client with id
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

// Sends a message to all the clients
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

// Stops the server
func (si *ServerInstance) Stop() {
	
	si.Running = false
}

// Disconnects the client from the server
func (si *ServerInstance) Disconnect(id uint32) {
	
	if conn, ok := si.Clients[id]; ok {
		conn.Close()
		delete(si.Clients, id)
	}
}

// Disconnects all clients from the server
func (si *ServerInstance) DisconnectAll() {
	
	for id, conn := range si.Clients {
		conn.Close()
		delete(si.Clients, id)
	}
}


type ServerFunctionalityProvider interface {
	
	// you can specify client state in your ServerFunctionalityProvider implementation
	
	// Is called when a new client connect to the server,
	// 'id' is the new clients's id
	OnClientConnect(si *ServerInstance, id uint32)
	
	// Is called when a client disconnects from the server,
	// 'id' is the disconnected clients's id
	OnClientDisconnect(si *ServerInstance, id uint32)
	
	// Is called when the server receives a message from a client,
	// 'id' is the clients's id
	// 'message' is the received message
	OnClientMessage(si *ServerInstance, id uint32, message Message)
	
	// Is called once, when the server application starts,
	OnServerInit(serverInstance *ServerInstance)
	
	// Is called on each server update, used as a 'main' logic function,
	// e.g. getting an input from the user of the server application
	OnServerUpdate(serverInstance *ServerInstance)
}


// Starts a server on port
// returns error on failure
// serverState parameter can be of any type and will be accessible via *hurlean.ServerInstance
func StartServer(port string, serverFunctionalityProvider ServerFunctionalityProvider) error {
	
	ln, err := net.Listen("tcp", ":" + port)
	if err != nil {
		return errors.New("__hurlean__  Failed to set up server application: " + err.Error())
	}
	defer ln.Close()
	
	serverInstance := ServerInstance{
		IDCounter: 0,
		Running: true,
		Clients: make(map[uint32]net.Conn),
		clientsMutex: sync.RWMutex{},
	}
	
	serverFunctionalityProvider.OnServerInit(&serverInstance)
	
	var serverUpdateWaitGroup = sync.WaitGroup{}
	serverUpdateWaitGroup.Add(1)
	
	go func(serverInstance *ServerInstance, serverUpdateWaitGroup *sync.WaitGroup) {
		
		for serverInstance.Running {
			serverFunctionalityProvider.OnServerUpdate(serverInstance)
		}
		
		// DEBUG MESSAGE
		if (debug) { fmt.Println("__hurlean__  ServerUpdate has stopped") }
		
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
			go handleClient(&serverInstance, &clientConnectionsWaitGroup, newId, conn, serverFunctionalityProvider)
			
			serverFunctionalityProvider.OnClientConnect(&serverInstance, newId)
		}
	}
	
	// DEBUG MESSAGE
	if (debug) { fmt.Println("__hurlean__  ServerListen has stopped") }
	
	clientConnectionsWaitGroup.Wait()
	
	serverUpdateWaitGroup.Wait()
	
	return nil
}

func handleClient(
	serverInstance *ServerInstance,
	clientConnectionsWaitGroup *sync.WaitGroup,
	id uint32, conn net.Conn, serverFunctionalityProvider ServerFunctionalityProvider) {
	
	serverInstance.Clients[id] = conn
	
	messageChannel := make(chan Message)
	doneChannel := make(chan struct{})
	
	var wg = sync.WaitGroup{}
	
	wg.Add(2)
	go listenToMessages(serverInstance, messageChannel, doneChannel, &wg, conn)
	go handleMessage(serverInstance, messageChannel, doneChannel, &wg, id, conn, serverFunctionalityProvider)
	
	wg.Wait()
	
	disconnectClient(serverInstance, id, conn, serverFunctionalityProvider)
	
	// DEBUG MESSAGE
	if (debug) { fmt.Println("__hurlean__  HandleClient has stopped") }
	
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
	
	var message Message
	
	for serverInstance.Running {
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
	if (debug) { fmt.Println("__hurlean__  Sender has stopped") }
	
	wg.Done()
}

// receiver
func handleMessage(
	serverInstance *ServerInstance,
	messageChannel <-chan Message, doneChannel <-chan struct{},
	wg *sync.WaitGroup,
	id uint32, conn net.Conn, serverFunctionalityProvider ServerFunctionalityProvider) {
	
	loop:
	for {
		select {
			case message := <- messageChannel:
				serverFunctionalityProvider.OnClientMessage(serverInstance, id, message)
			
			case <- doneChannel:
				break loop
		}
	}
	
	// DEBUG MESSAGE
	if (debug) { fmt.Println("__hurlean__  Receiver has stopped") }
	
	wg.Done()
}

func disconnectClient(serverInstance *ServerInstance, id uint32, conn net.Conn, serverFunctionalityProvider ServerFunctionalityProvider) {
	
	serverFunctionalityProvider.OnClientDisconnect(serverInstance, id)
	conn.Close()
}