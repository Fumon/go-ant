package main

// This is a marshallable datastructure for constructing and decoding ant packets

// List of message IDs. Incomplete.
const (
	// Config Messages, HOST -> ANT
	unassignChannel = 0x41
	assignChannel   = 0x42

	// Notifications ANT -> HOST
	startupMessage     = 0x6F
	serialErrorMessage = 0xAE

	// Control Messages
	resetSystem = 0x4A
)

// Encapsulates the basic structure of a standard ant packet.
//
// Byte #	-	Label
// 0	-	Sync
// 1	-	Message Length
// 2	- Message ID
// 3-(Message Length+2) - Data bytes (LSB ORDERING!)
// (Message Length + 3) - Checksum (XOR of all previous bytes including SYNC)
type antpacket struct {
	sync     byte
	msglen   byte
	id       byte
	data     []byte
	checksum byte
}
