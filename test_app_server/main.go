// this file is not a part of the library, it's ised to test the execution

package main


import (
	"github.com/danielblagy/hurlean"
	"fmt"
)


type MyClientHandler struct{}

func (ch MyClientHandler) OnClientConnect() {
	
	fmt.Println("A new client has been connected to the server")
}

func (ch MyClientHandler) OnClientDisconnect() {
	
	fmt.Println("Client has disconnected from the server")
}

func (ch MyClientHandler) OnClientMessage(message []byte) {
	
	//fmt.Println(reqLen)
	fmt.Println(string(message))
}


func main() {
	
	var myClientHandler hurlean.ClientHandler = MyClientHandler{}
	
	if err := hurlean.StartServer(8080, myClientHandler); err != nil {
		fmt.Println(err)
	}
}