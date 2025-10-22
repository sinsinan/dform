package dform

import (
	"reflect"
	"testing"

	"github.com/sinsinan/dform/compression"
	"github.com/sinsinan/dform/datainput"
)

func TestDFormRoundTrip_NoCompression(t *testing.T) {
	in := datainput.DataInput{"a", int32(1), datainput.DataInput{"b"}}
	outStr, err := EncodeDFormWithCompression(in, compression.CompressionTypeNone)
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}
	out, err := DecodeDForm(outStr)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if !equalDataInput(in, out) {
		t.Fatalf("mismatch: want %#v got %#v", in, out)
	}
}

func TestDFormRoundTrip_LZ4(t *testing.T) {
	in := datainput.DataInput{"x", int32(-42)}
	outStr, err := EncodeDFormWithCompression(in, compression.CompressionTypeLZ4)
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}
	out, err := DecodeDForm(outStr)
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if !equalDataInput(in, out) {
		t.Fatalf("mismatch: want %#v got %#v", in, out)
	}
}

// helper uses simple deep-equality for datainput
func equalDataInput(a, b datainput.DataInput) bool {
	if len(a) != len(b) {
		return false
	}
	reflectA := reflect.ValueOf(a)
	reflectB := reflect.ValueOf(b)
	return reflect.DeepEqual(reflectA.Interface(), reflectB.Interface())
}

func TestDecodeVersionMismatch(t *testing.T) {
	in := datainput.DataInput{"v"}
	outStr, err := EncodeDFormWithCompression(in, compression.CompressionTypeNone)
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}
	// tamper version byte
	b := []byte(outStr)
	b[0] = 0x99
	_, err = DecodeDForm(string(b))
	if err == nil {
		t.Fatalf("expected version mismatch error")
	}
}

func TestDecodeInvalidCompressionType(t *testing.T) {
	in := datainput.DataInput{"c"}
	outStr, err := EncodeDFormWithCompression(in, compression.CompressionTypeNone)
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}
	b := []byte(outStr)
	// set compression type byte to an unknown value
	b[1] = 0x99
	_, err = DecodeDForm(string(b))
	if err == nil {
		t.Fatalf("expected compression type not found error")
	}
}

func TestDecodeTruncatedPayload(t *testing.T) {
	in := datainput.DataInput{"t"}
	outStr, err := EncodeDFormWithCompression(in, compression.CompressionTypeNone)
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}
	b := []byte(outStr)
	if len(b) < 5 {
		t.Skip("encoded payload too small for truncation test")
	}
	truncated := string(b[:len(b)/2])
	_, err = DecodeDForm(truncated)
	if err == nil {
		t.Fatalf("expected decode error for truncated payload")
	}
}
