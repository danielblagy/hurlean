// this file is not a part of the library, it's ised to test the execution

package main


import (
	"github.com/danielblagy/hurlean"
	"fmt"
)


func main() {
	if err := hurlean.StartServer(8080); err != nil {
		fmt.Println(err)
	}
}