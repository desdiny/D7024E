/*
fixa för nätverk:
	addtoring // micke
	lookup
	update_others
	kanske (update_fingertable)

kolla vad som behövs i msg samt transport structen


uppdatera fingrar automagist .... läs i handout filen // trolle


*/

package dht

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"strings"
	"testing"
	"time"
)

//////////////////////////////////////
//egen inlagt
//var antalfingrar int = 3

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
	Type, Key, Src, Dst, Origin,Time string
}

//struct for Transport from lab handout
/*

vi vill ta emot ett meddelande.
läsa ut meddelandet och sedan returnera svaret till source addressen

*/
type Transport struct {
	node        *DHTNode
	bindAddress string
	channel map[int64]chan Msg

	// chan,,,, mutexlås
}

// listen function from lab handout
func (transport *Transport) listen() {
	udpAddr, err := net.ResolveUDPAddr("udp", transport.bindAddress)
	conn, err := net.ListenUDP("udp", udpAddr)
	defer conn.Close()
	dec := json.NewDecoder(conn)
	for {
		msg := Msg{}
		err := dec.Decode(&msg)
		//Parse(msg)
		// if type is response check timestamp and call the channel
		//we got a message maby baby?
		//Parse vad det är för metod (lookup, addToring)

	}

}

// send function from lab handout
func (transport *Transport) send(msg *Msg, ch chan Msg) {
	if ch != nil{
		transport.channel[msg.Time] = ch
	}

	udpAddr, err := net.ResolveUDPAddr("udp", msg.Dst)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	defer conn.Close()
	_, err = conn.Write(msg.Bytes())
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

	//###################################//
   //									//
  // DHT NODER OCH DESS FUNKTIONER     //
 //								      //
//###################################//
type DHTNode struct {
	id, address, port      string
	successor, predecessor *DHTNode
	finger                 []*Fingers //links to Fingers struct
	Transport 				*Transport
}

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


func MakeDHTNode(idcheck *string, address string, port string) *DHTNode {
	n := new(DHTNode)
	if idcheck == nil {
		n.id = generateNodeId()
		n.address = address
		n.port = port
		n.successor = n
		n.predecessor = n
		n.finger = make([]*Fingers, 160) //change to use for 3 and 160
		n.Transport = makeTransport(n, n.address)
		n.Transport.listen()

	} else {
		n.id = *idcheck
		n.address = address
		n.port = port
		n.successor = n
		n.predecessor = n
		n.finger = make([]*Fingers, 160) //change to use for 3 and 160
		n.Transport = makeTransport(n, n.address)
		n.Transport.listen()
	}
	return n

}

func (n *DHTNode) initFingerTable(newnode *DHTNode) {
		if n.finger[0] == nil {
		// fixar fingrar special första gången
		for i := 1; i <= len(n.finger); i++ {
			fingerID, _ := calcFinger([]byte(n.id), i, len(n.finger))
			if len(fingerID) < len(n.id) {
				fingerID = strings.Repeat("0", len(n.id)-len(fingerID)) + fingerID
			}
			tempnode := n.lookup(fingerID)

			if tempnode.id != fingerID {
				tempnode = tempnode.successor

			}
			n.finger[i-1] = &Fingers{fingerID, tempnode}

			fmt.Println(n.finger[i-1].node.id)
		}

	}
	// fixar fingrar där  för att fylla på med nollor på rätt ställen etc
	for i := 1; i <= len(n.finger); i++ {
		fingerID, _ := calcFinger([]byte(newnode.id), i, len(n.finger))
		if len(fingerID) < len(n.id) {
			fingerID = strings.Repeat("0", len(n.id)-len(fingerID)) + fingerID
		}
		tempnode := n.lookup(fingerID)
		if tempnode.id != fingerID {
			tempnode = tempnode.successor

		}
		newnode.finger[i-1] = &Fingers{fingerID, tempnode}
		fmt.Println(newnode.finger[i-1].node.id)

	}
	return
	
}
func makeMsg(Type string, Dst string, Key string, Origin string) *Msg{
	m := new(Msg)
	m.Type = Type
	m.Dst = Dst
	m.Key = Key
	m.Origin = Origin
	m.Time = time.Now().UnixNano()
	return m

}


func makeTransport(node *DHTNode, bindAddress string) *Transport {
	s := new(Transport)
	s.node = node
	s.bindAddress = bindAddress
	s.channel = make(map[int64]chan Msg)
	return s
}

/////////////////////////////////////////////////
	////////////////////////////////////////////
   // new func for addToRing for networking  //
  ////////////////////////////////////////////
/////////////////////////////////////////////


func (n *DHTNode) joinRing(networkaddr string) {
	channel := make (chan Msg)
	fmt.Println("calling node on address: ", networkaddr)
	m := makeMsg("lookup", networkaddr, n.id, n.address)
	n.Transport.send(m, channel)

	req := <- channel
	joinidandaddress := n.id +","+ n.address
	m = makeMsg("join", req.Src, joinidandaddress, n.address)
	n.Transport.send(m, channel)

	// waiting for answer
	req = <-channel

	// split req (id and address)
	a := strings.Split(req, ",")
	// create a new node
	s:= new(DHTNode)
	s.id = a[0]
	s.address = a[1]

	n.predecessor = s

	//inte än fixat
	n.initFingerTable(newnode)


//contacts node in ring
//	node := n.lookup(newnode.id)

//	oldnode := node.successor
//	node.successor = newnode
//	newnode.successor = oldnode
//	newnode.predecessor = node
//	oldnode.predecessor = newnode
	n.update_others()
}
// the node that jumps on the node
func (n *DHTNode) join(msg *Msg) {
	channel := make (chan Msg)
	// splits the incomming keys
	a := strings.Split(msg.Key, ",")

	fmt.Println("the joining has begun, calling to set predecessor on next node")
	joinidandaddress := a[0] + "," + a[1]
	m := makeMsg("changePredecessor", n.successor.address, joinidandaddress, n.address)
	n.Transport.send(m, channel)

	//creates a new node 
	s:= new(DHTNode)
	s.id = a[0]
	s.address = a[1]

	// adds the new node as the nodes succsessor
	n.successor = s

	//adding both to one variable so we can send it in the key value
	// have to concatinate when message is recived
	joinidandaddress = n.id + "," + n.address

	//creates message
	m = makeMsg("joinRing", newnode.address, joinidandaddress, n.address)

	// sends message
	n.Transport.send(m, channel)

}

func (n *DHTNode) changePredecessor(msg *Msg) {

	//split incomming key
	a := strings.Split(msg.Key, ",")

	//create a new node on this instance
	s:= new(DHTNode)
	s.id = a[0]
	s.address = a[1]

	// adds the node to n's predecessor
	n.predecessor = s
	
}




func (n *DHTNode) addToRing(newnode *DHTNode) {
	fmt.Println("Nodens id: ", newnode.id)
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

func (d *DHTNode) lookup(hash string) *DHTNode {

	if between([]byte(d.id), []byte(d.successor.id), []byte(hash)) {
		// returns that this node should be responible for this
		// how to use type in this case?
		// can we just send
		makeMsg(, Dst, Key, Origin)
		return d
	}

	dist := distance(d.id, hash, len(d.finger))
	index := dist.BitLen() - 1
	if index < 0 {
		return d
	}
	fmt.Println("INDEX", index)

	//stegar ner tills fingret inte pekar på sig själv
	for ; index > 0 && d.finger[index].node == d; index-- {

	}
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
		return d.successor.lookup(hash)

	}
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
	return d.finger[index].node.lookup(hash)
	/*
	   här under har vi den förra funktionen för att köra utan fingrar
	*/
	//	return d.successor.lookup(hash)
}

//om s är i (någon av) n  fingrar, uppdatera n's fingrar med s
func (n *DHTNode) update_finger_table(s *DHTNode, i int) {
	fmt.Println("updating finger", i, "on", n.id)
	if s.successor == n.finger[i-1].node {
		n.finger[i-1].node = s
		p := n.predecessor
		if p != n {
			p.update_finger_table(s, i)
		}

	}

}


// H2 har kastat bort hela update finger table
//update all nodes whose finger should refer to n
func (n *DHTNode) update_others() {
	for i := 1; i <= len(n.finger); i++ {
		big_n := big.Int{}
		sub_big_int := big.Int{}
		result := big.Int{}

		big_n.SetString(n.id, 16)
		sub_big_int.Exp(big.NewInt(2), big.NewInt(int64(i-1)), nil)

		//big_n.Sub(big_n, sub_big_int)
		//bigString := big_n.String()
		result.Sub(&big_n, &sub_big_int)
		if result.Sign() < 0 {
			fmt.Println("fixar negativa tal")
			//will be used for 2^(nodes to be used)
			big_totalnodes := big.Int{}
			//the amount of nodes to be used
			//big_nodes := big.Int{}
			//used to do the calculation for sub
			big_negative := result

			//sets the nodes variable to a big int from the size of n.fingers

			big_totalnodes.Exp(big.NewInt(2), big.NewInt(int64(len(n.finger))), nil)
			//

			fmt.Println("totalt antal noder: ", big_totalnodes)
			//calculate result
			fmt.Println("big_negative: ", big_negative)
			result.Add(&big_totalnodes, &big_negative)

			fmt.Println("här kommer det färdiga talet!: ")
			/////HÄR MÅSTE DET CHECKAS SÅ ATT VI INTE TAR -2 när det ska vara node 7 t.ex
		}
		bigString := fmt.Sprintf("%x", result.Bytes())
		fmt.Println(bigString)
		fmt.Println()
		fmt.Println()
		p := n.lookup(bigString)
		if p != n {
			p.update_finger_table(n, i)
		}

	}

}
func (n *DHTNode) testCalcFingers(k int, m int) {
	bigN := big.Int{}
	bigN.SetString(n.id, 16)

	fmt.Println(calcFinger(bigN.Bytes(), k, m))

}

////////////////////////////////////
//slut på egenediterat
///////////////////////////////////

// test cases can be run by calling e.g. go test -test.run TestRingSetup
// go run test will run all tests

/*
func TestRingSetup(t *testing.T) {
	// note nil arg means automatically generate ID, e.g. f38f3b2dcc69a2093f258e31902e40ad33148385
	node1 := makeDHTNode(nil, "localhost", "1111")
	node2 := makeDHTNode(nil, "localhost", "1112")
	node3 := makeDHTNode(nil, "localhost", "1113")
	node4 := makeDHTNode(nil, "localhost", "1114")
	node5 := makeDHTNode(nil, "localhost", "1115")
	node6 := makeDHTNode(nil, "localhost", "1116")
	node7 := makeDHTNode(nil, "localhost", "1117")
	node8 := makeDHTNode(nil, "localhost", "1118")
	node9 := makeDHTNode(nil, "localhost", "1119")
*/
/*
	node1.addToRing(node2)
	node1.addToRing(node3)
	node1.addToRing(node4)
	node4.addToRing(node5)
	node3.addToRing(node6)
	node3.addToRing(node7)
	node3.addToRing(node8)
	node7.addToRing(node9)
*/
/*
	node1.addToRing(node2)
	node2.addToRing(node3)
	node3.addToRing(node4)
	node4.addToRing(node5)
	node5.addToRing(node6)
	node6.addToRing(node7)
	node7.addToRing(node8)
	node8.addToRing(node9)
	//node9.addToRing(node1)

	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Println("RING STRUCTURE")
	//fmt.Println(node2.successor)
	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Println("Print Ring node 1 :")
	node2.printRing()
	fmt.Println("------------------------------------------------------------------------------------------------")
}
*/
/*
 * Example of expected output.
 *
 * str=hello students!
 * hashKey=cba8c6e5f208b9c72ebee924d20f04a081a1b0aa
 * c588f83243aeb49288d3fcdeb6cc9e68f9134dce is respoinsible for cba8c6e5f208b9c72ebee924d20f04a081a1b0aa
 * c588f83243aeb49288d3fcdeb6cc9e68f9134dce is respoinsible for cba8c6e5f208b9c72ebee924d20f04a081a1b0aa
 */
/*
func TestLookup(t *testing.T) {
	node1 := makeDHTNode(nil, "localhost", "1111")
	node2 := makeDHTNode(nil, "localhost", "1112")
	node3 := makeDHTNode(nil, "localhost", "1113")
	node4 := makeDHTNode(nil, "localhost", "1114")
	node5 := makeDHTNode(nil, "localhost", "1115")
	node6 := makeDHTNode(nil, "localhost", "1116")
	node7 := makeDHTNode(nil, "localhost", "1117")
	node8 := makeDHTNode(nil, "localhost", "1118")
	node9 := makeDHTNode(nil, "localhost", "1119")

	node1.addToRing(node2)
	node1.addToRing(node3)
	node1.addToRing(node4)
	node4.addToRing(node5)
	node3.addToRing(node6)
	node3.addToRing(node7)
	node3.addToRing(node8)
	node7.addToRing(node9)

	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Println("RING STRUCTURE")
	fmt.Println("------------------------------------------------------------------------------------------------")
	node1.printRing()
	fmt.Println("------------------------------------------------------------------------------------------------")

	str := "hello students!"
	hashKey := sha1hash(str)
	fmt.Println("str = " + str)
	fmt.Println("hashKey = " + hashKey)

	fmt.Println("node 1: " + node1.lookup(hashKey).id + " is respoinsible for " + hashKey)
	fmt.Println("node 5: " + node5.lookup(hashKey).id + " is respoinsible for " + hashKey)

	fmt.Println("------------------------------------------------------------------------------------------------")

}
*/
/*
 * Example of expected output.
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            0
 * k            1
 * m            3
 * 2^(k-1)      1
 * (n+2^(k-1))  1
 * 2^m          8
 * result       1
 * result (hex) 01
 * successor    01
 * distance     1
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            0
 * k            2
 * m            3
 * 2^(k-1)      2
 * (n+2^(k-1))  2
 * 2^m          8
 * result       2
 * result (hex) 02
 * successor    02
 * distance     2
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            0
 * k            3
 * m            3
 * 2^(k-1)      4
 * (n+2^(k-1))  4
 * 2^m          8
 * result       4
 * result (hex) 04
 * successor    04
 * distance     4
 */

func TestFinger3bits(t *testing.T) {
	id0 := "00"
	id1 := "01"
	id2 := "02"
	id3 := "03"
	id4 := "04"
	id5 := "05"
	id6 := "06"
	id7 := "07"

	node0 := makeDHTNode(&id0, "localhost", "1111")
	node1 := makeDHTNode(&id1, "localhost", "1112")
	node2 := makeDHTNode(&id2, "localhost", "1113")
	node3 := makeDHTNode(&id3, "localhost", "1114")
	node4 := makeDHTNode(&id4, "localhost", "1115")
	node5 := makeDHTNode(&id5, "localhost", "1116")
	node6 := makeDHTNode(&id6, "localhost", "1117")
	node7 := makeDHTNode(&id7, "localhost", "1118")

	node0.addToRing(node1)
	node1.addToRing(node2)
	node1.addToRing(node3)
	node1.addToRing(node4)
	node4.addToRing(node5)
	node3.addToRing(node6)
	node3.addToRing(node7)

	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Println("RING STRUCTURE")
	fmt.Println("------------------------------------------------------------------------------------------------")
	node0.printRing()
	fmt.Println("------------------------------------------------------------------------------------------------")

	//	node0.testCalcFingers(1, 3)
	fmt.Println("")
	//	node0.testCalcFingers(2, 3)
	fmt.Println("")
	//	node0.testCalcFingers(3, 3)
}

/*
 * Example of expected output.
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            682874255151879437996522856919401519827635625586
 * k            0
 * m            160
 * 2^(k-1)      1
 * (n+2^(k-1))  682874255151879437996522856919401519827635625587
 * 2^m          1461501637330902918203684832716283019655932542976
 * finger       682874255151879437996522856919401519827635625587
 * finger (hex) 779d240121ed6d5e8bd0cb6529b08e5c617b5e73
 * successor    779d240121ed6d5e8bd0cb6529b08e5c617b5e72
 * distance     0

 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            682874255151879437996522856919401519827635625586
 * k            1
 * m            160
 * 2^(k-1)      1
 * (n+2^(k-1))  682874255151879437996522856919401519827635625587
 * 2^m          1461501637330902918203684832716283019655932542976
 * finger       682874255151879437996522856919401519827635625587
 * finger (hex) 779d240121ed6d5e8bd0cb6529b08e5c617b5e73
 * successor    779d240121ed6d5e8bd0cb6529b08e5c617b5e72
 * distance     0
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            682874255151879437996522856919401519827635625586
 * k            80
 * m            160
 * 2^(k-1)      604462909807314587353088
 * (n+2^(k-1))  682874255151879437996523461382311327142222978674
 * 2^m          1461501637330902918203684832716283019655932542976
 * finger       682874255151879437996523461382311327142222978674
 * finger (hex) 779d240121ed6d5e8bd14b6529b08e5c617b5e72
 * successor    779d240121ed6d5e8bd0cb6529b08e5c617b5e72
 * distance     0
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            682874255151879437996522856919401519827635625586
 * k            120
 * m            90
 * 2^(k-1)      664613997892457936451903530140172288
 * (n+2^(k-1))  682874255152544051994415314855853423357775797874
 * 2^m          1237940039285380274899124224
 * finger       1180872106465109536036052594
 * finger (hex) 03d0cb6529b08e5c617b5e72
 * successor    f880fb198b7059ae92a69968727d84da9c94dd15
 * distance     877444087302148207702277795
 *
 * calulcating result = (n+2^(k-1)) mod (2^m)
 * n            682874255151879437996522856919401519827635625586
 * k            160
 * m            160
 * 2^(k-1)      730750818665451459101842416358141509827966271488
 * (n+2^(k-1))  1413625073817330897098365273277543029655601897074
 * 2^m          1461501637330902918203684832716283019655932542976
 * finger       1413625073817330897098365273277543029655601897074
 * finger (hex) f79d240121ed6d5e8bd0cb6529b08e5c617b5e72
 * successor    d0a43af3a433353909e09739b964e64c107e5e92
 * distance     508258282811496687056817668076520806659544776736
 */

/*func TestFinger160bits(t *testing.T) {
	// note nil arg means automatically generate ID, e.g. f38f3b2dcc69a2093f258e31902e40ad33148385
	node1 := makeDHTNode(nil, "localhost", "1111")
	node2 := makeDHTNode(nil, "localhost", "1112")
	node3 := makeDHTNode(nil, "localhost", "1113")
	node4 := makeDHTNode(nil, "localhost", "1114")
	node5 := makeDHTNode(nil, "localhost", "1115")
	node6 := makeDHTNode(nil, "localhost", "1116")
	node7 := makeDHTNode(nil, "localhost", "1117")
	node8 := makeDHTNode(nil, "localhost", "1118")
	node9 := makeDHTNode(nil, "localhost", "1119")

	node1.addToRing(node2)
	node1.addToRing(node3)
	node1.addToRing(node4)
	node4.addToRing(node5)
	node3.addToRing(node6)
	node3.addToRing(node7)
	node3.addToRing(node8)
	node7.addToRing(node9)

	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Println("RING STRUCTURE")
	fmt.Println("------------------------------------------------------------------------------------------------")
	node1.printRing()
	fmt.Println("------------------------------------------------------------------------------------------------")

	node3.testCalcFingers(0, 160)
	fmt.Println("")
	node3.testCalcFingers(1, 160)
	fmt.Println("")
	node3.testCalcFingers(80, 160)
	fmt.Println("")
	node3.testCalcFingers(120, 90)
	fmt.Println("")
	node3.testCalcFingers(160, 160)
	fmt.Println("")
}
*/
