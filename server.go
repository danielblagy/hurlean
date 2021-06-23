package hurlean


import (
	"net"
	"errors"
	"strconv"
	"fmt"
	"sync"
	"encoding/gob"
)


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


func StartServer(port int, clientHandler ClientHandler) error {
	
	ln, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		return errors.New("Failed to set up server application: " + err.Error())
	}
	defer ln.Close()
	
	var idCounter uint32 = 0
	
	for {
		conn, err := ln.Accept()
		if err != nil {
			return errors.New("Failed to accept a client: " + err.Error())
		}
		
		newId := idCounter
		idCounter += 1
		
		clientHandler.OnClientConnect(newId)
		
		go handleClient(newId, conn, clientHandler)
	}
	
	return nil
}

func handleClient(id uint32, conn net.Conn, clientHandler ClientHandler) {
	
	defer disconnectClient(id, conn, clientHandler)
	
	messageChannel := make(chan Message)
	doneChannel := make(chan struct{})
	
	var wg = sync.WaitGroup{}
	
	wg.Add(2)
	go listenToMessages(messageChannel, doneChannel, &wg, conn)
	go handleMessage(messageChannel, doneChannel, &wg, id, conn, clientHandler)
	
	wg.Wait()
}

// sender
func listenToMessages(messageChannel chan<- Message, doneChannel chan<- struct{}, wg *sync.WaitGroup, conn net.Conn) {
	
	for {
		var message Message
		decoder := gob.NewDecoder(conn)
		err := decoder.Decode(&message)
		
		if err != nil {
			fmt.Println("Server: ", err)
			//doneChannel <- struct{}{}
			close(doneChannel)
			break
		} else {
			messageChannel <- message
		}
	}
	
	wg.Done()
}

// receiver
func handleMessage(messageChannel <-chan Message, doneChannel <-chan struct{}, wg *sync.WaitGroup, id uint32, conn net.Conn, clientHandler ClientHandler) {
	
	loop:
	for {
		select {
			case message := <- messageChannel:
				if responseMessage, sendResponse := clientHandler.OnClientMessage(id, message); sendResponse {
					// TODO : check for errors in Write
					encoder := gob.NewEncoder(conn)
					encoder.Encode(responseMessage)
				}
			
			case <- doneChannel:
				break loop
		}
	}
	
	wg.Done()
}

func disconnectClient(id uint32, conn net.Conn, clientHandler ClientHandler) {
	
	clientHandler.OnClientDisconnect(id)
	conn.Close()
}