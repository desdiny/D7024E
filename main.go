package main

import "D7024E.git/branches/Objective-2/dht"

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter id: ")
	Id, _ := reader.ReadString('\n')
	//Id := ""
	//fmt.Println(text)

	fmt.Println("Enter Ip: ")
	Ip, _ := reader.ReadString('\n')
	fmt.Scanln(Ip)

	fmt.Println("Enter port: ")
	Port, _ := reader.ReadString('\n')
	fmt.Scanln(Port)

	//id0 := "00"
	fmt.Println(Id)
	fmt.Println(Ip)
	fmt.Println(Port)

	if &Ip != nil {
		dht.MakeDHTNode(&Id, Ip, Port)

	}

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
