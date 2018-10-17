package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

// the static number of bytes in a KademliaID
const IDLength = 20

// type definition of a KademliaID
type KademliaID [IDLength]byte

// NewKademliaID returns a new instance of a KademliaID based on the string input
func NewKademliaID(data string) *KademliaID {
	decoded, _ := hex.DecodeString(data)

	newKademliaID := KademliaID{}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = decoded[i]
	}
	fmt.Println("newID = ")
	fmt.Println(&newKademliaID)
	return &newKademliaID
}

//func testNewKadid(t *testing.T) {
//	testid := NewKademliaID("1337")
//}

// NewRandomKademliaID returns a new instance of a random KademliaID,
// change this to a better version if you like
func NewRandomKademliaID() *KademliaID {
	newKademliaID := KademliaID{}
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = uint8(rand.Intn(256))
	}
	return &newKademliaID
}

func NewRandomSerial() int32 {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Int31()
}

func NewRandomHash(data string) *KademliaID {
	letterBytes := data
	newKademliaID := KademliaID{}
	b := make([]byte, 20)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	for i := 0; i < IDLength; i++ {
		newKademliaID[i] = b[i]
	}
	return &newKademliaID
	//fmt.Println(string(&newKademliaID))
	//fmt.Println(len(string(b)))

}

// Less returns true if kademliaID < otherKademliaID (bitwise)
func (kademliaID KademliaID) Less(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return kademliaID[i] < otherKademliaID[i]
		}
	}
	return false
}

// Equals returns true if kademliaID == otherKademliaID (bitwise)
func (kademliaID KademliaID) Equals(otherKademliaID *KademliaID) bool {
	for i := 0; i < IDLength; i++ {
		if kademliaID[i] != otherKademliaID[i] {
			return false
		}
	}
	return true
}

// CalcDistance returns a new instance of a KademliaID that is built
// through a bitwise XOR operation betweeen kademliaID and target
func (kademliaID KademliaID) CalcDistance(target *KademliaID) *KademliaID {
	result := KademliaID{}
	for i := 0; i < IDLength; i++ {
		result[i] = kademliaID[i] ^ target[i]
	}
	return &result
}

// String returns a simple string representation of a KademliaID
func (kademliaID *KademliaID) String() string {
	return hex.EncodeToString(kademliaID[0:IDLength])
}

func (kademliaID *KademliaID) ToBytes() []byte {
	kadID := kademliaID[0:IDLength]
	return kadID
}

func IdFromBytes(idInBytes []byte) *KademliaID {

	kadId := KademliaID{}
	for i, abyte := range idInBytes {
		kadId[i] = abyte
	}
	return &kadId
}
