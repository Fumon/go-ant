package main

type anterror string

func (a anterror) Error() string {
	return string(a)
}

//Errors
const (
	ErrArgumentsNil        = anterror("Arguments list must not be nil")
	ErrArgumentsLen        = anterror("Arguments length mismatch")
	ErrUnknownClass        = anterror("Unknown message class")
	ErrMinimumPacketLength = anterror("Packet is smaller than minimum length")
	ErrChecksumMismatch    = anterror("Checksum Mismatch")
	ErrNetworkKeyLength    = anterror("Network key not of correct length")
	ErrAntTimedout         = anterror("Timed out waiting for a reply from ant stick")
	ErrAntInvalidMisc      = anterror("There was a problem, there should be a better error message than this but it hasn't been coded yet")
)
