package dht

import (
	//"encoding/json" // used for networking
	"fmt"
	"math/big"  // used for fingers
	"math/rand" //used for updating fingers
	//"net"
	"strings"
	//"testing"
	"time" // used to update fingers and to set time for msg
)

//###################################//
//									//
// DHT NODER OCH DESS FUNKTIONER     //
//								      //
//###################################//

//######################################//
//										//
// Denna DHT NOD SOM sätts i porgrammet	//
//										//
//######################################//

type DHTNode struct {
	id, address, port      string
	successor, predecessor *DHTNode
	//successor, predecessor *OutsideNode
	finger    []*Fingers //links to Fingers struct
	Transport *Transport
}

//##############################//
//								//
// Noder som sätts från utsidan	//
//								//
//##############################//

//type OutsideNode struct {
//	id, address, port string
//}

//added Fingers struct.. we say that every DHTNODE have finger witch is
// populated by fingers (ie. a start string and a pointer to a DHTNODE)
//so a DHTNode will now look like this:
//
//		id:00 address:nil port:nil
//		successor:01 predecessor:09
//		finger [start,node],[start,node],[start,node]
type Fingers struct {
	start string
	node  *DHTNode
}

func (node *DHTNode) autoFingers() {
	channel := make(chan Msg)
	i := rand.Intn(3) //vet inte ifall det behövs en random var i intn(???)

	//create autofingers message
	m := makeMsg("lookupNetwork", node.Address(), node.finger[i].start, node.Address(), time.Now().UnixNano(), node.Address())
	node.Transport.send(m, channel)
	//waitning for answer
	req := <-channel
	//split address and port
	a := strings.Split(req.Src, ":")
	node.finger[i].node.id = req.Key
	node.finger[i].node.address = a[0]
	node.finger[i].node.port = a[1]
	//finger := node.lookupNetwork(node.finger[i].node.id)

	//if finger != nil {
	//	node.finger[i].node = finger
	//}

}

//#############################//
//							   //
//		Make local node		   //
//							   //
//#############################//
func MakeDHTNode(idcheck *string, address string, port string) *DHTNode {
	n := new(DHTNode)
	if idcheck == nil {
		n.id = generateNodeId()
		n.address = address
		n.port = port
		n.successor = n
		n.predecessor = n
		n.finger = make([]*Fingers, 3) //change to use for 3 and 160
		n.initFingerTable()
		n.Transport = makeTransport(n, n.address, n.port)
		go n.Transport.listen()

	} else {
		n.id = *idcheck
		n.address = address
		n.port = port
		n.successor = n
		n.predecessor = n
		n.finger = make([]*Fingers, 3) //change to use for 3 and 160
		n.initFingerTable()
		n.Transport = makeTransport(n, n.address, n.port)
		go n.Transport.listen()

	}
	return n

}

func (n *DHTNode) initFingerTable() {

	// fixar fingrar där  för att fylla på med nollor på rätt ställen etc
	for i := 1; i <= len(n.finger); i++ {
		fingerID, _ := calcFinger([]byte(n.id), i, len(n.finger))
		if len(fingerID) < len(n.id) {
			fingerID = strings.Repeat("0", len(n.id)-len(fingerID)) + fingerID
		}
		tempfinger := new(Fingers)
		tempfinger.start = fingerID
		tempfinger.node = n
		n.finger[i-1] = tempfinger
	}
	return

}
func makeMsg(Type string, Dst string, Key string, Origin string, Time int64, Src string) *Msg {
	m := new(Msg)
	m.Type = Type
	m.Dst = Dst
	m.Key = Key
	m.Origin = Origin
	m.Time = Time
	m.Src = Src
	return m

}

func makeTransport(node *DHTNode, Address string, port string) *Transport {
	s := new(Transport)
	s.node = node

	s.bindAddress = Address + ":" + port
	s.channel = make(map[int64]chan Msg)
	return s
}

/////////////////////////////////////////////////
////////////////////////////////////////////
// new func for addToRing for networking  //
////////////////////////////////////////////
/////////////////////////////////////////////

//node that wants to join ring
func (n *DHTNode) JoinRing(networkaddr string) {
	channel := make(chan Msg)
	fmt.Println("JoinRing in Progress")
	fmt.Println("calling node on address: ", networkaddr)
	m := makeMsg("lookupNetwork", networkaddr, n.id, n.Address(), time.Now().UnixNano(), n.Address())
	n.Transport.send(m, channel)
	fmt.Println("node has been called")

	req := <-channel
	fmt.Println("receved a answer on JoinRing with KEY: ", req.Key)
	joinidandaddress := n.id + "," + n.address
	m = makeMsg("join", req.Src, joinidandaddress, n.Address(), time.Now().UnixNano(), n.Address())
	n.Transport.send(m, channel)
	fmt.Println("Sending message to join function....")

	// waiting for answer
	req = <-channel
	fmt.Println("recived ansver from Joinfunction")

	// split req (id and address)
	a := strings.Split(req.Key, ",")
	// create a new node
	s := new(DHTNode)
	b := strings.Split(req.Src, ":")
	//s := new(OutsideNode)
	s.id = a[0]
	s.address = a[1]
	s.port = b[1]
	fmt.Println("added the predecessor with ip: ", s.address, "id: ", s.id, "port: ", s.port)
	n.predecessor = s
	fmt.Println("Ending JoinRing")
	fmt.Println("---------------------------------------")
	//inte än fixat
	/*n.initFingerTable(newnode)

	//contacts node in ring
		node := n.lookup(newnode.id)
		oldnode := node.successor
		node.successor = newnode
		newnode.successor = oldnode
		newnode.predecessor = node
		oldnode.predecessor = newnode
		n.update_others()
	*/
}

// the node that jumps on the node
func (n *DHTNode) Join(msg *Msg) {
	fmt.Println("-----------------------------------------")
	fmt.Println("Join operation initializing")
	fmt.Println("")
	//channel := make(chan Msg)
	// splits the incomming keys
	a := strings.Split(msg.Key, ",")
	k := strings.Split(msg.Src, ":")

	fmt.Println("the joining has begun, calling to set predecessor on next node")
	fmt.Println("")
	oldsuccessor := n.successor
	joinidandaddress := a[0] + "," + a[1] + "," + k[1]
	m := makeMsg("changePredecessor", n.successor.Address(), joinidandaddress, n.Address(), time.Now().UnixNano(), n.Address())
	n.Transport.send(m, nil)
	fmt.Println("Sending changePredecessor to successor")
	fmt.Println("")
	//creates a new node
	//s := new(DHTNode)
	v := strings.Split(msg.Origin, ":")
	//s.id = a[0]
	//	s.address = a[1]
	//s.port = v[1]
	n.successor.id = a[0]
	n.successor.address = a[1]
	n.successor.port = v[1]
	fmt.Println("added the new node: ", " id: ", n.successor.id, " address: ", n.successor.address, " port: ", n.successor.port)
	fmt.Println("")
	// adds the new node as the nodes succsessor
	//	n.successor = s
	n.finger[0].node = n.successor

	//adding both to one variable so we can send it in the key value
	// have to concatinate when message is recived
	joinidandaddress = n.id + "," + n.address

	//creates message
	m = makeMsg("response", msg.Origin, joinidandaddress, n.Address(), msg.Time, n.Address())

	// sends message
	n.Transport.send(m, nil)

	key := oldsuccessor.id + ":" + oldsuccessor.address + ":" + oldsuccessor.port
	m = makeMsg("changeSuccessor", msg.Origin, key, n.Address(), TimeNow(), n.Address())
	n.Transport.send(m, nil)

	fmt.Println("sends respons to JoinRing to let it add with our cred")
	fmt.Println("")
	fmt.Println("Join operation complete")
	fmt.Println("------------------------------------------")
}

func (n *DHTNode) changePredecessor(msg *Msg) {

	//split incomming key
	a := strings.Split(msg.Key, ",")

	//create a new node on this instance
	//s := new(DHTNode)
	//s := new(OutsideNode)
	//s.id = a[0]
	//s.address = a[1]
	//s.port = a[2]
	n.predecessor.id = a[0]
	n.predecessor.address = a[1]
	n.predecessor.port = a[2]

	// adds the node to n's predecessor
	//n.predecessor = s

	//m := makeMsg("changeSuccessor", s.Address(), n.id, n.Address(), msg.Time, n.Address())

	// sends message
	//	n.Transport.send(m, nil)

}

func (n *DHTNode) changeSuccessor(msg *Msg) {
	a := strings.Split(msg.Key, ":")
	n.successor.id = msg.Key
	n.successor.address = a[1]
	n.successor.port = a[0]
	//s := new(DHTNode)
	//s.id = msg.Key
	//a := strings.Split(msg.Key, ":")
	//s.id = a[0]
	//s.address = a[1]
	//s.port = a[2]
	//n.successor = s
	n.finger[0].node = n.successor
}

//func (n *DHTNode) addToRing(newnode *DHTNode) {
//	fmt.Println("Nodens id: ", newnode.id)
//	if n.finger[0] == nil {
// fixar fingrar special första gången
//		for i := 1; i <= len(n.finger); i++ {
//			fingerID, _ := calcFinger([]byte(n.id), i, len(n.finger))
//			if len(fingerID) < len(n.id) {
//				fingerID = strings.Repeat("0", len(n.id)-len(fingerID)) + fingerID
//			}
//			tempnode := n.lookup(fingerID)
//
//			if tempnode.id != fingerID {
//				tempnode = tempnode.successor
//
//			}
//			n.finger[i-1] = &Fingers{fingerID, tempnode}
//
//			fmt.Println(n.finger[i-1].node.id)
//		}

//	}
//nyinlaggt den 14/10 vet inte om jag tänker rätt
//n.initFingerTable(newnode)

// fixar fingrar där  för att fylla på med nollor på rätt ställen etc
//for i := 1; i <= len(n.finger); i++ {
//	fingerID, _ := calcFinger([]byte(newnode.id), i, len(n.finger))
//	if len(fingerID) < len(n.id) {
//		fingerID = strings.Repeat("0", len(n.id)-len(fingerID)) + fingerID
//	}
//	tempnode := n.lookup(fingerID)
//	if tempnode.id != fingerID {
//		tempnode = tempnode.successor
//
//		}
//		newnode.finger[i-1] = &Fingers{fingerID, tempnode}
//		fmt.Println(newnode.finger[i-1].node.id)

//	}
//skapa ett meddelande som skall köra lookup för vilken nod vi vill joina på
//då kör man join på den ringen

//	node := n.lookup(newnode.id)
//	oldnode := node.successor
//	node.successor = newnode
//	newnode.successor = oldnode
//	newnode.predecessor = node
//	oldnode.predecessor = newnode
//	newnode.update_others()

//}

func (n *DHTNode) printRing() {

	nextNode := n.successor
	fmt.Println("id: ", n.id, "fingers: ", n.finger)
	for nextNode != n {
		fmt.Printf("id: %s fingers: ", nextNode.id)
		for i := 0; i < len(nextNode.finger); i++ {
			fmt.Printf("%s ", nextNode.finger[i].node.id)

		}
		fmt.Println()

		//fmt.Println(nextNode.id)
		nextNode = nextNode.successor

	}
}

func (d *DHTNode) tostring() (out string) {
	out = "DHTNode{id: " + d.id + ", address: " + d.address + ", port: " + d.port + "}"

	return
}

//func (d *DHTNode) lookup(hash string) *DHTNode {

//	if between([]byte(d.id), []byte(d.successor.id), []byte(hash)) {
// returns that this node should be responible for this
// how to use type in this case?
// can we just send
//		makeMsg(, Dst, Key, Origin)
//		return d
//	}

//	dist := distance(d.id, hash, len(d.finger))
//	index := dist.BitLen() - 1
//	if index < 0 {
//		return d
//	}
//	fmt.Println("INDEX", index)

//stegar ner tills fingret inte pekar på sig själv
//	for ; index > 0 && d.finger[index].node == d; index-- {

//	}
// Kollar så vi inte hamnar för långt
//	diff := big.Int{}
//	diff.Sub(dist, distance(d.id, d.finger[index].node.id, len(d.finger)))
//	for index > 0 && diff.Sign() < 0 {
//		index--
//		diff.Sub(dist, distance(d.id, d.finger[index].node.id, len(d.finger)))
//	}
//kollar så vi inte pekar på oss själva
//	if d.finger[index].node == d || diff.Sign() < 0 {
//		fmt.Println("ERROR ERROR alles gebort auf the baut")
//		return d.successor.lookup(hash)

//	}
/* här skall vi alltså lägga in att hoppa till en annan nod med
   ett msg sedan skicka det msget till send
   a = den här noden vi är i
   b = noden som skall plaseras
   msget skall då alltså innehålla:

   Type = lookup
   KEY = b.id
   Src = a.ip
   Dst = fingerindex[x].ip
   Origin = b.ip

*/
//	return d.finger[index].node.lookup(hash)
/*
   här under har vi den förra funktionen för att köra utan fingrar
*/
//	return d.successor.lookup(hash)
//}

//om s är i (någon av) n  fingrar, uppdatera n's fingrar med s
//unc (n *DHTNode) update_finger_table(s *DHTNode, i int) {
//	fmt.Println("updating finger", i, "on", n.id)
//	if s.successor == n.finger[i-1].node {
//		n.finger[i-1].node = s
//		p := n.predecessor
//		if p != n {
//			p.update_finger_table(s, i)
//		}

//	}

//}

//////////////////////////////////////////////////////////
//				func for lookup 						//
//														//
//	if it cant find it here it uses fingers				//
//	if unsing fingers it will jump to the closest		//
//	to our hash and then run lookupNetwork on that one	//
//////////////////////////////////////////////////////////
// func (d *DHTNode) lookup(hash string) *DHTNode {
// 	channel := make(chan Msg)
// 	//if d is  responsible for id
// 	if between([]byte(d.id), []byte(d.successor.id), []byte(hash)) {
// 		//returns d
// 		return d
// 	}
// 	//otherwise use fingers
// 	dist := distance(d.id, hash, len(d.finger))
// 	index := dist.BitLen() - 1
// 	if index < 0 {
// 		return d
// 	}
// 	fmt.Println("INDEX", index)

// 	//stegar ner tills fingret inte pekar på sig själv
// 	for ; index > 0 && d.finger[index].node == d; index-- {

// 	}
// 	// Kollar så vi inte hamnar för långt
// 	diff := big.Int{}
// 	diff.Sub(dist, distance(d.id, d.finger[index].node.id, len(d.finger)))
// 	for index > 0 && diff.Sign() < 0 {
// 		index--
// 		diff.Sub(dist, distance(d.id, d.finger[index].node.id, len(d.finger)))
// 	}
// 	//kollar så vi inte pekar på oss själva
// 	if d.finger[index].node == d || diff.Sign() < 0 {
// 		fmt.Println("ERROR ERROR alles gebort auf the baut")
// 		// send message to the successor node to do a lookup
// 		m := makeMsg("lookupNetwork", d.successor.address, hash, d.address)
// 		d.Transport.send(m, channel)

// 		//väntar på att vi ska få tillbaka ett svar
// 		req := <-channel
// 		//får tillbaka en nod req
// 		//creates a new node
// 		s := new(DHTNode)
// 		//s := new(OutsideNode)
// 		s.id = req.Key
// 		s.address = req.Src

// 		return s

// 		//return d.successor.lookup(hash)

// 	}

// 	// if nothing of the above works
// 	m := makeMsg("lookupNetwork", d.finger[index].node.address, hash, d.address)
// 	d.Transport.send(m, channel)

// 	//chilling for response
// 	req := <-channel
// 	//create new node
// 	s := new(DHTNode)
// 	//s := new(OutsideNode)
// 	s.id = req.Key
// 	s.address = req.Src

// 	return s

// 	//return d.finger[index].node.lookup(hash)

// }

//////////////////////////////////////////////////////////
//				func for lookupNetwork					//
//														//
//	beeing called from ether lookup on another computer	//
//	or lookupNetwork on another computer				//
//	sends the query forward if this isnt the right node	//
// 	or answers to the node who ran the lookup req from 	//
//	the begining whit help of msg.Origin				//
//////////////////////////////////////////////////////////
//node contacted over network
func (d *DHTNode) lookupNetwork(msg *Msg) {
	fmt.Println("-----------------------------------")
	fmt.Println("Starting lookupNetwork")
	//if d is  responsible for id
	if between([]byte(d.id), []byte(d.successor.id), []byte(msg.Key)) {
		m := makeMsg("response", msg.Origin, d.id, d.Address(), msg.Time, d.Address())
		fmt.Println("Trying to Send 1 lookup")
		fmt.Println("")
		d.Transport.send(m, nil)
		fmt.Println("Sending 1 done")
		fmt.Println("")
		return
		//return d
	}
	//otherwise use fingers
	dist := distance(d.id, msg.Key, len(d.finger))
	index := dist.BitLen() - 1
	if index < 0 {
		m := makeMsg("response", msg.Origin, d.id, d.Address(), msg.Time, d.Address())
		fmt.Println("Trying to send 2 lookup")
		fmt.Println("")
		d.Transport.send(m, nil)
		fmt.Println("Sending 2 done")
		fmt.Println("")

		return
	}
	//fmt.Println("INDEX", index)

	fmt.Println("TEST1 LOOKUP!")
	//stegar ner tills fingret inte pekar på sig själv

	for ; index > 0 && d.finger[index].node == d; index-- {

	}
	fmt.Println(index)
	fmt.Println("TEST 2 LOOKUP")
	// Kollar så vi inte hamnar för långt
	diff := big.Int{}
	diff.Sub(dist, distance(d.id, d.finger[index].node.id, len(d.finger)))
	for index > 0 && diff.Sign() < 0 {
		index--
		diff.Sub(dist, distance(d.id, d.finger[index].node.id, len(d.finger)))
	}
	//kollar så vi inte pekar på oss själva
	if d.finger[index].node == d || diff.Sign() < 0 {
		fmt.Println("ERROR ERROR alles gebort auf the baut")
		// send message to the successor node to do a lookup
		m := makeMsg("lookupNetwork", d.successor.Address(), msg.Key, msg.Origin, msg.Time, d.Address())
		fmt.Println("Sending to lookup again 1")
		fmt.Println("")
		d.Transport.send(m, nil)
		fmt.Println("Sending to lookup again 1 done")
		fmt.Println("")
		return
		//return d.successor.lookup(hash)

	}

	// if nothing of the above works
	m := makeMsg("lookupNetwork", d.finger[index].node.Address(), msg.Key, msg.Origin, msg.Time, d.Address())
	fmt.Println("Sending to lookup again 2")
	fmt.Println("")
	d.Transport.send(m, nil)
	fmt.Println("sending to lookup again")
	fmt.Println("")
	return

	//return d.finger[index].node.lookup(hash)

}

// H2 har kastat bort hela update finger table
//update all nodes whose finger should refer to n
//func (n *DHTNode) update_others() {
//	for i := 1; i <= len(n.finger); i++ {
//		big_n := big.Int{}
//		sub_big_int := big.Int{}
//		result := big.Int{}
//
//		big_n.SetString(n.id, 16)
//		sub_big_int.Exp(big.NewInt(2), big.NewInt(int64(i-1)), nil)

//big_n.Sub(big_n, sub_big_int)
//bigString := big_n.String()
//		result.Sub(&big_n, &sub_big_int)
//		if result.Sign() < 0 {
//			fmt.Println("fixar negativa tal")
//will be used for 2^(nodes to be used)
//			big_totalnodes := big.Int{}
//the amount of nodes to be used
//big_nodes := big.Int{}
//used to do the calculation for sub
//			big_negative := result

//sets the nodes variable to a big int from the size of n.fingers

//			big_totalnodes.Exp(big.NewInt(2), big.NewInt(int64(len(n.finger))), nil)
//

//			fmt.Println("totalt antal noder: ", big_totalnodes)
//calculate result
//			fmt.Println("big_negative: ", big_negative)
//			result.Add(&big_totalnodes, &big_negative)

//			fmt.Println("här kommer det färdiga talet!: ")
/////HÄR MÅSTE DET CHECKAS SÅ ATT VI INTE TAR -2 när det ska vara node 7 t.ex
//		}
//		bigString := fmt.Sprintf("%x", result.Bytes())
//		fmt.Println(bigString)
//		fmt.Println()
//		fmt.Println()
//		p := n.lookup(bigString)
//		if p != n {
//			p.update_finger_table(n, i)
//		}

//	}
//
//}
//func (n *DHTNode) testCalcFingers(k int, m int) {
//	bigN := big.Int{}
//	bigN.SetString(n.id, 16)
//
//	fmt.Println(calcFinger(bigN.Bytes(), k, m))
//
//}

func (d *DHTNode) Address() string {
	return d.address + ":" + d.port

}

func (d *DHTNode) FingerPrint() {
	fmt.Println("Här är dina fingrar")

	for i := 0; i < 160; i++ {
		fmt.Println("finger nr:", i, " ", d.finger[i].start)

	}
}

func TimeNow() int64 {
	return time.Now().UnixNano()
}
