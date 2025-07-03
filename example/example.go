package main

import (
	_"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nkeys"
	securemsg "github.com/devonberta/natssecuremsg"
)

type MessagePayload struct {
	From string `json:"from"`
	Msg  string `json:"msg"`
}

func main() {
	// Key generation
	aliceSig, _ := nkeys.CreateUser()
	aliceEnc, _ := nkeys.CreateCurveKeys()

	bobSig, _ := nkeys.CreateUser()
	bobEnc, _ := nkeys.CreateCurveKeys()

	// Public keys
	aliceEncPub, _ := aliceEnc.PublicKey()
	aliceSigPub, _ := aliceSig.PublicKey()
	bobEncPub, _ := bobEnc.PublicKey()
	bobSigPub, _ := bobSig.PublicKey()

	// Alice sends message to Bob
	alicePayload := &MessagePayload{From: "Alice", Msg: "Hello Bob"}

	aliceMsg, err := securemsg.NewSecureMessage(alicePayload, aliceSig)
	if err != nil {
		log.Fatal("Alice signing error:", err)
	}
	if err := aliceMsg.Encrypt(aliceEnc, bobEncPub); err != nil {
		log.Fatal("Alice encryption error:", err)
	}

	// Simulate transmission to Bob
	var msgToBob securemsg.SecureMessage[MessagePayload]
	msgToBob.Encrypted = aliceMsg.Encrypted

	if err := msgToBob.Decrypt(bobEnc, aliceEncPub, aliceSigPub); err != nil {
		log.Fatal("Bob decryption error:", err)
	}
	fmt.Println("📥 Bob received from Alice:", msgToBob.Payload)

	// Bob responds to Alice
	bobPayload := &MessagePayload{From: "Bob", Msg: "Hello Alice"}
	bobMsg, err := securemsg.NewSecureMessage(bobPayload, bobSig)
	if err != nil {
		log.Fatal("Bob signing error:", err)
	}
	if err := bobMsg.Encrypt(bobEnc, aliceEncPub); err != nil {
		log.Fatal("Bob encryption error:", err)
	}

	// Simulate transmission to Alice
	var msgToAlice securemsg.SecureMessage[MessagePayload]
	msgToAlice.Encrypted = bobMsg.Encrypted

	if err := msgToAlice.Decrypt(aliceEnc, bobEncPub, bobSigPub); err != nil {
		log.Fatal("Alice decryption error:", err)
	}
	fmt.Println("📥 Alice received from Bob:", msgToAlice.Payload)
}
