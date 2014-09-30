package main

import "testing"

func TestCalculateChecksumZero(t *testing.T) {
	// Create a packet with only 0s
	pkt := &antpacket{}
	if pkt.genChecksum() != byte(0) {
		t.Fail()
	}
}

func TestValidateChecksum(t *testing.T) {
	// Create deliberately wrong checksum
	pkt := &antpacket{}
	pkt.checksum = 7

	if pkt.validateChecksum() != false {
		t.Fail()
	}
}
