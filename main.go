package main

import "D7024E/dht"

//import "D7024E/dht"
import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	flag.Parse()
	reader := bufio.NewReader(os.Stdin)

	//id := flag.String(name, value, usage)
	id := flag.Arg(0)
	port := flag.Arg(1)
	//fmt.Print("Enter id: ")
	//	Id, _ := reader.ReadString('\n')
	//Id := ""
	//fmt.Println(text)

	//fmt.Println("Enter Ip: ")
	//Ip, _ := reader.ReadString('\n')
	//Ip := "localhost"
	//fmt.Scanln(Ip)

	//fmt.Println("Enter port: ")
	//Port, _ := reader.ReadString('\n')
	//fmt.Scanln(Port)

	//id0 := "00"
	//fmt.Println(Id)
	//fmt.Println(Ip)
	//fmt.Println(Port)

	//id := strings.TrimSpace(Id)
	//ip := strings.TrimSpace(Ip)
	//port := strings.TrimSpace(Port)
	fmt.Println("Detta är ditt id:", id)
	fmt.Println("Detta är din port: ", port)
	ip := "localhost"

	n := dht.MakeDHTNode(&id, ip, port)
	go n.JoinRing("localhost:1111")

	go func() {

		http.HandleFunc("/chord/", dht.Chord)

		http.HandleFunc("/chord/post/", func(w http.ResponseWriter, r *http.Request) {
			dht.Post(w, r, n)
		})

		http.HandleFunc("/chord/get/", func(w http.ResponseWriter, r *http.Request) {
			dht.Get(w, r, n)
		})

		http.HandleFunc("/chord/put/", func(w http.ResponseWriter, r *http.Request) {
			dht.Put(w, r, n)
		})

		http.HandleFunc("/chord/delete/", func(w http.ResponseWriter, r *http.Request) {
			dht.Del(w, r, n)
		})

		http.HandleFunc("/chord/list/", func(w http.ResponseWriter, r *http.Request) {
			dht.List(w, r, n)
		})

		http.ListenAndServe(":"+port, nil)

		fmt.Println("The page is rolling")

	}()

	go func() {
		c := time.Tick(3 * time.Second)
		for {
			select {
			case <-c:
				n.AutoFingers()
				//node.autoFingers()
			}
		}
	}()

	for {
		fmt.Println("Enter command: ")
		Input, _ := reader.ReadString('\n')

		fmt.Scanln(Input)
		input := strings.TrimSpace(Input)

		switch input {
		//		case "join":
		//			go n.Join(input)

		case "joinRing":
			go n.JoinRing("localhost:1112")

			//		case "changePredecessor":
			//			go n.changePredecessor(input)

		case "id":
			go n.IdPrint()

		case "preid":
			go n.PreID()

		case "sucid":
			go n.SucID()

		case "ping":
			fmt.Println("hate")
			go n.Ping()

		}

	}

}
