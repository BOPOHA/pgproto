package pgproto

import (
	"bytes"
	"io"
)

type PasswordMessage struct {
	HeaderMessage []byte
	BodyMessage   []byte
}

func (p *PasswordMessage) client() {}

func ParsePasswordMessage(r io.Reader) (*PasswordMessage, error) {
	b := newReadBuffer(r)

	// 'p' [int32 - length] [string] \0
	err := b.ReadTag('p')
	if err != nil {
		return nil, err
	}

	buf, err := b.ReadLength()
	if err != nil {
		return nil, err
	}
	headerMessage, err := buf.ReadString(true)
	if err != nil {
		return nil, err
	}
	//Checking SASLInitialResponse (F)
	//Byte1('p') Identifies the message as an initial SASL response.
	//Int32 Length of message contents in bytes, including self.
	//String Name of the SASL authentication mechanism that the client selected.
	//Int32 Length of SASL mechanism specific "Initial Client Response" that follows, or -1 if there is no Initial Response.
	//Byten SASL mechanism specific "Initial Response".
	_, err = buf.ReadInt()

	if err != nil {
		// seems to be it will be always EOF, and ...
		// we have not any data in buffer, so it is valid Password Messsage
		return &PasswordMessage{HeaderMessage: headerMessage}, nil
	}
	bodyMessage, err := buf.ReadString(true)
	if err != nil {
		return nil, err
	}
	return &PasswordMessage{HeaderMessage: headerMessage, BodyMessage: bodyMessage}, nil
}

func (p *PasswordMessage) PasswordValid(user []byte, password []byte, salt []byte) bool {
	hash := HashPassword(user, password, salt)
	return bytes.Equal(p.HeaderMessage, hash)
}

func (p *PasswordMessage) SetPassword(user []byte, password []byte, salt []byte) {
	p.HeaderMessage = HashPassword(user, password, salt)
}

func (p *PasswordMessage) Encode() []byte {
	// PasswordMessage
	// 'p' [int32 - length] [string] \0
	//
	// SASLInitialResponse
	// 'p' [int32 - length] [string] [int32 - length] []byte \0
	w := newWriteBuffer()
	w.WriteString(p.HeaderMessage, true)
	if len(p.BodyMessage) > 0 {
		w.WriteInt(len(p.BodyMessage))
		w.WriteBytes(p.BodyMessage)
	}
	w.Wrap('p')
	return w.Bytes()
}

func (p *PasswordMessage) AsMap() map[string]interface{} {
	if len(p.BodyMessage) > 0 {
		return map[string]interface{}{
			"Type": "SASL " + string(p.HeaderMessage),
			"Payload": map[string]interface{}{
				"Mechanism": p.HeaderMessage,
				"Message":   p.BodyMessage,
			},
		}
	}
	return map[string]interface{}{
		"Type": "PasswordMessage",
		"Payload": map[string]interface{}{
			"Password": p.HeaderMessage,
		},
	}
}

func (p *PasswordMessage) String() string { return messageToString(p) }
