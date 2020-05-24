package mecom

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/sg3des/stob"
)

func TestUINT16(t *testing.T) {
	var x UINT16 = 10

	_, err := io.Copy(&x, bytes.NewReader([]byte("000F")))
	if err != nil && err != io.EOF {
		t.Error(err)
	}

	if x != 15 {
		t.Error("not equal", x)
	}
}

func TestUINT8(t *testing.T) {
	var x UINT8 = 10

	_, err := io.Copy(&x, bytes.NewReader([]byte("0F")))
	if err != nil && err != io.EOF {
		t.Error(err)
	}

	if x != 15 {
		t.Error("not equal", x)
	}
}

func TestFLOAT32(t *testing.T) {
	var x FLOAT32 = 37.5

	_, err := io.Copy(&x, bytes.NewReader([]byte("42160000")))
	if err != nil && err != io.EOF {
		t.Error(err)
	}

	if x != 37.5 {
		t.Error("not equal", x)
	}
}

func TestSmartValuesStruct(t *testing.T) {
	var data struct {
		Address UINT16
	}

	s, err := stob.NewStruct(&data)
	if err != nil {
		t.Error(err)
	}

	if _, err := io.Copy(s, bytes.NewReader([]byte("000A"))); err != nil {
		t.Error(err)
	}

	b, err := ioutil.ReadAll(s)
	if err != nil {
		t.Error(err)
	}

	if string(b) != "000A" {
		t.Error("not equal", string(b))
	}

	t.Log(string(b), data)
}

func TestRespValue(t *testing.T) {
	var x RespValue

	_, err := io.Copy(&x, bytes.NewReader([]byte("42160000")))
	if err != nil {
		t.Error(err)
	}

	if x.Float32() != 37.5 {
		t.Error("not equal", x.Float32())
	}

	_, err = io.Copy(&x, bytes.NewReader([]byte("00000001")))
	if err != nil {
		t.Error(err)
	}

	if x.Uint16() != 1 {
		t.Error("not equal", x.Uint16())
	}
}
