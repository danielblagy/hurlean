// this file is not a part of the library, it's ised to test the execution

package main


import (
	"github.com/danielblagy/hurlean"
	"fmt"
)


type MyServerMessageHandler struct{}

func (mh MyServerMessageHandler) OnServerMessage(message []byte) {
	
	fmt.Println(string(message))
}


func main() {
	
	var mh hurlean.ServerMessageHandler = MyServerMessageHandler{}
	
	if err := hurlean.ConnectToServer("localhost", 8080, mh); err != nil {
		fmt.Println(err)
	}
}