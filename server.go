package hurlean


import (
	"net"
	"errors"
	"strconv"
	"fmt"
	"sync"
)


type ClientHandler interface {
	
	OnClientConnect(id uint32)
	OnClientDisconnect(id uint32)
	OnClientMessage(id uint32, message []byte) ([]byte, bool)	// returns (responseMessage, sendResponse)
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
	
	messageChannel := make(chan []byte)
	doneChannel := make(chan struct{})
	
	var wg = sync.WaitGroup{}
	
	wg.Add(2)
	go listenToMessages(messageChannel, doneChannel, &wg, conn)
	go handleMessage(messageChannel, doneChannel, &wg, id, conn, clientHandler)
	
	wg.Wait()
}

// sender
func listenToMessages(messageChannel chan<- []byte, doneChannel chan<- struct{}, wg *sync.WaitGroup, conn net.Conn) {
	
	buffer := make([]byte, 1024)
	
	for {
		_, err := conn.Read(buffer)
		
		if err != nil {
			fmt.Println("Server: ", err)
			//doneChannel <- struct{}{}
			close(doneChannel)
			break
		} else {
			messageChannel <- buffer
		}
	}
	
	wg.Done()
}

// receiver
func handleMessage(messageChannel <-chan []byte, doneChannel <-chan struct{}, wg *sync.WaitGroup, id uint32, conn net.Conn, clientHandler ClientHandler) {
	
	loop:
	for {
		select {
			case buffer := <- messageChannel:
				if responseMessage, sendResponse := clientHandler.OnClientMessage(id, buffer); sendResponse {
					// TODO : check for errors
					conn.Write(responseMessage)
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