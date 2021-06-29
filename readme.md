# hurlean

## TCP Networking Framework

## Installing

```
$ go get github.com/danielblagy/hurlean
```

## Import

```golang
import(
	"github.com/danielblagy/hurlean"
	//...
)
```

## Console Chat Example

[Chat Server Application](test_app_server/main.go)\
[Chat Client Application](test_app_client/main.go)

## How To Use

Let's create a simple server and client apps.
The client app will be able to query the server for current time.

### Example Server App

Let's start by creating **the server program**

First we import packages that we'll need

```golang
package main

import (
	"github.com/danielblagy/hurlean"
	"fmt"
	"time"
	"bufio"
	"os"
)
```

Then we implement `hurlean.ClientHandler` interface.\
With this interface we define the behavior of our server in the case of any client activity

```golang
type ExampleClientHandler struct{}

// executed once for each client
func (ch ExampleClientHandler) OnClientConnect(si *hurlean.ServerInstance, id uint32) {
	
	fmt.Printf("Client %v has connected to the server\n", id)
}

// executed once for each client
func (ch ExampleClientHandler) OnClientDisconnect(si *hurlean.ServerInstance, id uint32) {
	
	fmt.Printf("Client %v has disconnected from the server\n", id)
}

// executed each time the server gets a message from a client
func (ch ExampleClientHandler) OnClientMessage(si *hurlean.ServerInstance, id uint32, message hurlean.Message) {
	
	fmt.Printf("Message from %v: %v\n", id, message)
	
	if message.Type == "get current time" {
		currentTimeMessage := hurlean.Message{
			Type: "time",
			Body: time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006"),	// current time
		}
		si.Send(id, currentTimeMessage)
		
		fmt.Printf("Current time has been sent to %v\n", id)
	} else {
		fmt.Printf("Unknown message type '%v'\n", message.Type)
	}
}
```

We proceed by implementing `hurlean.ServerUpdater` interface.\
It allows us to specify the logic of our server program, e.g. getting console input which controls the server.

```golang
type ExampleServerUpdater struct{}

// executed continuously in a loop when the server is running
func (su ExampleServerUpdater) OnServerUpdate(serverInstance *hurlean.ServerInstance) {
	
	// get data from the server state
	scanner := serverInstance.State.(ExampleServerState).scanner
	
	// get console input from the server administrator (the app user)
	if scanner.Scan() {
		switch (scanner.Text()) {
		case "exit":
			serverInstance.Stop()
			
		case "disconnect":
			serverInstance.DisconnectAll()
		}
	}
}
```

We can define the state of our server program, which will be accessible via `hurlean.ServerInstance` pointer in
`hurlean.ClientHandler` methods implementations.

State can be anything from int or string to struct.

In this simple example we won't need much, only a *bufio.Scanner variable to store the initialized console scanner.

```golang
// used to store application-specific data in the server instance
type ExampleServerState struct{
	scanner *bufio.Scanner
}
```

Now, to the main function. Here we start the server on port 4545 and provide interfaces implementations and the state.\
In the case of failure to start the server, error will be returned.

```golang
func main() {
	
	// init server state
	var exampleServerState ExampleServerState = ExampleServerState{
		scanner: bufio.NewScanner(os.Stdin),
	}
	
	if err := hurlean.StartServer(4545, ExampleClientHandler{}, ExampleServerUpdater{}, exampleServerState); err != nil {
		fmt.Println(err)
	}
}
```

[Full source](#example-server-app-full-source)

### Example Client App

Import packages first

```golang
package main

import (
	"github.com/danielblagy/hurlean"
	"fmt"
	"bufio"
	"os"
)
```

Then we need to implement `hurlean.ServerMessageHandler` interface.\
Here we define how we should respond to the messages coming from the server.

```golang
type ExampleServerMessageHandler struct{}

// executed each time the client gets a message from the server
func (mh ExampleServerMessageHandler) OnServerMessage(message hurlean.Message) {
	
	if message.Type == "time" {
		fmt.Printf("Current time: %v\n\n", message.Body)
	} else {
		fmt.Printf("Unknown message type '%v'", message.Type)
	}
}
```

We proceed by implementing `hurlean.ClientUpdater` interface.\
It allows us to specify the logic of the client program, e.g. getting console input from the user.

```golang
type ExampleClientUpdater struct{}

// executed continuously in a loop when the client is running
func (cu ExampleClientUpdater) OnClientUpdate(clientInstance *hurlean.ClientInstance) {
	
	// get data from the server state
	scanner := clientInstance.State.(ExampleClientState).scanner
	
	// get console input from the user
	if scanner.Scan() {
		switch (scanner.Text()) {
		case "time":
			getTimeMessage := hurlean.Message{
				Type: "get current time",
				Body: "",
			}
			clientInstance.Send(getTimeMessage)
			
		case "disconnect":
			clientInstance.Disconnect()
		}
	}
}
```

Much like with server state, we can define client state as well.\
It'll be accessible via `hurlean.ClientInstance` pointer in `hurlean.ServerMessageHandler` method implementation.

```golang
// used to store application-specific data in the client instance
type ExampleClientState struct{
	scanner *bufio.Scanner
}
```

And finally, the main function where we start the client application and provide interfaces implementations and the state.\
In the case of failure to connect to the server, error will be returned.

``` golang
func main() {
	
	// init client state
	var exampleClientState ExampleClientState = ExampleClientState{
		scanner: bufio.NewScanner(os.Stdin),
	}
	
	if err := hurlean.ConnectToServer(
		"localhost", 4545, ExampleServerMessageHandler{}, ExampleClientUpdater{}, exampleClientState); err != nil {
		fmt.Println(err)
	}
}
```

[Full source](#example-client-app-full-source)


### Example Server App Full Source

```golang
package main


import (
	"github.com/danielblagy/hurlean"
	"fmt"
	"time"
	"bufio"
	"os"
)


type ExampleClientHandler struct{}

// executed once for each client
func (ch ExampleClientHandler) OnClientConnect(si *hurlean.ServerInstance, id uint32) {
	
	fmt.Printf("Client %v has connected to the server\n", id)
}

// executed once for each client
func (ch ExampleClientHandler) OnClientDisconnect(si *hurlean.ServerInstance, id uint32) {
	
	fmt.Printf("Client %v has disconnected from the server\n", id)
}

// executed each time the server gets a message from a client
func (ch ExampleClientHandler) OnClientMessage(si *hurlean.ServerInstance, id uint32, message hurlean.Message) {
	
	fmt.Printf("Message from %v: %v\n", id, message)
	
	if message.Type == "get current time" {
		currentTimeMessage := hurlean.Message{
			Type: "time",
			Body: time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006"),	// current time
		}
		si.Send(id, currentTimeMessage)
		
		fmt.Printf("Current time has been sent to %v\n", id)
	} else {
		fmt.Printf("Unknown message type '%v'\n", message.Type)
	}
}


type ExampleServerUpdater struct{}

// executed continuously in a loop when the server is running
func (su ExampleServerUpdater) OnServerUpdate(serverInstance *hurlean.ServerInstance) {
	
	// get data from the server state
	scanner := serverInstance.State.(ExampleServerState).scanner
	
	// get console input from the server administrator (the app user)
	if scanner.Scan() {
		switch (scanner.Text()) {
		case "exit":
			serverInstance.Stop()
			
		case "disconnect":
			serverInstance.DisconnectAll()
		}
	}
}


// used to store application-specific data in the server instance
type ExampleServerState struct{
	scanner *bufio.Scanner
}


func main() {
	
	// init server state
	var exampleServerState ExampleServerState = ExampleServerState{
		scanner: bufio.NewScanner(os.Stdin),
	}
	
	if err := hurlean.StartServer(4545, ExampleClientHandler{}, ExampleServerUpdater{}, exampleServerState); err != nil {
		fmt.Println(err)
	}
}
```

### Example Client App Full Source

```golang
package main


import (
	"github.com/danielblagy/hurlean"
	"fmt"
	"bufio"
	"os"
)


type ExampleServerMessageHandler struct{}

// executed each time the client gets a message from the server
func (mh ExampleServerMessageHandler) OnServerMessage(message hurlean.Message) {
	
	if message.Type == "time" {
		fmt.Printf("Current time: %v\n\n", message.Body)
	} else {
		fmt.Printf("Unknown message type '%v'", message.Type)
	}
}


type ExampleClientUpdater struct{}

// executed continuously in a loop when the client is running
func (cu ExampleClientUpdater) OnClientUpdate(clientInstance *hurlean.ClientInstance) {
	
	// get data from the server state
	scanner := clientInstance.State.(ExampleClientState).scanner
	
	// get console input from the user
	if scanner.Scan() {
		switch (scanner.Text()) {
		case "time":
			getTimeMessage := hurlean.Message{
				Type: "get current time",
				Body: "",
			}
			clientInstance.Send(getTimeMessage)
			
		case "disconnect":
			clientInstance.Disconnect()
		}
	}
}


// used to store application-specific data in the client instance
type ExampleClientState struct{
	scanner *bufio.Scanner
}


func main() {
	
	// init client state
	var exampleClientState ExampleClientState = ExampleClientState{
		scanner: bufio.NewScanner(os.Stdin),
	}
	
	if err := hurlean.ConnectToServer(
		"localhost", 4545, ExampleServerMessageHandler{}, ExampleClientUpdater{}, exampleClientState); err != nil {
		fmt.Println(err)
	}
}
```


For more complex example check out [the console chat example](#console-chat-example)