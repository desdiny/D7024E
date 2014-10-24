package main

//import "D7024E/dht"
import "D7024E.git/branches/Objective-2/dht"
import (
	"fmt"
	"time"
)

func main() {
	id0 := "00"
	n := dht.MakeDHTNode(&id0, "aaa", "aaa")

	go func() {
		c := time.Tick(3 * time.Second)
		for {
			select {
			case <-c:
				//node.autoFingers()
			}
		}
	}()

}
