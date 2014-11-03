package main

import "D7024E/dht"

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter id: ")
	text, _ := reader.ReadString('\n')
	Id := ""
	fmt.Println(text)

	fmt.Println("Enter port: ")
	Port := ""
	fmt.Scanln(Port)

	fmt.Println("Enter Ip: ")
	Ip := ""
	fmt.Scanln(Ip)

	n := dht.MakeDHTNode(id, Ip, Port)

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
