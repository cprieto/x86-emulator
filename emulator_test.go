package x86_emulator

import (
	"testing"
	"bytes"
	"io"
	)

func rawHeader() []byte {
	// 32 bytes
	return []byte{
		0x4d, 0x5a, 0x2b, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x01, 0x01, 0xff, 0xff, 0x01, 0x00,
		0x00, 0x10, 0x00, 0x00, 0x00, 0x01, 0x02, 0x00, 0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
}

func TestExample1(t *testing.T) {
	t.Log("example1")
}

func TestParseHeaderSignature(t *testing.T) {
	var reader io.Reader = bytes.NewReader(rawHeader())
	actual := ParseHeader(reader)
	expected := [2]byte{'M', 'Z'}
	if actual.exSignature != expected {
		t.Errorf("expected %v but actual %v", expected, actual.exSignature)
	}
}

func TestParseHeaderSize(t *testing.T) {
	var reader io.Reader = bytes.NewReader(rawHeader())
	actual := ParseHeader(reader)
	expected := word(2)
	if actual.exHeaderSize != expected {
		t.Errorf("expected %v but actual %v", expected, actual.exHeaderSize)
	}
}

func TestParseHeaderInitSS(t *testing.T) {
	var reader io.Reader = bytes.NewReader(rawHeader())
	actual := ParseHeader(reader)
	expected := word(0x0001)
	if actual.exInitSS != expected {
		t.Errorf("expected %v but actual %v", expected, actual.exInitSS)
	}
}

func TestParseHeaderInitSP(t *testing.T) {
	var reader io.Reader = bytes.NewReader(rawHeader())
	actual := ParseHeader(reader)
	expected := word(0x1000)
	if actual.exInitSP != expected {
		t.Errorf("expected %v but actual %v", expected, actual.exInitSP)
	}
}

func TestParseHeaderInitIP(t *testing.T) {
	var reader io.Reader = bytes.NewReader(rawHeader())
	actual := ParseHeader(reader)
	expected := word(0x0100)
	if actual.exInitIP != expected {
		t.Errorf("expected %v but actual %v", expected, actual.exInitIP)
	}
}

func TestParseHeaderInitCS(t *testing.T) {
	var reader io.Reader = bytes.NewReader(rawHeader())
	actual := ParseHeader(reader)
	expected := word(0x0002)
	if actual.exInitCS != expected {
		t.Errorf("expected %v but actual %v", expected, actual.exInitCS)
	}
}
