package mecom

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/apex/log"
	"github.com/sg3des/stob"
	"github.com/snksoft/crc"
	"github.com/tarm/serial"
)

type Controller struct {
	port *serial.Port

	buf     []byte
	seqno   uint16
	address uint16

	sync.Mutex
}

func Dial(devname string) (*Controller, error) {
	port, err := serial.OpenPort(&serial.Config{
		Name: devname,
		Baud: 57600,
	})
	if err != nil {
		return nil, err
	}

	bb := &Controller{
		port: port,
		buf:  make([]byte, 128),
	}

	bb.address, err = bb.LoopStatus()

	return bb, err
}

// execute Command - write and read result
func (bb *Controller) Execute(cmd interface{}) (resp Response, err error) {
	// prepare header
	cmdLine, err := stob.Marshal(&Header{
		Control: ControlHost,
		Address: UINT8(bb.address),
		SeqNo:   UINT16(bb.SeqNo()),
	})
	if err != nil {
		return resp, err
	}

	// convert command structure to the bytes
	cmdPayload, err := stob.Marshal(cmd)
	if err != nil {
		return resp, err
	}

	// append command and calculate crc
	cmdLine = append(cmdLine, cmdPayload...)
	cmdLine = append(cmdLine, bb.CRC(cmdLine)...)
	cmdLine = append(cmdLine, '\r')

	log.Debug(string(cmdLine))

	bb.Lock()
	defer bb.Unlock()

	// write command line
	if _, err := bb.port.Write(cmdLine); err != nil {
		return resp, err
	}

	respLine := bb.readResponse()

	log.Debug(string(respLine))

	err = stob.Unmarshal(respLine, &resp)
	return
}

func (bb *Controller) readResponse() []byte {
	s := bufio.NewScanner(bb.port)
	s.Split(bb.commandsSplitter)

	s.Scan()

	return s.Bytes()
}

func (bb *Controller) commandsSplitter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if i := bytes.Index(data, []byte{EOL}); i > -1 {
		return i + 1, data[:i], nil
	}

	return 0, nil, nil
}

// increment and return next sequence number
func (bb *Controller) SeqNo() uint16 {
	bb.seqno++
	return bb.seqno
}

func (bb *Controller) CRC(data []byte) []byte {
	crc16 := crc.CalculateCRC(crc.XMODEM, data)

	crcdata := make([]byte, 2)
	binary.BigEndian.PutUint16(crcdata, uint16(crc16))

	// return hex.EncodeToString(crcdata)
	return []byte(fmt.Sprintf("%04X", crcdata))
}

//
// Command
//

func (bb *Controller) LoopStatus() (address uint16, err error) {
	resp, err := bb.Execute(&ValueRead{
		Command:     CommandVR,
		ParameterID: 2051,
		Instance:    1,
	})

	return resp.Uint16()
}

func (bb *Controller) GetObjectTemperature() (temp float32, err error) {
	resp, err := bb.Execute(&ValueRead{
		Command:     CommandVR,
		ParameterID: 1000,
		Instance:    1,
	})

	return resp.Float32()
}

func (bb *Controller) GetTargetTemperature() (temp float32, err error) {
	resp, err := bb.Execute(&ValueRead{
		Command:     CommandVR,
		ParameterID: 1010,
		Instance:    1,
	})

	return resp.Float32()
}

func (bb *Controller) SetTemperature(temp float32) error {
	_, err := bb.Execute(&ValueSet{
		Command:     CommandVS,
		ParameterID: 3000,
		Instance:    1,
		Value:       FLOAT32(temp),
	})

	return err
}

func (bb *Controller) SetTECVoltage(v float32) error {
	_, err := bb.Execute(&ValueSet{
		Command:     CommandVS,
		ParameterID: 50002,
		Instance:    1,
		Value:       FLOAT32(v),
	})

	return err
}

func (bb *Controller) SetTECCurrent(v float32) error {
	_, err := bb.Execute(&ValueSet{
		Command:     CommandVS,
		ParameterID: 50001,
		Instance:    1,
		Value:       FLOAT32(v),
	})

	return err
}

//
// Frame
//

const ControlHost = '#'
const ControlDevice = '!'
const EOL = '\r'
const CommandVS = "VS"
const CommandVR = "?VR"

type Header struct {
	Control byte
	Address UINT8
	SeqNo   UINT16
}

type ValueSet struct {
	// Header
	Command     STR
	ParameterID UINT16
	Instance    UINT8
	Value       FLOAT32
	// CRC         UINT16
}

type ValueRead struct {
	// Header
	Command     STR
	ParameterID UINT16
	Instance    UINT8
	// CRC         UINT16
}

type Response struct {
	Header
	Value []byte
	// CRC         UINT16
}

var LenResponseValue = 8
var ErrResponseValue = errors.New("unexpected response value")

func (r Response) Bytes() ([]byte, error) {
	if len(r.Value) < LenResponseValue+4 {
		return nil, ErrResponseValue
	}

	return hex.DecodeString(string(r.Value[:LenResponseValue]))
}

func (r Response) Float32() (float32, error) {
	b, err := r.Bytes()
	if err != nil {
		return 0, err
	}

	return math.Float32frombits(binary.BigEndian.Uint32(b)), nil
}

func (r Response) Uint16() (uint16, error) {
	b, err := r.Bytes()
	if err != nil {
		return 0, err
	}

	return uint16(binary.BigEndian.Uint32(b)), nil
}
