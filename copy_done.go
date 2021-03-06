package pgproto

import (
	"fmt"
	"io"
)

type CopyDone struct {
}

func (c *CopyDone) server() {}

func ParseCopyDone(r io.Reader) (*CopyDone, error) {
	b := newReadBuffer(r)

	// 'c' [int32 - length]
	err := b.ReadTag('c')
	if err != nil {
		return nil, err
	}

	l, err := b.ReadInt()
	if err != nil {
		return nil, err
	}

	if l != 4 {
		return nil, fmt.Errorf("expected message length of 4")
	}

	return &CopyDone{}, nil
}

// Encode will return the byte representation of this message
func (c *CopyDone) Encode() []byte {
	// 'c' [int32 - length]
	return []byte{
		// Tag
		'c',
		// Length
		'\x00', '\x00', '\x00', '\x04',
	}
}

func (c *CopyDone) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type":    "CopyDone",
		"Payload": nil,
	}
}

func (c *CopyDone) String() string { return messageToString(c) }
