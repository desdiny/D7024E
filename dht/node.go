
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