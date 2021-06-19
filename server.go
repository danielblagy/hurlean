package hurlean


import (
	"net"
	"errors"
	"strconv"
	"fmt"
)


// TODO : client handling interface (onConnect, onDicsonnec, onMessage)


func StartServer(port int) error {
	
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
		
		fmt.Println("A new client has been connected to the server")
		
		go handleClient(conn)
	}
	
	return nil
}

func handleClient(conn net.Conn) {
	
	defer conn.Close()
	
	buffer := make([]byte, 1024)
	
	for {
		reqLen, err := conn.Read(buffer)
		
		if err != nil {
			fmt.Println(err)
			
			return
		} else {
			fmt.Println(reqLen)
		}
	}
}