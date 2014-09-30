package main

import (
	"fmt"
	"github.com/yokujin/gousb/usb"
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

	// Check for kernel driver involvement and detach if necessary

	fmt.Println("Opening Endpoint...")
	_, err = antdev.OpenEndpoint(
		uconf,
		uiface,
		usetup,
		uint8(uep)|uint8(usb.ENDPOINT_DIR_OUT),
	)

	if err != nil {
		fmt.Println("Error opening endpoint, ", err)
		return
	}

}
