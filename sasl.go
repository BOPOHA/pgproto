package pgproto

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
	// 'p' [int32 - length] []byte
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
