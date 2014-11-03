package main

import "D7024E/dht"

//import "fmt"

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
