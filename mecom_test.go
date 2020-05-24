package mecom

import (
	"encoding/binary"
	"testing"

	"github.com/apex/log"
	"github.com/sg3des/stob"
	"github.com/snksoft/crc"
)

func TestCRC(t *testing.T) {
	data := "#020003?VR03E801" // AC05

	crc16 := crc.CalculateCRC(crc.XMODEM, []byte(data))
	log.Debug(crc16)

	crcdata := make([]byte, 2)
	binary.BigEndian.PutUint16(crcdata, uint16(crc16))

	t.Logf("%02x", crcdata)
}

func TestFrameVR(t *testing.T) {
	framedata, err := stob.Marshal(&ValueRead{
		Command:     CommandVR,
		ParameterID: 1000,
		Instance:    1,
	})
	if err != nil {
		t.Error(err)
	}

	t.Logf("%s", framedata)
}
