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

func TestReadValues(t *testing.T) {
	testPkt := &antpacket{
		syncByte,
		0x4,
		0x32,
		[]byte{1, 2, 3, 4},
		0,
	}
	testPkt.setChecksum()

	buf := new(bytes.Buffer)
	testPkt.toBinary(buf)

	readPacket, err := readAntpacket(buf.Bytes())
	if err != nil {
		t.Fatal("Error reading packet, ", err)
	}

	if readPacket.sync != syncByte || readPacket.msglen != 0x4 || readPacket.id != 0x32 {
		t.Fail()
	}

	if len(readPacket.data) != len(testPkt.data) {
		t.Fail()
	}

	for i, x := range readPacket.data {
		if testPkt.data[i] != x {
			t.Fail()
		}
	}

	if readPacket.checksum != testPkt.checksum {
		t.Fail()
	}
}

func TestReadChecksumValidate(t *testing.T) {
	// Sane packet but incorrect checksum
	testPkt := &antpacket{
		syncByte,
		0x4,
		0x32,
		[]byte{1, 2, 3, 4},
		0,
	}
	buf := new(bytes.Buffer)
	testPkt.toBinary(buf)

	_, err := readAntpacket(buf.Bytes())

	if err == nil {
		t.Fatal("Failed to reject incorrect checksum")
	}
}

func TestReadLength(t *testing.T) {
}
