package main

// Config Messages, HOST -> ANT
const (
	UnassignChannel             = 0x41
	AssignChannel               = 0x42
	SetChannelID                = 0x51
	SetChannelPeriod            = 0x43
	SetSearchTimeout            = 0x44
	SetChannelRFFrequency       = 0x45
	SetNetwork                  = 0x46
	SetTransmitPower            = 0x47
	IDListAdd                   = 0x59
	IDListConfig                = 0x5A
	SetChannelTransmitPower     = 0x60
	SetLowPrioritySearchTimeout = 0x63
	SetSerialNumberSetChannelID = 0x65
	EnableExtRXMesgs            = 0x66
	EnableLED                   = 0x68
	CrystalEnable               = 0x6D
	LibConfig                   = 0x6E
	FrequencyAgility            = 0x70
	SetProximitySearch          = 0x71
	SetChannelSearchPriority    = 0x75
)

// Notifications ANT -> HOST
const (
	StartupMessage     = 0x6F
	SerialErrorMessage = 0xAE
)

// Control Messages HOST->ANT
const (
	SystemReset    = 0x4A
	OpenChannel    = 0x4B
	CloseChannel   = 0x4C
	OpenRxScanMode = 0x5B
	RequestMessage = 0x4D
	SleepMessage   = 0xC5
)

// Data Messages HOST<->ANT
const (
	BroadcastData     = 0x4E
	AcknowledgeData   = 0x4F
	BurstTransferData = 0x50
)

// Channel/Event Messages
const (
	ChannelResponseOrEvent = 0x40
)

// Requested Response ANT->HOST
const (
	ChannelStatus = 0x52
	ChannelID     = 0x51
	ANTVersion    = 0x3E
	Capabilities  = 0x54
	SerialNumber  = 0x61
)

// Test Mode HOST->ANT
const (
	CWInit = 0x53
	CWTest = 0x48
)

type msgClass struct {
	name          string
	datalength    int
	dataFieldDesc []string
}

var msgClasses = map[int]msgClass{
	UnassignChannel: msgClass{
		"Unassign Channel",
		1,
		[]string{
			"Channel Number",
		},
	},
	AssignChannel: msgClass{
		"Assign Channel",
		3,
		[]string{
			"Channel Number",
			"Channel Type",
			"Network Number",
		},
	},
	SetChannelID: msgClass{
		"Assign Channel",
		3,
		[]string{
			"Channel Number",
			"Channel Type",
			"Network Number",
		},
	},
}
