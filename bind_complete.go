package pgproto

import (
	"bytes"
	"fmt"
	"io"
)

// '2' [int32 - length]
var rawBindCompleteMessage = [5]byte{
	// Tag
	'2',
	// Length
	'\x00', '\x00', '\x00', '\x04',
}

// BindComplete represents a server response message
type BindComplete struct{}

func (b *BindComplete) server() {}

// ParseBindComplete will attempt to read an BindComplete message from the io.Reader
func ParseBindComplete(r io.Reader) (*BindComplete, error) {
	b := newReadBuffer(r)

	var msg [5]byte
	_, err := b.Read(msg[:])
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(msg[:], rawBindCompleteMessage[:]) {
		return nil, fmt.Errorf("invalid bind complete message")
	}

	return &BindComplete{}, nil
}

// Encode will return the byte representation of this message
func (b *BindComplete) Encode() []byte {
	// '2' [int32 - length]
	return rawBindCompleteMessage[:]
}

// AsMap method returns a common map representation of this message:
//
//   map[string]interface{}{
//     "Type": "BindComplete",
//     "Payload": nil,
//     },
//   }
func (b *BindComplete) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type":    "BindComplete",
		"Payload": nil,
	}
}

func (b *BindComplete) String() string { return messageToString(b) }
