package main

import "D7024E/dht"

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

	id := strings.TrimSpace(Id)
	ip := strings.TrimSpace(Ip)
	port := strings.TrimSpace(Port)

	n := dht.MakeDHTNode(&id, ip, port)
	//n.JoinRing("localhost:1112")
	fmt.Println("penis1")
	for {
		fmt.Println("Enter command: ")
		Input, _ := reader.ReadString('\n')
		fmt.Println("penis2")
		fmt.Scanln(Input)
		fmt.Println("penis3")
		input := strings.TrimSpace(Input)
		fmt.Println("penis4")

		switch input {
		//		case "join":
		//			go n.Join(input)

		case "joinRing":
			go n.JoinRing("localhost:1112")

			//		case "changePredecessor":
			//			go n.changePredecessor(input)
		case "fingers":
			go n.FingerPrint()

		}

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
