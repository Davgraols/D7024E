package main

import (
	"fmt"
	"testing"
)

func testNewRandomKademliaID(t *testing.T) {
	testone := NewRandomKademliaID()
	testtwo := NewRandomKademliaID()

	if testone == testtwo {

		t.Errorf("error, generate two same ID. ")
	}
	fmt.Println("nNewRandomKademliaID: pass")
}

func testNewRandomSerial(t *testing.T) {
	testone := NewRandomSerial()
	testtwo := NewRandomSerial()

	if testone == testtwo {

		t.Errorf("error, generate two same ID. ")
	}
	fmt.Println("NewRandomSerial: pass")
}

func testNewRandomHash(t *testing.T) {
	testone := NewRandomHash("1337")
	testtwo := NewRandomHash("1337")

	if testone == testtwo {

		t.Errorf("error, generate two same ID. ")
	}
	fmt.Println("NewRandomSerial: pass")
}
