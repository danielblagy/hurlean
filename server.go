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
	OnClientMessage(message []byte)
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
	
	defer conn.Close()
	
	buffer := make([]byte, 1024)
	
	for {
		_, err := conn.Read(buffer)
		
		if err != nil {
			fmt.Println("Server: ", err)
			
			return
		} else {
			clientHandler.OnClientMessage(buffer)
		}
	}
}