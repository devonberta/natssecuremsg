package natssecuremsg

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"github.com/nats-io/nkeys"
)

const SignatureSize = 64

type envelope[T any] struct {
	Timestamp time.Time `json:"ts"`
	Nonce     string    `json:"nonce"`
	Payload   T         `json:"payload"`
}

type SecureMessage[T any] struct {
	Payload   *T     `json:"-"`
	Nonce     string `json:"nonce,omitempty"`
	Timestamp time.Time `json:"ts,omitempty"`
	Signature []byte `json:"sig,omitempty"`
	Encrypted []byte `json:"enc,omitempty"`
}

// NewSecureMessage wraps the payload and signs it with the given signature keypair
func NewSecureMessage[T any](payload *T, signer nkeys.KeyPair) (*SecureMessage[T], error) {
	ts := time.Now().UTC()
	nonce := generateNonce()

	wrapped := envelope[T]{
		Timestamp: ts,
		Nonce:     nonce,
		Payload:   *payload,
	}

	encoded, err := json.Marshal(wrapped)
	if err != nil {
		return nil, err
	}

	sig, err := signer.Sign(encoded)
	if err != nil {
		return nil, err
	}

	return &SecureMessage[T]{
		Payload:   payload,
		Timestamp: ts,
		Nonce:     nonce,
		Signature: sig,
	}, nil
}

// Encrypt seals the signed message using X25519 for the recipient
func (m *SecureMessage[T]) Encrypt(sender nkeys.KeyPair, recipientPub string) error {
	wrapped := envelope[T]{
		Timestamp: m.Timestamp,
		Nonce:     m.Nonce,
		Payload:   *m.Payload,
	}
	plain, err := json.Marshal(wrapped)
	if err != nil {
		return err
	}
	combined := append(plain, m.Signature...)

	encrypted, err := sender.Seal(combined, recipientPub)
	if err != nil {
		return err
	}
	m.Encrypted = encrypted
	return nil
}

// Decrypt unseals the message using the recipient's private key and verifies the signature
func (m *SecureMessage[T]) Decrypt(receiver nkeys.KeyPair, senderEncPub, senderSigPub string) error {
	decrypted, err := receiver.Open(m.Encrypted, senderEncPub)
	if err != nil {
		return err
	}
	if len(decrypted) < SignatureSize {
		return errors.New("invalid decrypted length")
	}
	content := decrypted[:len(decrypted)-SignatureSize]
	m.Signature = decrypted[len(decrypted)-SignatureSize:]

	pubKP, err := nkeys.FromPublicKey(senderSigPub)
	if err != nil {
		return err
	}
	if err := pubKP.Verify(content, m.Signature); err != nil {
		return err
	}

	var wrapped envelope[T]
	if err := json.Unmarshal(content, &wrapped); err != nil {
		return err
	}
	m.Payload = &wrapped.Payload
	m.Nonce = wrapped.Nonce
	m.Timestamp = wrapped.Timestamp
	return nil
}

func generateNonce() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
