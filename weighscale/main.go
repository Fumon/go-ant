package main

import (
	"bytes"
	"fmt"
	"github.com/yokujin/gousb/usb"
	"log"
	"time"
)

const (
	dynastreamUsbVendid = 0x0fcf
	antstick            = 0x1008
	// USB endpoint info
	uconf  = 01
	uiface = 00
	usetup = 00
	uep    = 0x01
)

func main() {
	fmt.Println("- Life Begins -\n")

	// Get context
	ctx := usb.NewContext()
	defer ctx.Close()

	ctx.Debug(4)

	// Find and open the device
	devs, err := ctx.ListDevices(func(desc *usb.Descriptor) bool {
		if desc.Vendor == dynastreamUsbVendid && desc.Product == antstick {
			fmt.Println("Found antstick")
			return true
		}
		return false
	})

	defer func() {
		for _, d := range devs {
			d.Close()
		}
	}()

	if err != nil {
		log.Fatalln("ERROR! ", err)
		return
	}

	// Exit if no devices opened
	if len(devs) == 0 {
		log.Fatalln("No devices found")
		return
	}

	// Pick off the first device

	antdev := devs[0]

	log.Println("Opening Endpoints...")
	epRead, err := antdev.OpenEndpoint(
		uconf,
		uiface,
		usetup,
		uint8(uep)|uint8(usb.ENDPOINT_DIR_IN),
	)
	epWrite, err := antdev.OpenEndpoint(
		uconf,
		uiface,
		usetup,
		uint8(uep)|uint8(usb.ENDPOINT_DIR_OUT),
	)

	if err != nil {
		log.Println("Error opening endpoint, ", err)
		return
	}

	// Create read channel
	readChan := make(chan []byte, 20)

	// Launch listener daemon
	go func() {
		// Read forever
		for {
			buf := make([]byte, maxDataLength, maxDataLength)
			_, err := epRead.Read(buf)
			if err == usb.ERROR_TIMEOUT {
				// Timeout
				continue
			}

			if err != nil {
				log.Fatalln("Error reading from endpoint, ", err)
				break
			}
			// Send out
			readChan <- buf
		}
	}()

	log.Println("Sending reset packet...")
	// Send a reset
	outBuf := new(bytes.Buffer)
	resetPacket := &antpacket{
		syncByte,
		1,
		resetSystem,
		[]byte{0},
		0,
	}
	resetPacket.setChecksum()
	_, err = resetPacket.toBinary(outBuf)
	if err != nil {
		log.Fatalln("Error in writing packet to binary, ", err)
	}

	epWrite.Write(outBuf.Bytes())

	log.Println("Waiting for reply...")
	select {
	case read := <-readChan:
		pkt, err := readAntpacket(read)
		if err != nil {
			log.Fatalln("Error reading packet, ", err)
		}
		fmt.Printf("Reply: % X\n", pkt)
	case <-time.After(1 * time.Second):
		log.Println("Timedout")
	}

	// Exiting
	fmt.Println("Exiting...")
}
