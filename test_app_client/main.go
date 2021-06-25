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
	
	var input string
	fmt.Scanln(&input)
	switch (input) {
	case "/disconnect":
		clientInstance.Disconnect()
	default:
		message := hurlean.Message{
			Type: "chat message",
			Body: input,
		}
		clientInstance.Send(message)
	}
}


func main() {
	
	var myServerMessageHandler hurlean.ServerMessageHandler = MyServerMessageHandler{}
	var myClientUpdater hurlean.ClientUpdater = MyClientUpdater{}
	
	if err := hurlean.ConnectToServer("localhost", 8080, myServerMessageHandler, myClientUpdater); err != nil {
		fmt.Println(err)
	}
}