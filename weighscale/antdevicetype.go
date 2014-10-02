package main

// Convenience struct for defining channel properties
type Antdevicetype struct {
	ChannelType      byte
	RFChannelFreq    byte
	TransmissionType byte
	DeviceType       byte
	DeviceNumber     uint16
	ChannelPeriod    uint16
	SearchTimeout    byte
}

// TODO: new antchannel function

var weighscale *Antdevicetype = &Antdevicetype{
	0x00,
	57,
	0,
	119,
	0,
	8192, // 8192 counts
	0xFF, // Timeout should be as long as possible
}

var heartrate *Antdevicetype = &Antdevicetype{
	0x00,
	57,
	0,
	120,
	0,
	8070, // 8070 counts
	12,   // Search timeout 30 seconds
}
