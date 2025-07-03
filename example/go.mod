module example

go 1.23.3

require (
	github.com/devonberta/natssecuremsg v0.0.0
	github.com/nats-io/nkeys v0.4.11
)

require (
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
)

replace github.com/devonberta/natssecuremsg => ../
