package hurlean


// hurlean.Message objects is data that can be sent and received via network
type Message struct {
	Type string
	Body string
}


// controls the debug prints
var debug bool = true