package gopigpio

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"
)

const (
	WAVE_CLEAR           CmdID = 27
	WAVE_TX_BUSY         CmdID = 32
	WAVE_ADD_NEW         CmdID = 53
	WAVE_ADD_GENERIC     CmdID = 28
	WAVE_CREATE          CmdID = 49
	WAVE_TRANSMIT        CmdID = 51
	WAVE_TRANSMIT_REPEAT CmdID = 52
	WAVE_CHAIN           CmdID = 93
	WAVE_TRANSMIT_MODE   CmdID = 100
	WAVE_TRANSMIT_AT     CmdID = 101
)

type Pulse struct {
	OnPins  []Pin
	OffPins []Pin
	Delay   time.Duration
}

func (p Pulse) encodeToBytes() []byte {
	bytes := make([]byte, 12)

	var onMask uint32
	for _, pin := range p.OnPins {
		onMask |= 1 << pin
	}
	binary.LittleEndian.PutUint32(bytes[0:4], onMask)

	var offMask uint32
	for _, pin := range p.OffPins {
		offMask |= 1 << pin
	}
	binary.LittleEndian.PutUint32(bytes[4:8], offMask)

	delay := uint32(p.Delay.Nanoseconds() / 1000)
	binary.LittleEndian.PutUint32(bytes[8:12], delay)

	return bytes
}

func WaveClear(p io.ReadWriter) error {
	cmd := Cmd{
		ID: WAVE_CLEAR,
	}
	res, err := sendCmd(p, cmd)
	if err != nil {
		return err
	}
	if res < 0 {
		return fmt.Errorf("Error while executing WaveClear command: Error Code was %d", res)
	}
	return nil
}

func WaveAddNew(p io.ReadWriter) error {
	cmd := Cmd{
		ID: WAVE_ADD_NEW,
	}
	res, err := sendCmd(p, cmd)
	if err != nil {
		return err
	}
	if res < 0 {
		return fmt.Errorf("Error while executing WaveAddNew command: Error Code was %d", res)
	}
	return nil
}

func WaveAddGeneric(p io.ReadWriter, pulses []Pulse) (int32, error) {
	ext := []byte{}
	for _, pulse := range pulses {
		ext = append(ext, pulse.encodeToBytes()...)
	}
	cmd := Cmd{
		ID:  WAVE_ADD_GENERIC,
		Ext: ext,
	}

	res, err := sendCmd(p, cmd)
	if err != nil {
		return res, err
	}
	if res < 0 {
		return res, fmt.Errorf("Error while executing WaveAddGeneric command: Error Code was %d", res)
	}
	return res, nil
}

type WaveID uint32

func WaveCreate(p io.ReadWriter) (WaveID, error) {
	cmd := Cmd{
		ID: WAVE_CREATE,
	}
	res, err := sendCmd(p, cmd)
	if err != nil {
		return WaveID(res), err
	}
	if res < 0 {
		return WaveID(res), fmt.Errorf("Error while executing WaveCreate command: Error Code was %d", res)
	}
	return WaveID(res), nil
}

func WaveTransmit(p io.ReadWriter, waveID WaveID) (int32, error) {
	cmd := Cmd{
		ID: WAVE_TRANSMIT,
		P1: uint32(waveID),
	}
	res, err := sendCmd(p, cmd)
	if err != nil {
		return res, err
	}
	if res < 0 {
		return res, fmt.Errorf("Error while executing WaveTransmit command: Error Code was %d", res)
	}
	return res, nil
}

func WaveTransmitRepeat(p io.ReadWriter, waveID int32) (int32, error) {
	cmd := Cmd{
		ID: WAVE_TRANSMIT_REPEAT,
		P1: uint32(waveID),
	}
	res, err := sendCmd(p, cmd)
	if err != nil {
		return res, err
	}
	if res < 0 {
		return res, fmt.Errorf("Error while executing WaveTransmitRepeat command: Error Code was %d", res)
	}
	return res, nil
}

type Chainer interface {
	encodeChain() []byte
}

var (
	LOOP_START   = []byte{0xff, 0}
	LOOP_REPEAT  = []byte{0xff, 1, 0, 0}
	LOOP_DELAY   = []byte{0xff, 2, 0, 0}
	LOOP_FOREVER = []byte{0xff, 3}
)

type ChainLoopN struct {
	Chain Chainer
	Count uint16
}

func (cln ChainLoopN) encodeChain() []byte {
	chain := LOOP_START
	chain = append(chain, cln.Chain.encodeChain()...)
	chain = append(chain, LOOP_REPEAT...)
	binary.LittleEndian.PutUint16(chain[len(chain)-2:], cln.Count)
	return chain
}

type ChainLoopForever struct {
	Chain Chainer
}

func (clf ChainLoopForever) encodeChain() []byte {
	chain := LOOP_START
	chain = append(chain, clf.Chain.encodeChain()...)
	chain = append(chain, LOOP_FOREVER...)
	return chain
}

type ChainDelay struct {
	Delay time.Duration
}

func (cd ChainDelay) encodeChain() []byte {
	chain := LOOP_DELAY
	binary.LittleEndian.PutUint16(chain[len(chain)-2:], uint16(cd.Delay.Nanoseconds()/1000))
	return chain
}

type Chainers []Chainer

func (cs Chainers) encodeChain() []byte {
	chain := []byte{}
	for _, c := range cs {
		chain = append(chain, c.encodeChain()...)
	}
	return chain
}

type ChainWaveID uint8

func (cwi ChainWaveID) encodeChain() []byte {
	return []byte{byte(cwi)}
}

func WaveChain(p io.ReadWriter, chain Chainer) (int32, error) {
	cmd := Cmd{
		ID:  WAVE_CHAIN,
		Ext: chain.encodeChain(),
	}
	res, err := sendCmd(p, cmd)
	if err != nil {
		return res, err
	}
	if res < 0 {
		return res, fmt.Errorf("Error while executing WaveChain command: Error Code was %d", res)
	}
	return res, nil

}

type WaveMode uint32

const (
	WAVE_MODE_ONE_SHOT      WaveMode = 0
	WAVE_MODE_REPEAT        WaveMode = 1
	WAVE_MODE_ONE_SHOT_SYNC WaveMode = 2
	WAVE_MODE_REPEAT_SYNC   WaveMode = 3
)

func WaveTransmitMode(p io.ReadWriter, waveID int32, mode WaveMode) (int32, error) {
	cmd := Cmd{
		ID: WAVE_TRANSMIT_MODE,
		P1: uint32(waveID),
		P2: uint32(mode),
	}
	res, err := sendCmd(p, cmd)
	if err != nil {
		return res, err
	}
	if res < 0 {
		return res, fmt.Errorf("Error while executing WaveTransmitMode command: Error Code was %d", res)
	}
	return res, nil
}

func WaveTransmitAt(p io.ReadWriter) (int32, error) {
	cmd := Cmd{
		ID: WAVE_TRANSMIT_AT,
	}
	res, err := sendCmd(p, cmd)
	if err != nil {
		return res, err
	}
	if res < 0 {
		return res, fmt.Errorf("Error while executing WaveTransmitAt command: Error Code was %d", res)
	}
	return res, nil
}

func WaveTxBusy(p io.ReadWriter) (int32, error) {
	cmd := Cmd{
		ID: WAVE_TX_BUSY,
	}
	res, err := sendCmd(p, cmd)
	if err != nil {
		return res, err
	}
	if res < 0 {
		return res, fmt.Errorf("Error while executing WaveTransmitAt command: Error Code was %d", res)
	}
	return res, nil
}
