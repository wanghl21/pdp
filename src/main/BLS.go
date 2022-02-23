package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/Nik-U/pbc"
	"github.com/W1llyu/ourjson"
)

// messageData represents a signed message sent over the network
type messageData struct {
	message   string
	signature []byte
}

// This example computes and verifies a Boneh-Lynn-Shacham signature in a simulated conversation between Alice and Bob.
func main1() {
	// The authority generates system parameters
	// In a real application, generate this once and publish it
	params := pbc.GenerateA(160, 512)
	//fmt.Println(params)

	pairing := params.NewPairing() // instantiates a  pairing
	g := pairing.NewG2().Rand()
	// The authority distributes params and g to Alice and Bob
	sharedParams := params.String()

	sharedG := g.Bytes()

	/*
		//pairing := params.NewPairing()
		sharedParams := Utils{}.readFile("../config/b.properties")
		jsonStringPub := Utils{}.readFile("../config/pubkey")

		pubkeyObject, _ := ourjson.ParseObject(jsonStringPub)
		g1, _ := pubkeyObject.GetString("g")
		fmt.Printf(g1)
		var u Utils
		sharedG := u.base64StringToElementBytes(g1)
		pairing, _ := pbc.NewPairingFromString(sharedParams) // loads pairing parameters from a string and instantiates a pairing
		fmt.Println(pairing)
	*/
	// Channel for messages. Normally this would be a network connection.
	messageChannel := make(chan *messageData)
	// Channel for public key distribution. This might be a secure out-of-band
	// channel or something like a web of trust. The public key only needs to
	// be transmitted and verified once. The best way to do this is beyond the
	// scope of this example.
	keyChannel := make(chan []byte)

	// Channel to wait until both simulations are done
	finished := make(chan bool)

	// Simulate the conversation participants
	go alice(sharedParams, sharedG, messageChannel, keyChannel, finished)
	go bob(sharedParams, sharedG, messageChannel, keyChannel, finished)
	<-finished
	<-finished

	// Wait for the communication to finish

}

// Alice generates a keypair and signs a message
func alice(sharedParams string, sharedG []byte, messageChannel chan *messageData, keyChannel chan []byte, finished chan bool) {
	fmt.Println(sharedG)
	fmt.Println(sharedParams)

	// Alice loads the system parameters
	pairing, _ := pbc.NewPairingFromString(sharedParams) // loads pairing parameters from a string and instantiates a pairing
	fmt.Println("12")
	jsonStringPub := Utils{}.readFile("../config/pubkey")
	//fmt.Println(jsonStringPub)
	pubkeyObject, _ := ourjson.ParseObject(jsonStringPub)
	gk, _ := pubkeyObject.GetString("g")
	l := pairing.NewG2().SetBytes(Utils{}.base64StringToElementBytes(gk))
	fmt.Println(l)
	fmt.Println(pairing)
	sharedG = Utils{}.base64StringToElementBytes(gk)
	g := pairing.NewG2().SetBytes(sharedG)
	fmt.Println("g")
	fmt.Println(sharedG)
	fmt.Println(g)

	// Generate keypair (x, g^x)
	privKey := pairing.NewZr().Rand()
	pubKey := pairing.NewG2().PowZn(g, privKey)
	fmt.Println(pubKey)

	// Send public key to Bob
	keyChannel <- pubKey.Bytes()

	// Some time later, sign a message, hashed to h, as h^x
	message := "some text to sign"
	h := pairing.NewG1().SetFromStringHash(message, sha256.New())
	signature := pairing.NewG2().PowZn(h, privKey)

	// Send the message and signature to Bob
	messageChannel <- &messageData{message: message, signature: signature.Bytes()}

	finished <- true
}

// Bob verifies a message received from Alice
func bob(sharedParams string, sharedG []byte, messageChannel chan *messageData, keyChannel chan []byte, finished chan bool) {
	fmt.Println(sharedG)
	fmt.Println(sharedParams)
	// Bob loads the system parameters
	pairing, _ := pbc.NewPairingFromString(sharedParams)
	g := pairing.NewG2().SetBytes(sharedG)

	// Bob receives Alice's public key (and presumably verifies it manually)
	pubKey := pairing.NewG2().SetBytes(<-keyChannel)

	// Some time later, Bob receives a message to verify
	data := <-messageChannel
	signature := pairing.NewG1().SetBytes(data.signature)

	// To verify, Bob checks that e(h,g^x)=e(sig,g)
	h := pairing.NewG1().SetFromStringHash(data.message, sha256.New())
	temp1 := pairing.NewGT().Pair(h, pubKey)
	temp2 := pairing.NewGT().Pair(signature, g)
	if !temp1.Equals(temp2) {

		fmt.Println("*BUG* Signature check failed *BUG*")
	} else {
		fmt.Println("Signature verified correctly")
	}
	finished <- true
}
