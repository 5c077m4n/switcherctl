package parse

import "encoding/hex"

const (
	MessageLengthDefault = 165
	MessageLengthBreeze  = 168
	MessageLengthRunner  = 159
)

type DatagramParser struct{ message []byte }

func (dp *DatagramParser) IsSwitcher() (bool, error) {
	msgHex := hex.EncodeToString(dp.message)

	decoded, err := hex.DecodeString(msgHex[0:4])
	if err != nil {
		return false, err
	}

	return string(decoded) == "fef0" &&
		(len(dp.message) == MessageLengthDefault || len(dp.message) == MessageLengthBreeze || len(dp.message) == MessageLengthRunner), nil
}

func New(msg []byte) DatagramParser {
	return DatagramParser{message: msg}
}
