package utils

import (
	"bytes"
	"testing"
)

func TestVarIntRoundTrip(t *testing.T) {
	cases := []int{0, 1, 127, 128, 129, 255, 256, 1024, 1 << 20}
	for _, c := range cases {
		buf, err := EncodeVarInt(c)
		if err != nil {
			t.Fatalf("EncodeVarInt(%d) error: %v", c, err)
		}
		b := bytes.NewBuffer(buf.Bytes())
		got, err := DecodeVarInt(b)
		if err != nil {
			t.Fatalf("DecodeVarInt for %d error: %v", c, err)
		}
		if got != c {
			t.Fatalf("roundtrip mismatch: want %d got %d", c, got)
		}
	}
}

func TestVarUInt32RoundTrip(t *testing.T) {
	cases := []uint32{0, 1, 127, 128, 129, 255, 256, 1024, 1 << 20}
	for _, c := range cases {
		buf, err := EncodeVarUInt32(c)
		if err != nil {
			t.Fatalf("EncodeVarUInt32(%d) error: %v", c, err)
		}
		b := bytes.NewBuffer(buf.Bytes())
		got, err := DecodeVarUInt32(b)
		if err != nil {
			t.Fatalf("DecodeVarUInt32 for %d error: %v", c, err)
		}
		if got != c {
			t.Fatalf("roundtrip mismatch: want %d got %d", c, got)
		}
	}
}

func TestZigZagSymmetry(t *testing.T) {
	values := []int32{0, 1, -1, 127, -127, 128, -128, 1<<15 - 1, -(1 << 15)}
	for _, v := range values {
		enc := EncodeZigZagInt32(v)
		dec := DecodeZigZagInt32(enc)
		if dec != v {
			t.Fatalf("zigzag symmetry failed for %d: got %d", v, dec)
		}
	}
}
