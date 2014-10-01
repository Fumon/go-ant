package main

type anterror string

func (a anterror) Error() string {
	return string(a)
}

//Errors
const (
	ErrArgumentsNil        = anterror("Argument list must not be nil")
	ErrArgumentsLen        = anterror("Insufficient arguments")
	ErrUnknownClass        = anterror("Unknown message class")
	ErrMinimumPacketLength = anterror("Packet is smaller than minimum length")
	ErrChecksumMismatch    = anterror("Checksum Mismatch")
)
