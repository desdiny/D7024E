package dht

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

/////////////////////////////////////
//egen inlagt
//var antalfingrar int = 3

type DHTNode struct {
	id, address, port      string
	successor, predecessor *DHTNode
	finger                 []*DHTNode
}

//var antalfingrar int = 3

func makeDHTNode(idcheck *string, address string, port string) *DHTNode {
	n := new(DHTNode)
	if idcheck == nil {
		n.id = generateNodeId()
		n.address = address
		n.port = port
		n.successor = n
		n.predecessor = n
		n.finger = make([]*DHTNode, 3)

	} else {
		n.id = *idcheck
		n.address = address
		n.port = port
		n.successor = n
		n.predecessor = n
		n.finger = make([]*DHTNode, 3)
	}
	return n

}

func (n *DHTNode) addToRing(newnode *DHTNode) {
	for i := 0; i < len(n.finger); i++ {
		fingerID, _ := calcFinger([]byte(n.id), i, len(n.finger))
		n.finger[i] = n.lookup(fingerID)
	}
	node := n.lookup(newnode.id)
	oldnode := node.successor
	node.successor = newnode
	newnode.successor = oldnode
	newnode.predecessor = node
	oldnode.predecessor = newnode

	fmt.Println(n)
}

func (n *DHTNode) printRing() {

	nextNode := n.successor
	fmt.Println(n.id)
	for nextNode != n {
		fmt.Println(nextNode.id)
		nextNode = nextNode.successor

	}
}

func (d *DHTNode) tostring() (out string) {
	out = "DHTNode{id: " + d.id + ", address: " + d.address + ", port: " + d.port + "}"

	return
}

func (d *DHTNode) lookup(hash string) *DHTNode {

	if between([]byte(d.id), []byte(d.successor.id), []byte(hash)) {
		return d
	}
	return d.successor.lookup(hash)
}

//om s är i (någon av) n  fingrar, uppdatera n's fingrar med s
func (n *DHTNode) update_finger_table(s *DHTNode, i int) {
	if s.successor == n.finger[i-1] {
		n.finger[i-1] = s
		p := n.predecessor
		if p != n {
			p.update_finger_table(s, i)
		}

	}

}

//update all nodes whose finger should refer to n
func (n *DHTNode) update_others() {
	for i := 1; i <= len(n.finger); i++ {
		big_n := big.Int{}
		sub_big_int := big.Int{}
		result := big.Int{}

		big_n.SetString(n.id, 16)
		sub_big_int.SetInt64(int64(math.Exp2(float64(i - 1))))

		//big_n.Sub(big_n, sub_big_int)
		//bigString := big_n.String()
		result.Sub(&big_n, &sub_big_int)
		bigString := result.String()

		p := n.lookup(bigString)
		p.update_finger_table(n, i)

	}

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
/*func TestFinger3bits(t *testing.T) {
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
	node1.printRing()
	fmt.Println("------------------------------------------------------------------------------------------------")

	node3.testCalcFingers(1, 3)
	fmt.Println("")
	node3.testCalcFingers(2, 3)
	fmt.Println("")
	node3.testCalcFingers(3, 3)
}
*/

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
