package dht

import (
	"github.com/boltdb/bolt"
	//"encoding/json" // used for networking
	"fmt"
	"math/big"  // used for fingers
	"math/rand" //used for updating fingers
	//"net"
	"log"
	"strings"
	//"testing"
	"time" // used to update fingers and to set time for msg
)

//###################################//
//									 //
// DHT NODER OCH DESS FUNKTIONER     //
//								     //
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

func (node *DHTNode) AutoFingers() {
	channel := make(chan Msg)
	i := rand.Intn(3) //vet inte ifall det behövs en random var i intn(???)
	//fmt.Println("Autouppdaterar finger: ", i)
	//create autofingers message
	m := makeMsg("lookupNetwork", node.Address(), node.finger[i].start, node.Address(), TimeNow(), node.Address())
	node.Transport.send(m, channel)
	//waitning for answer
	req := <-channel
	//split address and port
	a := strings.Split(req.Src, ":")
	if req.Src != node.finger[i].node.Address() && req.Key != node.finger[i].node.id {
		fmt.Println("Autouppdaterar finger: ", i)
		node.finger[i].node.id = req.Key
		node.finger[i].node.address = a[0]
		node.finger[i].node.port = a[1]
	}

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
	m := makeMsg("lookupNetwork", networkaddr, n.id, n.Address(), TimeNow(), n.Address())
	n.Transport.send(m, channel)
	fmt.Println("node has been called")

	req := <-channel
	fmt.Println("receved a answer on JoinRing with KEY: ", req.Key)
	joinidandaddress := n.id + "," + n.address + "," + n.port
	m = makeMsg("join", req.Src, joinidandaddress, n.Address(), TimeNow(), n.Address())
	n.Transport.send(m, channel)
	fmt.Println("Sending message to join function....")

	// waiting for answer
	req = <-channel
	fmt.Println("recived ansver from Joinfunction")

	// split req (id and address)
	a := strings.Split(req.Key, ",")

	//fmt.Println("")
	s := new(DHTNode)
	s.id = a[0]
	s.address = a[1]
	s.port = a[2]
	//	n.predecessor.id = a[0]
	//	n.predecessor.address = a[1]
	//	n.predecessor.port = a[2]
	// create a new node
	//s := new(DHTNode)
	//b := strings.Split(req.Src, ":")
	//s := new(OutsideNode)
	//s.id = a[0]
	//s.address = a[1]
	//s.port = b[1]
	fmt.Println("added the predecessor with id: ", n.predecessor.id, "address: ", n.predecessor.address, "port: ", n.predecessor.port)
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
	fmt.Println("")
	fmt.Println("-----------------------------------------")
	fmt.Println("Join operation initializing")
	fmt.Println("")
	//channel := make(chan Msg)
	// splits the incomming keys
	a := strings.Split(msg.Key, ",")

	fmt.Println("the joining has begun, calling to set predecessor on next node")
	fmt.Println("")
	oldsuccessor := n.successor
	//joinidandaddress := a[0] + "," + a[1] + "," + a[2]
	m := makeMsg("changePredecessor", n.successor.Address(), msg.Key, n.Address(), TimeNow(), n.Address())
	n.Transport.send(m, nil)
	fmt.Println("Sending changePredecessor to successor: ", n.successor.Address())
	fmt.Println("")
	//creates a new node
	s := new(DHTNode)
	//v := strings.Split(msg.Origin, ":")
	s.id = a[0]
	s.address = a[1]
	s.port = a[2]
	//n.successor.id = a[0]
	//n.successor.address = a[1]
	//n.successor.port = a[2]
	n.successor = s
	fmt.Println("added the successor: ", " id: ", n.successor.id, " address: ", n.successor.address, " port: ", n.successor.port)
	fmt.Println("")
	// adds the new node as the nodes succsessor
	//	n.successor = s
	n.finger[0].node = n.successor

	//adding both to one variable so we can send it in the key value
	// have to concatinate when message is recived
	joinidandaddresstest := n.id + "," + n.address + "," + n.port

	//creates message
	fmt.Println("<-----message content----->")
	fmt.Println("msg.origin: ", msg.Origin)
	fmt.Println("joinidandaddress: ", joinidandaddresstest)
	fmt.Println("n.Address(): ", n.Address())
	fmt.Println("<-----message content----->")

	o := makeMsg("response", msg.Origin, joinidandaddresstest, n.Address(), msg.Time, n.Address())

	// sends message
	n.Transport.send(o, nil)

	key := oldsuccessor.id + ":" + oldsuccessor.address + ":" + oldsuccessor.port
	m = makeMsg("changeSuccessor", msg.Origin, key, n.Address(), TimeNow(), n.Address())
	fmt.Println("Sending message to change Successor to: ", msg.Origin)
	n.Transport.send(m, nil)

	fmt.Println("sends respons to JoinRing to let it add with our cred")
	fmt.Println("")
	fmt.Println("Join operation complete")
	fmt.Println("------------------------------------------")
	fmt.Println("")
}

func (n *DHTNode) changePredecessor(msg *Msg) {
	fmt.Println("")
	fmt.Println("------------------------------------------")
	fmt.Println("Starting changePredecessor")
	//split incomming key
	a := strings.Split(msg.Key, ",")

	//create a new node on this instance
	s := new(DHTNode)
	s.id = a[0]
	s.address = a[1]
	s.port = a[2]
	//n.predecessor.id = a[0]
	//n.predecessor.address = a[1]
	//n.predecessor.port = a[2]

	// adds the node to n's predecessor
	n.predecessor = s
	fmt.Println("Have changed the predecessor to: ", " id: ", s.id, " address: ", s.address, " port: ", s.port)

	//m := makeMsg("changeSuccessor", s.Address(), n.id, n.Address(), msg.Time, n.Address())

	// sends message
	//	n.Transport.send(m, nil)
	fmt.Println("Changing Predecessor completed")
	fmt.Println("------------------------------------------")
	fmt.Println("")

}

func (n *DHTNode) changeSuccessor(msg *Msg) {
	fmt.Println("")
	fmt.Println("------------------------------------------")
	fmt.Println("Starting changeSuccessor")
	a := strings.Split(msg.Key, ":")
	//	n.successor.id = msg.Key
	//n.successor.address = a[1]
	//	n.successor.port = a[0]
	s := new(DHTNode)
	s.id = msg.Key
	//a := strings.Split(msg.Key, ":")
	s.id = a[0]
	s.address = a[1]
	s.port = a[2]
	n.successor = s
	fmt.Println("Have changed the successor to: ", " id: ", s.id, " address: ", s.address, " port: ", s.port)
	n.finger[0].node = n.successor
	fmt.Println("Changing Successor completed")
	fmt.Println("------------------------------------------")
	fmt.Println("")
}

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

	fmt.Println("Check between vaiables, d.id:", d.id, " d.successor.id: ", d.successor.id, " msg.Key: ", msg.Key)
	fmt.Println("This is our message: msg.Key", msg.Key, "msg.Origin: ", msg.Origin)
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

func (d *DHTNode) Address() string {
	return d.address + ":" + d.port

}

func (d *DHTNode) FingerPrint() {
	fmt.Println("Här är dina fingrar")

	for i := 0; i < 160; i++ {
		fmt.Println("finger nr:", i, " ", d.finger[i].start)

	}
}
func (d *DHTNode) IdPrint() {

	fmt.Println("Detta är ditt id:", d.id)

}

func (d *DHTNode) PreID() {

	fmt.Println("Detta är din predecessor, id: ", d.predecessor.id, "address: ", d.predecessor.address, "port: ", d.predecessor.port)

}

func (d *DHTNode) SucID() {

	fmt.Println("Detta är din successor, id: ", d.successor.id, "address: ", d.successor.address, "port: ", d.successor.port)

}

func (n *DHTNode) Ping() {
	fmt.Println("test")

	channel := make(chan Msg)

	m := makeMsg("Pong", n.successor.Address(),
		"ALLO", n.Address(), TimeNow(), n.Address())

	//fmt.Println(m.Type)
	n.Transport.send(m, channel)
	fmt.Println(n.Address())

	select {
	case req := <-channel:
		fmt.Println("Successor is responding", req)
	case <-time.After(2 * time.Second):
		fmt.Println("Successor is not responding")

	}

}
func (n *DHTNode) Pong(msg *Msg) {
	if msg.Key == "ALLO" {
		fmt.Println("Har fått pong, skickar ping")
		m := makeMsg("response", msg.Origin, "Ping", n.Address(), msg.Time, n.Address())
		fmt.Println("Pong value", m)
		n.Transport.send(m, nil)
	} else {
		fmt.Println("har inte fått ping :/")
	}
}

func TimeNow() int64 {
	return time.Now().UnixNano()
}

// initializing DB
func (n *DHTNode) initDB() {

	db, err := bolt.Open(n.id+".db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}

/*										*
				AddData

	Used to lookup where to store data
	then contacts that node and asks it
	to write the data

*										*/
func (n *DHTNode) AddData(key string, value string) {
	fmt.Println("")
	fmt.Println("--------------------------------------")
	fmt.Println("Starting AddData with key: ", key, " and value: ", value)
	fmt.Println("")
	channel := make(chan Msg)
	//hashKey := sha1hash(key)
	m := makeMsg("lookupNetwork", n.Address(), key, n.Address(), TimeNow(), n.Address())
	n.Transport.send(m, channel)
	fmt.Println("Lookup has been sent the key")
	fmt.Println("")

	req := <-channel
	fmt.Println("Recived answer that the value should be placed on node: ", req.Key, " with address: ", req.Src)
	fmt.Println("")
	data := hashKey + ":" + value
	m = makeMsg("writeData", req.Src, data, n.Address(), TimeNow(), n.Address())
	n.Transport.send(m, nil)

}

/*										*
				writeData

	Used to write data to the nodes
	database.
	It then calls the predecessor
	and commands it to replicate the
	input.

*										*/
func (n *DHTNode) writeData(msg *Msg) {
	//channel := make(chan Msg)
	//a := strings.Split(msg.Key, ":")
	key := a[0]
	value := a[1]
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(n.id))
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(key), []byte(value))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	m := makeMsg("replicateData", n.predecessor.Address(), msg.Key, n.Address(), TimeNow(), n.Address())
	n.Transport.send(m, nil)
	//
	// Send data to replication (predecessor)
	//

}

/*										*
				readData

	Lookup where the req data is saved
	and then contacts that node to ask
	it to send back the req data.
	When the node have sent its data
	we return it as a string.

*										*/
func (n *DHTNode) readData(key string) string {
	channel := make(chan Msg)
	hashKey := sha1hash(key)
	m := makeMsg("lookupNetwork", n.Address(), hashKey, n.Address(), TimeNow(), n.Address())
	n.Transport.send(m, channel)

	req := <-channel

	m = makeMsg("returnData", req.Src, hashKey, n.Address(), TimeNow(), n.Address())
	n.Transport.send(m, channel)

	req = <-channel

	return req.Key

}

/*										*
				returnData

	Retrives the data stored in key
	and then
	returns the data requested from
	readData.
*										*/
func (n *DHTNode) returnData(msg *Msg) {
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(n.id))
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", []byte(n.id))
		}

		value := bucket.Get([]byte(msg.Key))

		m := makeMsg("response", msg.Src, string(value), n.Address(), msg.Time, n.Address())
		n.Transport.send(m, nil)

		return nil
	})

}

/*										*
				deleteData

	Searches where the data that is
	going to be deleted is stored.
	Once that is found we send a
	removeData tho that node.
*										*/
func (n *DHTNode) deleteData() {

}

/*										*
				removeData

	Removes the data of the key sent
	from deleteData.
	When it have been deleted it schoud
	conntact the replicated data and
	remove that key/value to
*										*/
func (n *DHTNode) removeData(msg *Msg) {

}

/*										*
				replicateData

	Replicates a newley added key/value
	from its successor.
*										*/
func (n *DHTNode) replicateData(msg *Msg) {
	channel := make(chan Msg)
	a := strings.Split(msg.Key, ":")
	key := a[0]
	value := a[1]
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(n.id))
		if err != nil {
			return err
		}

		err = bucket.Put([]byte(key), []byte(value))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

/*										*
				lookupData


*										*/
func (n *DHTNode) lookupData() {
	//channel := make(chan Msg)

	//
	//check if data is between successor and successorsuccessor
	//if so send data to successor
	//

}

//to add something bucket wants []byte for key and a []byte for value
