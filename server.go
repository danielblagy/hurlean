package hurlean


import (
	"net"
	"errors"
	"strconv"
	"fmt"
	"sync"
)


type ClientHandler interface {
	
	OnClientConnect()
	OnClientDisconnect()
	OnClientMessage(message []byte) ([]byte, bool)	// returns (responseMessage, sendResponse)
}


func StartServer(port int, clientHandler ClientHandler) error {
	
	ln, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		return errors.New("Failed to set up server application: " + err.Error())
	}
	defer ln.Close()
	
	for {
		conn, err := ln.Accept()
		if err != nil {
			return errors.New("Failed to accept a client: " + err.Error())
		}
		
		clientHandler.OnClientConnect()
		
		go handleClient(conn, clientHandler)
	}
	
	return nil
}

func handleClient(conn net.Conn, clientHandler ClientHandler) {
	
	defer disconnectClient(conn, clientHandler)
	
	messageChannel := make(chan []byte)
	doneChannel := make(chan struct{})
	
	var wg = sync.WaitGroup{}
	
	wg.Add(2)
	go listenToMessages(messageChannel, &wg, conn)
	go handleMessage(messageChannel, doneChannel, &wg, conn, clientHandler)
	
	wg.Wait()
}

// sender
func listenToMessages(messageChannel chan <- []byte, wg *sync.WaitGroup, conn net.Conn) {
	
	buffer := make([]byte, 1024)
	
	for {
		_, err := conn.Read(buffer)
		
		if err != nil {
			fmt.Println("Server: ", err)
			
			return
		} else {
			messageChannel <- buffer
		}
	}
	
	wg.Done()
}

// receiver
func handleMessage(messageChannel <- chan []byte, doneChannel <- chan struct{}, wg *sync.WaitGroup, conn net.Conn, clientHandler ClientHandler) {
	
	for {
		select {
			case buffer := <- messageChannel:
				if responseMessage, sendResponse := clientHandler.OnClientMessage(buffer); sendResponse {
					// TODO : check for errors
					conn.Write(responseMessage)
				}
			
			case <- doneChannel:
				break
		}
	}
	
	wg.Done()
}

func disconnectClient(conn net.Conn, clientHandler ClientHandler) {
	
	clientHandler.OnClientDisconnect()
	conn.Close()
}