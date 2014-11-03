package dht

import (
	"encoding/json"
	"fmt"
	//	"math/big"
	"net"
	//	"strings"
	//	"testing"
	//	"time"
)

//###################################//
//									//
// Nätverk och dess funktioner	   //
//								  //
//###############################//
//
//struct for messages that we got from the lab handout

/*
vad är det vi vill skicka?
vem vi är?| vad vi vill att den ska göra/utföra? |
*/
type Msg struct {
	// Type = metoden som skall köras
	// KEY = värdet som skall köras
	// Src = noden som kallade (den som skickar meddelandet)
	// Dst = destinationsadressen
	// Origin = vem var det som ropa från början?? Vem var det??!!
	// Time = timestamp
	Type, Key, Src, Dst, Origin string
	Time                        int64
}

//struct for Transport from lab handout
/*

vi vill ta emot ett meddelande.
läsa ut meddelandet och sedan returnera svaret till source addressen

*/
type Transport struct {
	node        *DHTNode
	bindAddress string
	channel     map[int64]chan Msg

	// chan,,,, mutexlås
}

// listen function from lab handout
func (transport *Transport) listen() {
	udpAddr, err := net.ResolveUDPAddr("udp", transport.bindAddress)
	if err != nil {
		fmt.Println("Error in listen func: ", err)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error in listen func: ", err)
	}
	defer conn.Close()
	dec := json.NewDecoder(conn)
	for {
		msg := Msg{}
		err := dec.Decode(&msg)
		if err != nil {
			fmt.Println("Error in listen func: ", err)
		}
		//Parse(msg)
		// if type is response check timestamp and call the channel
		//we got a message maby baby?
		//Parse vad det är för metod (lookup, addToring)

	}

}

// send function from lab handout
func (transport *Transport) send(msg *Msg, ch chan Msg) {
	if ch != nil {
		transport.channel[msg.Time] = ch
	}

	udpAddr, err := net.ResolveUDPAddr("udp", msg.Dst)
	if err != nil {
		fmt.Println("Error in send func: ", err)
	}
	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error in send func: ", err)
	}
	defer conn.Close()
	_, err = conn.Write(msg.Bytes())
	if err != nil {
		fmt.Println("Error in send func: ", err)
	}
	//implementera msg.Bytes
	//encoda till ett json object
	//få det till en bytearray
	//alltså en bytearray som representerar ett json objekt

}

func (msg *Msg) Bytes() []byte {
	//encode to json
	jsonenconded, err := json.Marshal(msg)
	if err == nil {
		return jsonenconded
	}
	fmt.Println("Error in Bytes func: ", err)
	return nil

}

func (n *DHTNode) parse(msg *Msg) {

	msg := make(chan *Msg)
	//go n.Transport.listen(msg)
	go n.Transport.listen()
	switch msg.Type {

	case "join":
		go n.join(msg)
		break
	case "joinRing":
		go n.joinRing(networkaddr)
		break
	case "changePredecessor":
		go n.changePredecessor(msg)
		break
	case "findSuccessorFinger":
		go n.findSuccessorFinger()
		break
	}
}
