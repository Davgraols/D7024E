package main

import (
	"testing"
)

func testNewRandomKademliaID(t *testing.T) {
	testone := NewRandomKademliaID()
	testtwo := NewRandomKademliaID()

	if testone == testtwo {

		t.Errorf("error, generate two of the same ID. ")
	}
}

func testNewRandomSerial(t *testing.T) {
	testone := NewRandomSerial()
	testtwo := NewRandomSerial()

	if testone == testtwo {

		t.Errorf("error, generate two of the same Serial. ")
	}
}

func testNewRandomHash(t *testing.T) {
	testone := NewRandomHash("1337")
	testtwo := NewRandomHash("1337")

	if testone == testtwo {

		t.Errorf("error, generate two of the same hash. ")
	}
}

func testEquals(t *testing.T) {
	node := NewRandomKademliaID()
	if !(node.Equals(node)) {
		t.Errorf("the Equals du not return True if the nodes are equal")
	}
}

func testString(t *testing.T) {
	node := NewRandomKademliaID()
	nodestring := node.String()
	var i interface{} = nodestring
	_, ok := i.(string)
	if ok == false {
		t.Errorf("the String function dose not konvert kadmliaID to string")
	}
}

func testToBytes(t *testing.T) {
	node := NewRandomKademliaID()
	nodebyte := node.ToBytes()
	var i interface{} = nodebyte
	_, ok := i.([]byte)
	if ok == false {
		t.Errorf("the ToBytes function dose not konvert kadmliaID to byte")
	}

}

func testIdFromBytes(t *testing.T) {
	node := NewRandomKademliaID()
	nodebyte := node.ToBytes()
	nodes := IdFromBytes(nodebyte)
	var i interface{} = nodes
	_, ok := i.(*KademliaID)
	if ok == false {
		t.Errorf("the ToBytes function dose not konvert kadmliaID to byte")
	}
}
