// The Antbuffer buffers and manages the ongoing communication between the host and the ant stick
//
// The Antbuffer can be told to set up and listen on a channel. By default, the
// Antbuffer will keep that listen open forever until the end of the program.
package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"time"

	"github.com/yokujin/gousb/usb"
)

// The Antbuffer is the point of control over the serial interface to the ant stick.
type Antbuffer struct {
	epin              usb.Endpoint
	epout             usb.Endpoint
	readChan          chan []byte
	writeChan         chan *antpacket
	channelListenners []chan bytes.Buffer
}

// TODO: actually listen for errors
// Could be done through the read daemon sending messages to the error channel of the Antbuffer
// Then the wait function could get a message from a secondary channel

// NewAntbuffer creates a new Antbuffer and populates network key 0x01 with the given network key unless nil.
func NewAntbuffer(epin, epout usb.Endpoint, networkKey []byte) (*Antbuffer, error) {
	// Create read channel
	readChan := make(chan []byte, 20)
	// Create write channel
	writeChan := make(chan *antpacket, 20)

	// Initialize Antbuffer
	antbuf := &Antbuffer{
		epin,
		epout,
		readChan,
		writeChan,
		make([]chan bytes.Buffer, 6), //TODO: make actual device limit of channels
	}

	// Launch listener daemon
	go antbuf.readDaemon()

	// Reset
	_, err := antbuf.GenSendAndWait(SystemReset, 0)
	if err != nil {
		return nil, err
	}

	// Set network 1 with ant plus network key
	if len(networkKey) != 8 {
		return nil, ErrNetworkKeyLength
	}
	_, err = antbuf.GenSendAndWait(
		SetNetwork,
		0x01,
		networkKey[0],
		networkKey[1],
		networkKey[2],
		networkKey[3],
		networkKey[4],
		networkKey[5],
		networkKey[6],
		networkKey[7],
	)
	if err != nil {
		return nil, err
	}

	return antbuf, nil
}

// SetupChannel will begin listening for the device specified by dev, initializing it on given channel.
// Returns a channel which contains events generated on that channel.
func (a *Antbuffer) SetupChannel(channel byte, dev *Antdevicetype) (<-chan bytes.Buffer, error) {
	// Setup Channel Type (Assign Channel)
	// TODO: Network should not be a magic number
	_, err := a.GenSendAndWait(AssignChannel, channel, dev.ChannelType, 0x1)
	if err != nil {
		return nil, err
	}

	// Set Channel Frequency (ChannelRFFrequency)
	_, err = a.GenSendAndWait(SetChannelRFFrequency, channel, dev.RFChannelFreq)
	if err != nil {
		return nil, err
	}

	// Setup Chanel Device Number, Device Type & Transmission Type (Set Channel ID)
	// TODO: Library for this
	devnum := &bytes.Buffer{}
	err = binary.Write(devnum, binary.LittleEndian, dev.DeviceNumber)
	numbytes := devnum.Bytes()
	// TODO: Allow for pairing bit or not on DeviceType
	_, err = a.GenSendAndWait(SetChannelID, channel, numbytes[0], numbytes[1], dev.DeviceType, dev.TransmissionType)
	if err != nil {
		return nil, err
	}

	// Setup Channel Messaging Period (Channel Period)
	period := &bytes.Buffer{}
	err = binary.Write(period, binary.LittleEndian, dev.ChannelPeriod)
	perbytes := period.Bytes()
	_, err = a.GenSendAndWait(SetChannelPeriod, channel, perbytes[0], perbytes[1])
	if err != nil {
		return nil, err
	}

	// Setup Channel Search Timeout (Channel Search Timeout)
	_, err = a.GenSendAndWait(SetSearchTimeout, channel, dev.SearchTimeout)
	if err != nil {
		return nil, err
	}

	// Open Channel!
	_, err = a.GenSendAndWait(OpenChannel, channel)
	if err != nil {
		return nil, err
	}

	// Create listen channel
	retChannel := make(chan bytes.Buffer, 20)

	// Register listen channel with Antbuffer
	// TODO: Actually use this
	a.channelListenners[int(channel)] = retChannel

	return retChannel, err

}

// CloseChannel will close the specified channel
// Currently, this function will discard all other messages until closed
// TODO: non-blocking
func (a *Antbuffer) CloseChannel(channel byte) error {
	log.Println("Closing channels...")
	pkt, _ := a.GenSendAndWait(CloseChannel, channel)
	if pkt.id == ChannelResponseOrEvent {
		// If we didn't get an ack
		if pkt.data[2] != 0x00 {
			return ErrAntInvalidMisc
		}
	}
	// Wait for complete close
	log.Println("Waiting for confirmation...")
	for {
		pkt, err := a.Wait()
		if err != nil {
			log.Fatalln("Error while closing, ", err)
		}
		if pkt.id == ChannelResponseOrEvent {
			if pkt.data[2] == 0x07 && pkt.data[0] == channel {
				// Channel was successfully closed
				log.Println("Successfully closed channel ", pkt.data[0])
				break
			}
		}
	}

	return nil
}

// TODO: Error channel
// TODO: Submitting to parser
// Parser distributes to error handler, channel handlers and others

// readDaemon is the goroutine which holds the read endpoint of the ant stick.
// It forwards read antpackets to the antbuffer for parsing and distribution.
func (a *Antbuffer) readDaemon() {
	// Read forever
	for {
		buf := make([]byte, maxDataLength, maxDataLength)
		_, err := a.epin.Read(buf)
		if err == usb.ERROR_TIMEOUT {
			// Timeout
			continue
		}

		// TODO: make "no device" error more pretty

		if err != nil {
			log.Fatalln("Error reading from endpoint, ", err)
			break
		}
		// Send out
		a.readChan <- buf
	}
}

// TODO some kind of error returning

func (a *Antbuffer) writeDaemon() {
}

// GenSendAndWait - Generate an antpacket, send and await reply
func (a *Antbuffer) GenSendAndWait(pktdetails ...byte) (*antpacket, error) {
	// TODO: Debug flag for this
	pkt, err := GenerateAntpacket(pktdetails[0], pktdetails[1:]...)
	if err != nil {
		return nil, err
	}
	log.Println("OUT: ", pkt)

	// Send
	a.Send(pkt)
	if err != nil {
		return nil, err
	}

	// Wait
	pkt, err = a.Wait()
	if err != nil {
		return nil, err
	}

	return pkt, err
}

// Send packet
// TODO: use a daemon
func (a *Antbuffer) Send(pkt *antpacket) error {
	outBuf := new(bytes.Buffer)
	_, err := pkt.toBinary(outBuf)
	if err != nil {
		return err
	}

	a.epout.Write(outBuf.Bytes())
	// a.writeChan <- pkt
	return nil
}

// Wait blocks while listening for a reply. This function will be deprecated soon.
func (a *Antbuffer) Wait() (*antpacket, error) {
	log.Println("Waiting for reply...")
	select {
	case read := <-a.readChan:
		pkt, err := readAntpacket(read)
		if err != nil {
			return nil, err
		}
		log.Printf("IN: %v\n", pkt)
		return pkt, nil
	case <-time.After(1 * time.Second):
		return nil, ErrAntTimedout
	}
}

// RegisterHandler tegisters a handler on a channel for a specific class of ant packets.
func (a *Antbuffer) RegisterHandler(channel int, class byte, receiving chan<- *antpacket) {

}
