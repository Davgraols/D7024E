package main

import (
	"testing"
)

func testNewRandomKademliaID(t *testing.T) {
	testone := NewRandomKademliaID()
	testtwo := NewRandomKademliaID()

	if testone == testtwo {

		t.Errorf("error, generate two same ID. ")
	}
}
