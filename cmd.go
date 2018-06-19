package gopigpio

import (
	"fmt"
	"bytes"
	"encoding/binary"
	"io"
)

type Cmd struct {
  ID uint32
  P1 uint32
  P2 uint32
  Ext []byte
  ResponseHasExt bool // Indicates whether the response will have an extended.
}

func (c Cmd) encodeToBytes() []byte {
	buf := make([]byte, 16)
	binary.LittleEndian.PutUint32(buf[0:4], c.ID)
	binary.LittleEndian.PutUint32(buf[4:8], c.P1)
	binary.LittleEndian.PutUint32(buf[8:12], c.P2)
	binary.LittleEndian.PutUint32(buf[12:16], uint32(len(c.Ext)))
	if len(c.Ext) > 0 {
		buf = append(buf, c.Ext...)
	}

	return buf
}

type Result struct {
  res int32
  ext []byte
}


func sendCmd(p io.ReadWriter, c Cmd) (Result, error) {
	cmdAsBytes := c.encodeToBytes()
	if _, err := p.Write(cmdAsBytes); err != nil {
		return Result{}, err
	}

	// Read the first 32 bytes back
	responseBytes := make([]byte, 16)
	if _, err := io.ReadFull(p, responseBytes); err != nil {
		return Result{}, err
	}

	// Verify the first 12 bytes of the response match the first 12 bytes of the request
	if !bytes.Equal(cmdAsBytes[0:3], responseBytes[0:3]) {
		return Result{}, fmt.Errorf("Unexpected response header: want %v, got %v", cmdAsBytes[0:3], responseBytes[0:3])
	}

	result := Result{}
	resVal := binary.LittleEndian.Uint32(responseBytes[12:16])
	// If we are expecting an extended response, then resVal indicates the length of the extra data.
	if c.ResponseHasExt {
		result.ext = make([]byte, resVal)
		if _, err := io.ReadFull(p, result.ext); err != nil {
			return result, err
		}
	} else {
		result.res = int32(resVal)
	}

	return result, nil
}

