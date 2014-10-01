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
	class         string
	template      antpacketTemplate
	dataFieldDesc []string
}

var msgClasses = map[int]*msgClass{
	UnassignChannel: &msgClass{
		"Unassign Channel",
		"Config",
		antpacketTemplate{
			1,
			0x41,
		},
		[]string{
			"Channel Number",
		},
	},
	AssignChannel: &msgClass{
		"Assign Channel",
		"Config",
		antpacketTemplate{
			3,
			0x42,
		},
		[]string{
			"Channel Number",
			"Channel Type",
			"Network Number",
		},
	},
	ChannelID: &msgClass{
		"Set Channel ID / Respond",
		"Config / Requested Response",
		antpacketTemplate{
			5,
			0x51,
		},
		[]string{
			"Channel Number",
			"Device number(1/2)",
			"Device number(2/2)",
			"Device Type ID",
			"Trans. Type / Man ID",
		},
	},
	SetChannelPeriod: &msgClass{
		"Set Channel Period",
		"Config",
		antpacketTemplate{
			3,
			0x43,
		},
		[]string{
			"Channel Number",
			"Messaging Period(1/2)",
			"Messaging Period(2/2)",
		},
	},
	SetSearchTimeout: &msgClass{
		"Set Search Timeout",
		"Config",
		antpacketTemplate{
			2,
			0x44,
		},
		[]string{
			"Channel Number",
			"Search Timeout",
		},
	},
	SetChannelRFFrequency: &msgClass{
		"Set Channel RF Frequency",
		"Config",
		antpacketTemplate{
			2,
			0x45,
		},
		[]string{
			"Channel Number",
			"RF Frequency",
		},
	},
	SetNetwork: &msgClass{
		"Set Network",
		"Config",
		antpacketTemplate{
			9,
			0x46,
		},
		[]string{
			"Net #",
			"Key 0",
			"Key 1",
			"Key 2",
			"Key 3",
			"Key 4",
			"Key 5",
			"Key 6",
			"Key 7",
		},
	},
	SetTransmitPower: &msgClass{
		"Set Transmit Power",
		"Config",
		antpacketTemplate{
			2,
			0x47,
		},
		[]string{
			"0",
			"TX Power",
		},
	},
	IDListAdd: &msgClass{
		"ID List Add",
		"Config",
		antpacketTemplate{
			6,
			0x59,
		},
		[]string{
			"Channel Number",
			"Device number(1/2)",
			"Device number(2/2)",
			"Device Type ID",
			"Trans. Type",
			"List Index",
		},
	},
	IDListConfig: &msgClass{
		"ID List Config",
		"Config",
		antpacketTemplate{
			3,
			0x5A,
		},
		[]string{
			"Channel Number",
			"List Size",
			"Exclude",
		},
	},
	SetChannelTransmitPower: &msgClass{
		"Set Channel Transmit Power",
		"Config",
		antpacketTemplate{
			2,
			0x60,
		},
		[]string{
			"Channel Number",
			"TX Power",
		},
	},
	SetLowPrioritySearchTimeout: &msgClass{
		"Set Low Priority Search Timeout",
		"Config",
		antpacketTemplate{
			2,
			0x63,
		},
		[]string{
			"Channel Number",
			"Search Timeout",
		},
	},
	SetSerialNumberSetChannelID: &msgClass{
		"Set Serial Number Set Channel ID",
		"Config",
		antpacketTemplate{
			3,
			0x65,
		},
		[]string{
			"Channel Number",
			"Device Type ID",
			"Trans. Type",
		},
	},
	EnableExtRXMesgs: &msgClass{
		"Enable Ext RX Mesgs",
		"Config",
		antpacketTemplate{
			2,
			0x66,
		},
		[]string{
			"0",
			"Enable",
		},
	},
	EnableLED: &msgClass{
		"Enable LED",
		"Config",
		antpacketTemplate{
			2,
			0x68,
		},
		[]string{
			"0",
			"Enable",
		},
	},
	CrystalEnable: &msgClass{
		"Crystal Enable",
		"Config",
		antpacketTemplate{
			1,
			0x6D,
		},
		[]string{
			"0",
		},
	},
	LibConfig: &msgClass{
		"Lib Config",
		"Config",
		antpacketTemplate{
			2,
			0x6E,
		},
		[]string{
			"0",
			"Lib Config",
		},
	},
	FrequencyAgility: &msgClass{
		"Frequency Agility",
		"Config",
		antpacketTemplate{
			4,
			0x70,
		},
		[]string{
			"Channel Number",
			"Freq’ 1",
			"Freq’ 2",
			"Freq’ 3",
		},
	},
	SetProximitySearch: &msgClass{
		"Set Proximity Search",
		"Config",
		antpacketTemplate{
			2,
			0x71,
		},
		[]string{
			"Channel Number",
			"Search Threshold",
		},
	},
	SetChannelSearchPriority: &msgClass{
		"Set Channel Search Priority",
		"Config",
		antpacketTemplate{
			2,
			0x75,
		},
		[]string{
			"Channel Number",
			"Search Priority",
		},
	},
	StartupMessage: &msgClass{
		"Startup Message",
		"Notifications",
		antpacketTemplate{
			1,
			0x6F,
		},
		[]string{
			"Startup Message ",
		},
	},
	SerialErrorMessage: &msgClass{
		"Serial Error Message",
		"Notifications",
		antpacketTemplate{
			1,
			0xAE,
		},
		[]string{
			"Error Number ",
		},
	},
	SystemReset: &msgClass{
		"System Reset",
		"Control",
		antpacketTemplate{
			1,
			0x4A,
		},
		[]string{
			"0",
		},
	},
	OpenChannel: &msgClass{
		"Open Channel",
		"Control",
		antpacketTemplate{
			1,
			0x4B,
		},
		[]string{
			"Channel Number",
		},
	},
	CloseChannel: &msgClass{
		"Close Channel",
		"Control",
		antpacketTemplate{
			1,
			0x4C,
		},
		[]string{
			"Channel Number",
		},
	},
	OpenRxScanMode: &msgClass{
		"Open Rx Scan Mode",
		"Control",
		antpacketTemplate{
			1,
			0x5B,
		},
		[]string{
			"0",
		},
	},
	RequestMessage: &msgClass{
		"Request Message",
		"Control",
		antpacketTemplate{
			2,
			0x4D,
		},
		[]string{
			"Channel Number",
			"Message ID",
		},
	},
	SleepMessage: &msgClass{
		"Sleep Message",
		"Control",
		antpacketTemplate{
			1,
			0xC5,
		},
		[]string{
			"0",
		},
	},
	BroadcastData: &msgClass{
		"Broadcast Data",
		"Data",
		antpacketTemplate{
			9,
			0x4E,
		},
		[]string{
			"Channel Number",
			"Data0",
			"Data1",
			"Data2",
			"Data3",
			"Data4",
			"Data5",
			"Data6",
			"Data7",
		},
	},
	AcknowledgeData: &msgClass{
		"Acknowledge Data",
		"Data",
		antpacketTemplate{
			9,
			0x4F,
		},
		[]string{
			"Channel Number",
			"Data0",
			"Data1",
			"Data2",
			"Data3",
			"Data4",
			"Data5",
			"Data6",
			"Data7",
		},
	},
	BurstTransferData: &msgClass{
		"Burst Transfer Data",
		"Data",
		antpacketTemplate{
			9,
			0x50,
		},
		[]string{
			"Sequence/Channel Number",
			"Data0",
			"Data1",
			"Data2",
			"Data3",
			"Data4",
			"Data5",
			"Data6",
			"Data7",
		},
	},
	ChannelResponseOrEvent: &msgClass{
		"Channel Response / Event",
		"Channel / Event Messages",
		antpacketTemplate{
			3,
			0x40,
		},
		[]string{
			"Channel Number",
			"Message ID",
			"Message Code",
		},
	},
	ChannelStatus: &msgClass{
		"Channel Status",
		"Requested Response",
		antpacketTemplate{
			2,
			0x52,
		},
		[]string{
			"Channel Number",
			"Channel Status",
		},
	},
	ANTVersion: &msgClass{
		"ANT Version",
		"Requested Response",
		antpacketTemplate{
			11,
			0x3E,
		},
		[]string{
			"Ver0",
			"Ver1",
			"Ver2",
			"Ver 3|Ver 4",
			"Ver 5|Ver6",
			"Ver7",
			"Ver8",
			"Ver9",
			"Ver10",
		},
	},
	Capabilities: &msgClass{
		"Capabilities",
		"Requested Response",
		antpacketTemplate{
			6,
			0x54,
		},
		[]string{
			"Max Channels",
			"Max Networks",
			"Standard Options",
			"Advanced Options",
			"Adv’ Options 2",
			"Rsvd",
		},
	},
	SerialNumber: &msgClass{
		"Serial Number",
		"Requested Response",
		antpacketTemplate{
			4,
			0x61,
		},
		[]string{
			"Serial Number(1/4)",
			"Serial Number(2/4)",
			"Serial Number(3/4)",
			"Serial Number(4/4)",
		},
	},
	CWInit: &msgClass{
		"CW Init",
		"Test",
		antpacketTemplate{
			1,
			0x53,
		},
		[]string{
			"0",
		},
	},
	CWTest: &msgClass{
		"CW Test",
		"Test",
		antpacketTemplate{
			3,
			0x48,
		},
		[]string{
			"0",
			"TX Power",
			"RF Freq",
		},
	},
}
