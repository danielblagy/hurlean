// this file is not a part of the library, it's ised to test the execution

package main


import (
	"github.com/danielblagy/hurlean"
	"fmt"
	"strconv"
	//"bufio"
)


type MyClientHandler struct{}

func (ch MyClientHandler) OnClientConnect(si *hurlean.ServerInstance, id uint32) {
	
	fmt.Println("A new client (id", id, ") has been connected to the server", si)
	fmt.Println(si.State.(MyServerState).i)
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
	
	if message.Type == "chat message" {
		responseMessage := hurlean.Message{
			Type: "chat message",
			Body: strconv.FormatUint(uint64(id), 10) + ": " + message.Body,
		}
		
		si.SendAll(responseMessage)
	}
}


type MyServerUpdater struct{}

func (su MyServerUpdater) OnServerUpdate(serverInstance *hurlean.ServerInstance) {
	
	var input string
	fmt.Scanln(&input)
	switch (input) {
	case "exit":
		serverInstance.Stop()
	}
}


type MyServerState struct{
	i int
}


func main() {
	
	var myClientHandler hurlean.ClientHandler = MyClientHandler{}
	var myServerUpdater hurlean.ServerUpdater = MyServerUpdater{}
	
	if err := hurlean.StartServer(8080, myClientHandler, myServerUpdater, MyServerState{15}); err != nil {
		fmt.Println(err)
	}
}