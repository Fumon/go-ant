package main

import (
	"errors"
	"fmt"
	"github.com/yokujin/gousb/usb"
	"log"
	"os"
	"os/signal"
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

	ctx.Debug(3)

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

	// Get the ant plus network key
	key, err := getNetworkKey()
	if err != nil {
		log.Fatalln("Error getting key, ", err)
	}

	// Create antbuffer
	antbuf, err := NewAntbuffer(epRead, epWrite, key)
	if err != nil {
		log.Fatalln("Error in creating antbuffer, ", err)
	}

	// ListenForCheststrap
	// TODO: this should be using the returned channel to listen
	// All errors at a higher than channel level are to be handled
	// by the Antbuffer
	_, err = antbuf.SetupChannel(0x01, heartrate)
	if err != nil {
		log.Fatalln("Error listening to Heart Rate sensor, ", err)
	}
	// TODO: This should be implicit in the Antbuffer
	defer func() {
		log.Println("Closing channels...")
		antbuf.GenSendAndWait(CloseChannel, 0x01)
		// Wait for complete close
		log.Println("Waiting for confirmation...")
		for {
			pkt, err := antbuf.Wait()
			if err != nil {
				log.Fatalln("Error while closing, ", err)
			}
			if pkt.id == ChannelResponseOrEvent {
				if pkt.data[2] == 0x07 {
					// Channel was successfully closed
					log.Println("Successfully closed channel ", pkt.data[0], " proceeding to exit...")
					break
				}
			}
		}
	}()

	// TODO: Move this somewhere else
	// Catch close signal
	killchan := make(chan os.Signal, 1)
	signal.Notify(killchan, os.Interrupt, os.Kill)

	// Listen for everything forever
readloop:
	for {
		// Die if killed
		select {
		case <-killchan:
			log.Println("Recieved KILL!")
			break readloop
		case <-time.After(10 * time.Millisecond):
		}
		pkt, err := antbuf.Wait()
		if err == ErrAntTimedout {
			// Depending on stuff... might need to relisten for device
			// Device relisten should really be based off of error
			// events from the stick.
			continue
		} else if err != nil {
			log.Fatalln("Error in waiting, ", err)
		}
		log.Println(pkt)
	}

	// Exiting
	fmt.Println("Exiting...")
}

func getNetworkKey() (key []byte, err error) {
	// Network Key
	// TODO: Publish additional go binary to write the key from command line
	// TODO: Make flag for where this is
	// TODO: Note this in documentation
	file, err := os.Open("/etc/ant/antPlusNetworkKey")
	if err != nil {
		return nil, errors.New(fmt.Sprint("Error opening ant key file in /etc/ant/antPlusNetworkKey, ", err))
	}
	defer file.Close()

	// The network key is 8 bytes long
	key = make([]byte, 8)
	n, err := file.Read(key)
	if err != nil {
		return nil, err
	} else if n != 8 {
		return nil, ErrNetworkKeyLength
	}

	return key, nil
}
