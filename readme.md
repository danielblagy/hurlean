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

Then we implement `hurlean.ServerFunctionalityProvider` interface.\
With this interface we define the behavior of our server application.\

First we define the struct of our interface implementation. Here we can put our client application state.\
In this simple example we won't need much, only a *bufio.Scanner variable to store the initialized console scanner.
We the scanner we can get console input from the user of our server application. We'll provide two commands: 'exit',
which will stop the server, and 'disconnect' which will force-disconnect all the currently connected users.

```golang
type ExampleClientHandler struct{}

type ExampleServerFunctionalityProvider struct{
	// here we can store application-specific data
	scanner *bufio.Scanner
}
```

Next we implement functions that define the behavior of our server in the case of any client activity.

```golang
// executed once for each client
func (fp ExampleServerFunctionalityProvider) OnClientConnect(si *hurlean.ServerInstance, id uint32) {
	
	fmt.Printf("Client %v has connected to the server\n", id)
}

// executed once for each client
func (fp ExampleServerFunctionalityProvider) OnClientDisconnect(si *hurlean.ServerInstance, id uint32) {
	
	fmt.Printf("Client %v has disconnected from the server\n", id)
}

// executed each time the server gets a message from a client
func (fp ExampleServerFunctionalityProvider) OnClientMessage(si *hurlean.ServerInstance, id uint32, message hurlean.Message) {
	
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

Now let's implement the functions that define the logic of our server program.

```golang
// executed once when the server instance is initialized
func (fp ExampleServerFunctionalityProvider) OnServerInit(serverInstance *hurlean.ServerInstance) {
	// empty for this simple example
}

// executed continuously in a loop when the server is running
func (fp ExampleServerFunctionalityProvider) OnServerUpdate(serverInstance *hurlean.ServerInstance) {
	
	// get console input from the server administrator (the app user)
	if fp.scanner.Scan() {
		switch (fp.scanner.Text()) {
		case "exit":
			serverInstance.Stop()
			
		case "disconnect":
			serverInstance.DisconnectAll()
		}
	}
}
```

Now, to the main function. Here we start the server on port 4545 and provide interfaces implementations and the state.\
In the case of failure to start the server, error will be returned.

```golang
func main() {
	
	// init server state
	exampleServerFunctionalityProvider := ExampleServerFunctionalityProvider{
		scanner: bufio.NewScanner(os.Stdin),
	}
	
	if err := hurlean.StartServer("4545", exampleServerFunctionalityProvider); err != nil {
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

Then we need to implement `hurlean.ClientFunctionalityProvider` interface.\
Much like with server state, we're going to define client state as well.

```golang
type ExampleClientFunctionalityProvider struct{
	// here we can store application-specific data
	scanner *bufio.Scanner
}
```

We proceed by implementing the interface functions.\
Here we define how we should respond to the messages coming from the server.

```golang
// executed each time the client gets a message from the server
func (fp ExampleClientFunctionalityProvider) OnServerMessage(clientInstance *hurlean.ClientInstance, message hurlean.Message) {
	
	if message.Type == "time" {
		fmt.Printf("Current time: %v\n\n", message.Body)
	} else {
		fmt.Printf("Unknown message type '%v'", message.Type)
	}
}
```

Now we need to implement functions that define the client application behavior.

```golang
// executed once when the client instance is initialized
func (fp ExampleClientFunctionalityProvider) OnClientInit(clientInstance *hurlean.ClientInstance) {
	// empty for this simple example
}

// executed continuously in a loop when the client is running
func (fp ExampleClientFunctionalityProvider) OnClientUpdate(clientInstance *hurlean.ClientInstance) {
	
	// get console input from the user
	if fp.scanner.Scan() {
		switch (fp.scanner.Text()) {
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

And finally, the main function where we start the client application.\
In the case of failure to connect to the server, error will be returned.

``` golang
func main() {
	
	// init client state
	exampleClientFunctionalityProvider := ExampleClientFunctionalityProvider{
		scanner: bufio.NewScanner(os.Stdin),
	}
	
	if err := hurlean.ConnectToServer(
		"localhost", "4545", exampleClientFunctionalityProvider); err != nil {
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


type ExampleServerFunctionalityProvider struct{
	// here we can store application-specific data
	scanner *bufio.Scanner
}

// executed once for each client
func (fp ExampleServerFunctionalityProvider) OnClientConnect(si *hurlean.ServerInstance, id uint32) {
	
	fmt.Printf("Client %v has connected to the server\n", id)
}

// executed once for each client
func (fp ExampleServerFunctionalityProvider) OnClientDisconnect(si *hurlean.ServerInstance, id uint32) {
	
	fmt.Printf("Client %v has disconnected from the server\n", id)
}

// executed each time the server gets a message from a client
func (fp ExampleServerFunctionalityProvider) OnClientMessage(si *hurlean.ServerInstance, id uint32, message hurlean.Message) {
	
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

// executed once when the server instance is initialized
func (fp ExampleServerFunctionalityProvider) OnServerInit(serverInstance *hurlean.ServerInstance) {
	// empty for this simple example
}

// executed continuously in a loop when the server is running
func (fp ExampleServerFunctionalityProvider) OnServerUpdate(serverInstance *hurlean.ServerInstance) {
	
	// get console input from the server administrator (the app user)
	if fp.scanner.Scan() {
		switch (fp.scanner.Text()) {
		case "exit":
			serverInstance.Stop()
			
		case "disconnect":
			serverInstance.DisconnectAll()
		}
	}
}


func main() {
	
	// init server state
	exampleServerFunctionalityProvider := ExampleServerFunctionalityProvider{
		scanner: bufio.NewScanner(os.Stdin),
	}
	
	if err := hurlean.StartServer("4545", exampleServerFunctionalityProvider); err != nil {
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


type ExampleClientFunctionalityProvider struct{
	// here we can store application-specific data
	scanner *bufio.Scanner
}

// executed each time the client gets a message from the server
func (fp ExampleClientFunctionalityProvider) OnServerMessage(clientInstance *hurlean.ClientInstance, message hurlean.Message) {
	
	if message.Type == "time" {
		fmt.Printf("Current time: %v\n\n", message.Body)
	} else {
		fmt.Printf("Unknown message type '%v'", message.Type)
	}
}

// executed once when the client instance is initialized
func (fp ExampleClientFunctionalityProvider) OnClientInit(clientInstance *hurlean.ClientInstance) {
	// empty for this simple example
}

// executed continuously in a loop when the client is running
func (fp ExampleClientFunctionalityProvider) OnClientUpdate(clientInstance *hurlean.ClientInstance) {
	
	// get console input from the user
	if fp.scanner.Scan() {
		switch (fp.scanner.Text()) {
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


func main() {
	
	// init client state
	exampleClientFunctionalityProvider := ExampleClientFunctionalityProvider{
		scanner: bufio.NewScanner(os.Stdin),
	}
	
	if err := hurlean.ConnectToServer(
		"localhost", "4545", exampleClientFunctionalityProvider); err != nil {
		fmt.Println(err)
	}
}
```


For more complex example check out [the console chat example](#console-chat-example)