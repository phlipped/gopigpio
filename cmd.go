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

func sendCmd(p io.ReadWriter, c Cmd) (int32, error) {
	cmdAsBytes := c.encodeToBytes()
	if _, err := p.Write(cmdAsBytes); err != nil {
		return 0, err
	}

	// Read the first 32 bytes back
	responseBytes := make([]byte, 16)
	if _, err := io.ReadFull(p, responseBytes); err != nil {
		return 0, err
	}

	// Verify the first 12 bytes of the response match the first 12 bytes of the request
	if !bytes.Equal(cmdAsBytes[0:12], responseBytes[0:12]) {
		return 0, fmt.Errorf("Unexpected response header: want %v, got %v", cmdAsBytes[0:12], responseBytes[0:12])
	}

	res := binary.LittleEndian.Uint32(responseBytes[12:16])
	return int32(res), nil
}

// WARNING: Callers must check if res is negative, as well as checking for err, before assuming there is
// data in ext.
func sendCmdExtResponse(p io.ReadWriter, c Cmd) (res int32, ext []byte, err error) {
	res, err = sendCmd(p, c)
	if err != nil {
		return res, nil, err
	}

	// <res> indicates the length of the extended response data that follows. It should generally
	// be interpreted as a uint32
	// Unfortunately, <res> can also contain an error code when an error has occurred. Error values
	// are negative int32 values.
	// Thus, it is not possible to differentiate between a large uint32 value (indicating the length of data that follows), vs a negative int32 value (indicating an error)

	// For the moment, we will assume that if <res> is negative when interpreted as an int32, then
	// an error has indeed occurred, and therefore there is no extended data in the response waiting
	// to be read, and we should just return <res> to the caller.
	// The rationale for this is based on a few observations:
	//   1) Very large extended responses are unlikely, and so we are catering to the most common
	//      case
	//   2) It's a safer strategy than assuming a large uint32 and trying to read the corresponding
	//      data from <p>, which would block.
	//   3) Users who actually find they are actually getting bitten by this issue can still work
	//      around it by reading the data from <p> themselves.

	if res < 0 {
		return res, nil, nil
	}

	ext = make([]byte, res)
	_, err = io.ReadFull(p, ext)
	return res, ext, err
}
