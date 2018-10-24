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
