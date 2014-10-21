package main

import (
	"bytes"
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	_ "github.com/lib/pq"

	"github.com/streadway/amqp"
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

	// Open a file to store the info we get from the weighscale on the first connect. All of it.

	file, err := os.OpenFile(fmt.Sprint("/home/fumon/dk/weighscale_data/weighscale_log"), os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln("Error opening file, ", err)
	}
	defer file.Close()

	// Connect to DB
	db, err := sql.Open("postgres", "user=inserter dbname='quantifiedSelf' sslmode=disable")
	if err != nil {
		log.Fatalln("Error connecting to database, ", err)
	}
	defer db.Close()

	// Open AMQP
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalln("Failed to connect to amqp")
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalln("Failed to create channel, ")
	}
	defer ch.Close()

	channelName := "weights"

	log.Println("Opening channel ", channelName)
	q, err := ch.QueueDeclare(
		channelName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalln("Could not declare Queue")
	}

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

	// TODO: Move this somewhere else
	// Catch close signal
	killchan := make(chan os.Signal, 1)
	signal.Notify(killchan, os.Interrupt, os.Kill)

	// Listen for everything forever
	shouldWait := false
	die := false
readloop:
	for {
		if shouldWait {
			log.Println("Waiting a while so that we don't catch another weighin too soon")
			<-time.After(2 * time.Minute)
		}
		// TODO: this should be using the returned channel to listen
		// All errors at a higher than channel level are to be handled
		// by the Antbuffer
		_, err = antbuf.SetupChannel(0x01, weighscale)
		if err != nil {
			log.Fatalln("Error listening to Heart Rate sensor, ", err)
		}

		shouldWait = true
	weighinloop:
		for {
			// Die if killed
			select {
			case <-killchan:
				log.Println("Recieved KILL!")
				die = true
				break weighinloop
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
			// Check for weight packet
			if pkt.id == BroadcastData && pkt.data[1] == 0x01 {
				// Grab weight
				var weight uint16
				readbuf := bytes.NewBuffer(pkt.data[7:])
				binary.Read(readbuf, binary.LittleEndian, &weight)
				weightFactor := float64(weight) / 100.0

				// Write to file & db
				log.Println("Got a weight of ", weightFactor, "kg or ", (2.204 * weightFactor), "lbs\n\tWriting to file and shutting down.")
				_, err = file.WriteString(fmt.Sprint(time.Now().UTC().Unix(), "\t", time.Now().UTC(), "\t", weightFactor, "\t", (2.204 * weightFactor), "\n"))
				if err != nil {
					log.Fatalln("Problem writing to file, ", err)
				}

				var insertid int64
				err := db.QueryRow(fmt.Sprintf("INSERT INTO buffer.weight (date, weight) VALUES (to_timestamp(%v), %v) RETURNING did", time.Now().UTC().Unix(), weight)).Scan(&insertid)
				if err != nil {
					log.Println("ERROR inserting into db, ", err)
					break readloop
				}
				log.Println("DB did: ", insertid)

				// Send alert to processor
				buf := new(bytes.Buffer)
				err = binary.Write(buf, binary.LittleEndian, insertid)
				if err != nil {
					log.Println("Could not write binary to buffer, ", err)
					break readloop
				}

				err = ch.Publish(
					"",
					q.Name,
					false,
					false,
					amqp.Publishing{
						ContentType: "text/plain",
						Body:        buf.Bytes(),
					})
				if err != nil {
					log.Fatalln("Did not send correctly")
				}

				break weighinloop
			}
		}
		antbuf.CloseChannel(0x01)
		if die {
			break readloop
		}
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
