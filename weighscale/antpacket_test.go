package main

import (
	"bytes"
	"testing"
)

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

func TestToBinary(t *testing.T) {
	pkt := &antpacket{}
	buf := new(bytes.Buffer)

	_, err := pkt.toBinary(buf)
	if err != nil {
		t.Fatal("Function returned error, ", err)
	}

	byteslice := buf.Bytes()
	for _, x := range byteslice {
		if x != 0 {
			t.Fail()
		}
	}
}
