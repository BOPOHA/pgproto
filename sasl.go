package pgproto

import (
	"io"
)

type SASLInitialResponse struct {
	Mechanism string
	Message   []byte
}

func (s *SASLInitialResponse) client() {}

func ParseSASLInitialResponse(r io.Reader) (*SASLInitialResponse, error) {
	b := newReadBuffer(r)

	// 'I' [int32 - length]
	err := b.ReadTag('p')
	if err != nil {
		return nil, err
	}

	b, err = b.ReadLength()
	if err != nil {
		return nil, err
	}

	mechanism, err := b.ReadString(true)
	if err != nil {
		return nil, err
	}
	_, err = b.ReadInt()
	if err != nil {
		return nil, err
	}

	a := &SASLInitialResponse{Mechanism: string(mechanism)}
	a.Message, err = b.ReadString(true)
	if err != nil {
		return nil, err
	}

	return a, nil
}

//SASLInitialResponse (F)
//Byte1('p') Identifies the message as an initial SASL response.
//Int32 Length of message contents in bytes, including self.
//String Name of the SASL authentication mechanism that the client selected.
//Int32 Length of SASL mechanism specific "Initial Client Response" that follows, or -1 if there is no Initial Response.
//Byten SASL mechanism specific "Initial Response".

func (s *SASLInitialResponse) Encode() []byte {
	// 'p' [int32 - length] [string] [int32 - length] []byte \0
	w := newWriteBuffer()
	w.WriteString([]byte(s.Mechanism), true)
	w.WriteInt(len(s.Message))
	w.WriteBytes(s.Message)
	w.Wrap('p')
	return w.Bytes()
}

func (s *SASLInitialResponse) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type": "SASL " + s.Mechanism,
		"Payload": map[string]interface{}{
			"Mechanism": s.Mechanism,
			"Message":   s.Message,
		},
	}
}

func (s *SASLInitialResponse) String() string { return messageToString(s) }

type SASLResponse struct {
	Message []byte
}

func (s *SASLResponse) client() {}
func (s *SASLResponse) server() {}

//SASLResponse (F)
//Byte1('p') Identifies the message as a SASL response.
//Int32 Length of message contents in bytes, including self.
//Byten SASL mechanism specific message data.

func (s *SASLResponse) Encode() []byte {
	// 'p' [int32 - length] []byte \0
	w := newWriteBuffer()
	w.WriteBytes(s.Message)
	w.Wrap('p')
	return w.Bytes()
}

func (s *SASLResponse) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type": "SASLResponse",
		"Payload": map[string]interface{}{
			"Message": s.Message,
		},
	}
}

func (s *SASLResponse) String() string { return messageToString(s) }
