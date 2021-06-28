// this file is not a part of the library, it's ised to test the execution

package main


import (
	"github.com/danielblagy/hurlean"
	"fmt"
	"strconv"
	"bufio"
	"os"
	"sync"
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
	
	if message.Type == "chat message" {
		var nickname string
		if name, ok := si.State.(MyServerState).clientNames[id]; ok {
			nickname = name
		} else {
			nickname = strconv.FormatUint(uint64(id), 10)
		}
		
		responseMessage := hurlean.Message{
			Type: "chat message",
			Body: nickname + ": " + message.Body,
		}
		
		si.SendAll(responseMessage)
	} else if message.Type == "setname" {
		// a pointer method workaround
		state := si.State.(MyServerState)
		
		state.clientNamesMutex.Lock()
		state.clientNames[id] = message.Body
		state.clientNamesMutex.Unlock()
		
		fmt.Printf("client %v has set name to %v", id, message.Body)
	}
}


type MyServerUpdater struct{}

func (su MyServerUpdater) OnServerUpdate(serverInstance *hurlean.ServerInstance) {
	
	scanner := serverInstance.State.(MyServerState).scanner
	
	if scanner.Scan() {
		switch (scanner.Text()) {
		case "exit":
			serverInstance.Stop()
			
		case "disconnect":
			serverInstance.DisconnectAll()
		}
	}
}


type MyServerState struct{
	scanner *bufio.Scanner
	clientNames map[uint32]string
	clientNamesMutex sync.RWMutex
}


func main() {
	
	var myClientHandler hurlean.ClientHandler = MyClientHandler{}
	var myServerUpdater hurlean.ServerUpdater = MyServerUpdater{}
	
	// set the app-specific server's state
	var myServerState MyServerState = MyServerState{
		scanner: bufio.NewScanner(os.Stdin),
		clientNames: make(map[uint32]string),
		clientNamesMutex: sync.RWMutex{},
	}
	
	if err := hurlean.StartServer(8080, myClientHandler, myServerUpdater, myServerState); err != nil {
		fmt.Println(err)
	}
}