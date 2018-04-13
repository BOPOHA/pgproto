package pgproto

import (
	"io"
)

type CopyBothFormat byte

type CopyBothResponse struct {
	Format        CopyBothFormat
	ColumnFormats []int
}

func (c *CopyBothResponse) server() {}

func ParseCopyBothResponse(r io.Reader) (*CopyBothResponse, error) {
	b := newReadBuffer(r)

	// 'W' [int32 - length] [int16 - count] [int16 - format] ...
	err := b.ReadTag('W')
	if err != nil {
		return nil, err
	}

	buf, err := b.ReadLength()
	if err != nil {
		return nil, err
	}

	format, err := buf.ReadByte()

	count, err := buf.ReadInt16()
	if err != nil {
		return nil, err
	}

	c := &CopyBothResponse{
		Format:        CopyBothFormat(format),
		ColumnFormats: make([]int, count),
	}

	for i := 0; i < count; i++ {
		c.ColumnFormats[i], err = buf.ReadInt16()
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *CopyBothResponse) Encode() []byte {
	// 'W' [int32 - length] [int16 - count] [int16 - format] ...
	w := newWriteBuffer()
	w.WriteByte(byte(c.Format))
	w.WriteInt16(len(c.ColumnFormats))
	for _, format := range c.ColumnFormats {
		w.WriteInt16(format)
	}
	w.Wrap('W')
	return w.Bytes()
}

func (c *CopyBothResponse) AsMap() map[string]interface{} {
	return map[string]interface{}{
		"Type": "CopyBothResponse",
		"Payload": map[string]interface{}{
			"Format":        byte(c.Format),
			"ColumnFormats": c.ColumnFormats,
		},
	}
}

func (c *CopyBothResponse) WriteTo(w io.Writer) (int64, error) { return writeTo(c, w) }
func (c *CopyBothResponse) String() string                     { return messageToString(c) }