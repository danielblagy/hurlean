// this file is not a part of the library, it's ised to test the execution

package main


import (
	"github.com/danielblagy/hurlean"
	"fmt"
)


type MyServerMessageHandler struct{}

func (mh MyServerMessageHandler) OnServerMessage(message hurlean.Message) {
	
	fmt.Println("")
	fmt.Println("----------------")
	fmt.Println("Message from the server")
	fmt.Println("  Type:", message.Type)
	fmt.Println("  Body:", message.Body)
	fmt.Println("----------------")
}


type MyClientUpdater struct{}

func (cu MyClientUpdater) OnClientUpdate(clientInstance *hurlean.ClientInstance) {
	
	
}


func main() {
	
	var mh hurlean.ServerMessageHandler = MyServerMessageHandler{}
	
	if err := hurlean.ConnectToServer("localhost", 8080, mh); err != nil {
		fmt.Println(err)
	}
}