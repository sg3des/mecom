package mecom

import (
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/sg3des/stob"
)

type Value interface {
	io.Reader
	io.Writer
	stob.Size
}

//
//

type UINT16 uint16

func (n UINT16) Read(p []byte) (int, error) {
	return copy(p, []byte(fmt.Sprintf("%04X", uint16(n)))), io.EOF
}

func (n UINT16) Size() int { return 4 }

func (n *UINT16) Write(p []byte) (int, error) {
	size := n.Size()
	if len(p) < size {
		return 0, io.ErrUnexpectedEOF
	}

	_, err := fmt.Sscanf(string(p[:size]), "%04X", n)
	return n.Size(), err
}

//
//

type UINT8 uint8

func (n UINT8) Read(p []byte) (int, error) {
	return copy(p, []byte(fmt.Sprintf("%02X", uint8(n)))), io.EOF
}

func (n UINT8) Size() int { return 2 }

func (n *UINT8) Write(p []byte) (int, error) {
	size := n.Size()
	if len(p) < size {
		return 0, io.ErrUnexpectedEOF
	}

	_, err := fmt.Sscanf(string(p[:size]), "%02X", n)
	return n.Size(), err
}

//
//

type FLOAT32 float32

func (n FLOAT32) Read(p []byte) (int, error) {
	s := fmt.Sprintf("%08X", math.Float32bits(float32(n)))
	return copy(p, []byte(s)), io.EOF
}

func (n FLOAT32) Size() int { return 8 }

func (n *FLOAT32) Write(p []byte) (int, error) {
	size := n.Size()
	if len(p) < size {
		return 0, io.ErrUnexpectedEOF
	}

	var u uint32
	_, err := fmt.Sscanf(string(p[:size]), "%08X", &u)
	if err != nil {
		return 0, err
	}

	*n = FLOAT32(math.Float32frombits(u))

	return size, err
}

//
//

type STR string

func (s STR) Read(p []byte) (int, error) {
	return copy(p, []byte(s)), io.EOF
}

func (s STR) Size() int {
	return len(s)
}

func (s *STR) Write(p []byte) (int, error) {
	if string(p[:2]) == CommandVS {
		*s = CommandVS
		return 2, nil
	}

	if string(p[:3]) == CommandVR {
		*s = CommandVR
		return 3, nil
	}

	return 0, errors.New("unexpected string")
}

//
//
//

// type RespValue []byte

// func (v RespValue) Read(p []byte) (int, error) {
// 	return copy(p, []byte(fmt.Sprintf("%04X", v))), io.EOF
// }

// func (n RespValue) Size() int { return 0 }

// func (n *RespValue) Write(p []byte) (int, error) {
// 	size := n.Size()
// 	if len(p) < size {
// 		return 0, io.ErrUnexpectedEOF
// 	}

// 	*n = make([]byte, 4)

// 	_, err := hex.Decode(*n, p[:size])

// 	return size, err
// }

// func (n RespValue) Float32() float32 {
// 	return math.Float32frombits(binary.BigEndian.Uint32(n))
// }

// func (n RespValue) Uint16() uint16 {
// 	return uint16(binary.BigEndian.Uint32(n))
// }
