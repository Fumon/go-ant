package main

import (
	"fmt"
	"github.com/yokujin/gousb/usb"
	"log"
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
		fmt.Println("ERROR! ", err)
		return
	}

	// Exit if no devices opened
	if len(devs) == 0 {
		fmt.Println(devs)
		fmt.Println("No devices found")
		return
	}

	// Pick off the first device

	antdev := devs[0]

	fmt.Println("Opening Endpoint...")
	ep_read, err := antdev.OpenEndpoint(
		uconf,
		uiface,
		usetup,
		uint8(uep)|uint8(usb.ENDPOINT_DIR_IN),
	)

	if err != nil {
		fmt.Println("Error opening endpoint, ", err)
		return
	}

	// Create read channel
	readChan := make(chan []byte, 20)

	// Launch listener daemon
	go func() {
		// Read forever
		for {
			buf := make([]byte, maxDataLength, maxDataLength)
			_, err := ep_read.Read(buf)
			if err == usb.ERROR_TIMEOUT {
				// Timeout
				continue
			}

			if err != nil {
				fmt.Println("Error reading from endpoint, ", err)
				break
			}
			// Send out
			readChan <- buf
		}
	}()

	// Send a reset
	//outBuf := make([]byte, 12)

	//ep.Write()

	read := <-readChan
	log.Println(read)

	// Exiting
	fmt.Println("Exiting...")
}
