// this file is not a part of the library, it's ised to test the execution

package main


import (
	"github.com/danielblagy/hurlean"
	"fmt"
	"bufio"
	"os"
	"strings"
)


type MyClientFunctionalityProvider struct{
	scanner *bufio.Scanner
}

func (fp MyClientFunctionalityProvider) OnServerMessage(clientInstance *hurlean.ClientInstance, message hurlean.Message) {
	
	if message.Type == "chat message" {
		fmt.Printf("%v\n\n", message.Body)
	}
}

func (fp MyClientFunctionalityProvider) OnClientInit(clientInstance *hurlean.ClientInstance) {
	
	fmt.Printf("Welcome to the chat!\n\n")
}

func (fp MyClientFunctionalityProvider) OnClientUpdate(clientInstance *hurlean.ClientInstance) {
	
	scanner := fp.scanner
	
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


func main() {
	
	// set the app-specific client's state
	var myClientFunctionalityProvider MyClientFunctionalityProvider = MyClientFunctionalityProvider{
		scanner: bufio.NewScanner(os.Stdin),
	}
	
	if err := hurlean.ConnectToServer("localhost", "8080", myClientFunctionalityProvider); err != nil {
		fmt.Println(err)
	}
}