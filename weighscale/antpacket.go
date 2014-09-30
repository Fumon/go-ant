package main

// This is a marshallable datastructure for constructing and decoding ant packets

const (
	maxDataLength = 56
)

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

// Checksum the packet
func (a *antpacket) genChecksum() (chk byte) {
	// XOR everything
	chk = a.sync ^ a.msglen ^ a.id
	for _, e := range a.data {
		chk = chk ^ e
	}
	return
}

// Set the checksum for a constructed packet
func (a *antpacket) setChecksum() {
	a.checksum = a.genChecksum()
}

// Validate a packet by checksum
// Returns true if valid
func (a *antpacket) validateChecksum() bool {
	if a.genChecksum() == a.checksum {
		return true
	}
	return false
}
