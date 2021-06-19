package hurlean


import (
	"net"
	"errors"
	"strconv"
	"fmt"
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
	
	buffer := make([]byte, 1024)
	
	for {
		_, err := conn.Read(buffer)
		
		if err != nil {
			fmt.Println("Server: ", err)
			
			return
		} else {
			if responseMessage, sendResponse := clientHandler.OnClientMessage(buffer); sendResponse {
				// TODO : check for errors
				conn.Write(responseMessage)
			}
		}
	}
}

func disconnectClient(conn net.Conn, clientHandler ClientHandler) {
	
	clientHandler.OnClientDisconnect()
	conn.Close()
}