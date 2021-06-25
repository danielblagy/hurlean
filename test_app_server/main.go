// this file is not a part of the library, it's ised to test the execution

package main


import (
	"github.com/danielblagy/hurlean"
	"fmt"
)


type MyClientHandler struct{}

func (ch MyClientHandler) OnClientConnect(si *hurlean.ServerInstance, id uint32) {
	
	fmt.Println("A new client (id", id, ") has been connected to the server", si)
}

func (ch MyClientHandler) OnClientDisconnect(si *hurlean.ServerInstance, id uint32) {
	
	fmt.Println("Client (id", id, ") has disconnected from the server", si)
}

func (ch MyClientHandler) OnClientMessage(si *hurlean.ServerInstance, id uint32, message hurlean.Message) {
	
	fmt.Println("")
	fmt.Println("----------------")
	fmt.Println("Message from", id, ":")
	fmt.Println("  Type:", message.Type)
	fmt.Println("  Body:", message.Body)
	fmt.Println("----------------")
	
	responseMessage := hurlean.Message{
		Type: "echo",
		Body: "echo from server",
	}
	
	si.Send(id, responseMessage)
}


type MyServerUpdater struct{}

func (su MyServerUpdater) OnServerUpdate(serverInstance *hurlean.ServerInstance) {
	
	var input string
	fmt.Scanln(&input)
	switch (input) {
	case "exit":
		hurlean.Stop(serverInstance)
	}
}


func main() {
	
	var myClientHandler hurlean.ClientHandler = MyClientHandler{}
	var MyServerUpdater hurlean.ServerUpdater = MyServerUpdater{}
	
	if err := hurlean.StartServer(8080, myClientHandler, MyServerUpdater); err != nil {
		fmt.Println(err)
	}
}