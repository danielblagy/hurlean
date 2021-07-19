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


type MyServerFunctionalityProvider struct{
	scanner *bufio.Scanner
	clientNames map[uint32]string
	clientNamesMutex sync.RWMutex
}

func (fp MyServerFunctionalityProvider) OnClientConnect(si *hurlean.ServerInstance, id uint32) {
	
	fmt.Println("A new client (id", id, ") has been connected to the server", si)
}

func (fp MyServerFunctionalityProvider) OnClientDisconnect(si *hurlean.ServerInstance, id uint32) {
	
	fmt.Println("Client (id", id, ") has disconnected from the server", si)
}

func (fp MyServerFunctionalityProvider) OnClientMessage(si *hurlean.ServerInstance, id uint32, message hurlean.Message) {
	
	fmt.Println("")
	fmt.Println("----------------")
	fmt.Println("Message from", id, ":")
	fmt.Println("  Type:", message.Type)
	fmt.Println("  Body:", message.Body)
	fmt.Println("----------------")
	
	if message.Type == "chat message" {
		var nickname string
		if name, ok := fp.clientNames[id]; ok {
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
		fp.clientNamesMutex.Lock()
		fp.clientNames[id] = message.Body
		fp.clientNamesMutex.Unlock()
		
		fmt.Printf("client %v has set name to %v", id, message.Body)
	}
}

func (fp MyServerFunctionalityProvider) OnServerInit(serverInstance *hurlean.ServerInstance) {
	
	fmt.Printf("The server has been initialized!\n")
}

func (fp MyServerFunctionalityProvider) OnServerUpdate(serverInstance *hurlean.ServerInstance) {
	
	scanner := fp.scanner
	
	if scanner.Scan() {
		switch (scanner.Text()) {
		case "exit":
			serverInstance.Stop()
			
		case "disconnect":
			serverInstance.DisconnectAll()
		}
	}
}


func main() {
	
	// set the app-specific server's state
	var myServerFunctionalityProvider MyServerFunctionalityProvider = MyServerFunctionalityProvider{
		scanner: bufio.NewScanner(os.Stdin),
		clientNames: make(map[uint32]string),
		clientNamesMutex: sync.RWMutex{},
	}
	
	if err := hurlean.StartServer("8080", myServerFunctionalityProvider); err != nil {
		fmt.Println(err)
	}
}