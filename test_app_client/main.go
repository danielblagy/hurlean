// this file is not a part of the library, it's ised to test the execution

package main


import (
	"github.com/danielblagy/hurlean"
	"fmt"
	"bufio"
	"os"
	"strings"
)


type MyServerMessageHandler struct{}

func (mh MyServerMessageHandler) OnServerMessage(message hurlean.Message) {
	
	if message.Type == "chat message" {
		fmt.Printf("%v\n\n", message.Body)
	}
}


type MyClientUpdater struct{}

func (cu MyClientUpdater) OnClientUpdate(clientInstance *hurlean.ClientInstance) {
	
	scanner := clientInstance.State.(MyClientState).scanner
	
	if scanner.Scan() {
		input := scanner.Text()
		
		if len(input) > 0 && input[0] == '/' {
			input = strings.TrimPrefix(input, "/")
			input := strings.Split(input, " ")
			
			switch (input[0]) {
			case "disconnect":
				clientInstance.Disconnect()
			case "setname":
				message := hurlean.Message{
					Type: "setname",
					Body: input[1],
				}
				clientInstance.Send(message)
			default:
				fmt.Printf("Unrecognized command '%v'\n", input[0])
			}
		} else {
			message := hurlean.Message{
				Type: "chat message",
				Body: input,
			}
			clientInstance.Send(message)
		}
	}
}


type MyClientState struct{
	scanner *bufio.Scanner
}


func main() {
	
	var myServerMessageHandler hurlean.ServerMessageHandler = MyServerMessageHandler{}
	var myClientUpdater hurlean.ClientUpdater = MyClientUpdater{}
	
	// set the app-specific client's state
	var myClientState MyClientState = MyClientState{
		scanner: bufio.NewScanner(os.Stdin),
	}
	
	if err := hurlean.ConnectToServer("localhost", 8080, myServerMessageHandler, myClientUpdater, myClientState); err != nil {
		fmt.Println(err)
	}
}