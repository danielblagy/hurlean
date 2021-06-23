// this file is not a part of the library, it's ised to test the execution

package main


import (
	"github.com/danielblagy/hurlean"
	"fmt"
)


type MyClientHandler struct{}

func (ch MyClientHandler) OnClientConnect(id uint32) {
	
	fmt.Println("A new client (id", id, ") has been connected to the server")
}

func (ch MyClientHandler) OnClientDisconnect(id uint32) {
	
	fmt.Println("Client (id", id, ") has disconnected from the server")
}

func (ch MyClientHandler) OnClientMessage(id uint32, message []byte) ([]byte, bool) {
	
	convertedMessage := string(message)
	
	fmt.Println(id, ":", convertedMessage)
	
	//responseMessage := "echo from server: " + convertedMessage
	responseMessage := "echo from server"
	
	return []byte(responseMessage), true
}


func main() {
	
	var myClientHandler hurlean.ClientHandler = MyClientHandler{}
	
	if err := hurlean.StartServer(8080, myClientHandler); err != nil {
		fmt.Println(err)
	}
}