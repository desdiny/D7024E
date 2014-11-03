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

	if &Ip != nil {
		dht.MakeDHTNode(&id, ip, port)

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
