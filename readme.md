# SecureMsg - Encrypted & Signed Message Wrapper for Go

`nats-securemsg` is a lightweight Go module for securely exchanging typed messages using [nkeys](https://github.com/nats-io/nkeys). It supports:
- Signed messages with Ed25519
- Encrypted payloads using X25519 (sealed box)
- Typed generic payloads
- Automatic timestamp and nonce generation

## Installation

```bash
go get github.com/devonberta/nats-securemsg
```# natssecuremsg
