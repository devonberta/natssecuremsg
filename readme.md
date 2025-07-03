
# natssecuremsg – Secure Messaging in Go with nkeys

`natssecuremsg` is a lightweight Go module that wraps typed payloads in signed and encrypted messages using [NATS nkeys](https://github.com/nats-io/nkeys).

It is ideal for end-to-end secure messaging in systems where:
- Messages need to be encrypted with public keys (X25519)
- Payload authenticity is guaranteed using Ed25519 signatures
- Replay protection and verification are important

## ✨ Features

- ✅ Generic payload support via `SecureMessage[T]`
- ✅ Ed25519 digital signature (`sign` + `verify`)
- ✅ X25519 encryption (`seal` + `open`)
- ✅ Nonce and timestamp embedded in every message
- ✅ Stateless: no session needed to verify or decrypt

## 📥 Installation

```bash
go get github.com/devonberta/natssecuremsg
```

## 🧱 API Overview

### `NewSecureMessage[T](payload *T, signer nkeys.KeyPair)`
Creates a new message with payload `T`, signed using Ed25519.

### `(msg *SecureMessage[T]) Encrypt(sender nkeys.KeyPair, recipientPub string)`
Encrypts the signed message using X25519 and the recipient’s public encryption key.

### `(msg *SecureMessage[T]) Decrypt(receiver nkeys.KeyPair, senderEncPub, senderSigPub string)`
Decrypts and verifies a message using the sender’s public keys.

## 🚀 Example Application

Here's a full round-trip example with two parties (Alice and Bob), each with their own signing and encryption keypairs:

### ✍️ Step 1: Setup Keypairs

```go
aliceSig, _ := nkeys.CreateUser()
aliceEnc, _ := nkeys.CreateCurveKeys()

bobSig, _ := nkeys.CreateUser()
bobEnc, _ := nkeys.CreateCurveKeys()
```

### 📨 Step 2: Alice Sends Secure Message to Bob

```go
payload := &MessagePayload{From: "Alice", Msg: "Hello Bob"}

msg, _ := securemsg.NewSecureMessage(payload, aliceSig)
msg.Encrypt(aliceEnc, bobEncPub) // Bob's public key
```

### 📬 Step 3: Bob Receives, Decrypts, and Verifies

```go
var received securemsg.SecureMessage[MessagePayload]
received.Encrypted = msg.Encrypted
received.Decrypt(bobEnc, aliceEncPub, aliceSigPub)
fmt.Println("Received from Alice:", received.Payload)
```

### 🔁 Step 4: Bob Sends a Response to Alice

```go
reply := &MessagePayload{From: "Bob", Msg: "Hi Alice"}
response, _ := securemsg.NewSecureMessage(reply, bobSig)
response.Encrypt(bobEnc, aliceEncPub) // Alice's public key
```

### 📩 Step 5: Alice Receives, Decrypts, and Verifies

```go
var back securemsg.SecureMessage[MessagePayload]
back.Encrypted = response.Encrypted
back.Decrypt(aliceEnc, bobEncPub, bobSigPub)
fmt.Println("Received from Bob:", back.Payload)
```

## 🛡️ Security Properties

- 🔐 **Confidentiality:** via X25519 public key encryption
- 🧾 **Integrity:** payloads are signed and verified with Ed25519
- ⏱️ **Freshness:** timestamp and nonce prevent replay attacks
- 🧾 **Stateless Verification:** no shared secrets or sessions required

## 🧪 Running the Example

The `example/main.go` demonstrates the full round-trip between two parties.

```bash
cd example
go run main.go
```

Expected output:

```
📥 Bob received from Alice: &{From:Alice Msg:Hello Bob}
📥 Alice received from Bob: &{From:Bob Msg:Hello Alice}
```

