package datatype

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/sinsinan/dform/datainput"
	"github.com/sinsinan/dform/utils"
)

func TestStringDataTypeRoundTrip(t *testing.T) {
	df := NewDatatypeFactory()
	if err := df.Register(NewStringDataType()); err != nil {
		t.Fatalf("register string failed: %v", err)
	}

	buf := bytes.NewBuffer(nil)
	in := "hello world"
	if err := df.Encode(in, buf); err != nil {
		t.Fatalf("encode string failed: %v", err)
	}
	got, err := df.Decode(buf)
	if err != nil {
		t.Fatalf("decode string failed: %v", err)
	}
	if gotStr, ok := got.(string); !ok || gotStr != in {
		t.Fatalf("string mismatch: want %q got %v", in, got)
	}
}

func TestInt32DataTypeRoundTrip(t *testing.T) {
	df := NewDatatypeFactory()
	if err := df.Register(NewInt32DataType()); err != nil {
		t.Fatalf("register int32 failed: %v", err)
	}

	buf := bytes.NewBuffer(nil)
	in := int32(-123456)
	if err := df.Encode(in, buf); err != nil {
		t.Fatalf("encode int32 failed: %v", err)
	}
	got, err := df.Decode(buf)
	if err != nil {
		t.Fatalf("decode int32 failed: %v", err)
	}
	if gotInt, ok := got.(int32); !ok || gotInt != in {
		t.Fatalf("int32 mismatch: want %v got %v", in, got)
	}
}

func TestArrayDataTypeRoundTrip(t *testing.T) {
	df := NewDatatypeFactory()
	// register base types and array type
	if err := df.Register(NewStringDataType()); err != nil {
		t.Fatalf("register string failed: %v", err)
	}
	if err := df.Register(NewInt32DataType()); err != nil {
		t.Fatalf("register int32 failed: %v", err)
	}
	if err := df.Register(NewArrayDataType(df)); err != nil {
		t.Fatalf("register array failed: %v", err)
	}

	in := datainput.DataInput{"one", int32(2), datainput.DataInput{"inner"}}
	buf := bytes.NewBuffer(nil)
	if err := df.Encode(in, buf); err != nil {
		t.Fatalf("encode array failed: %v", err)
	}
	got, err := df.Decode(buf)
	if err != nil {
		t.Fatalf("decode array failed: %v", err)
	}
	if !reflect.DeepEqual(in, got) {
		t.Fatalf("array mismatch: want %#v got %#v", in, got)
	}
}

func TestStringDataTypeTooLong(t *testing.T) {
	df := NewDatatypeFactory()
	if err := df.Register(NewStringDataType()); err != nil {
		t.Fatalf("register string failed: %v", err)
	}
	// craft a buffer with a string length exceeding maxStringLength
	buf := bytes.NewBuffer(nil)
	// write type id
	buf.WriteByte(byte(stringDataTypeId))
	// write varint length > maxStringLength (use a large number)
	if varint, err := utils.EncodeVarInt(maxStringLength + 1); err == nil {
		buf.Write(varint.Bytes())
	} else {
		t.Fatalf("prepare varint failed: %v", err)
	}
	// decoding should error
	_, err := df.Decode(buf)
	if err == nil {
		t.Fatalf("expected error for too long string")
	}
}

func TestArrayDataTypeTooLarge(t *testing.T) {
	df := NewDatatypeFactory()
	if err := df.Register(NewArrayDataType(df)); err != nil {
		t.Fatalf("register array failed: %v", err)
	}
	// craft a buffer with array type id and length > maxArrayLength
	buf := bytes.NewBuffer(nil)
	buf.WriteByte(byte(arrayDataTypeId))
	if varint, err := utils.EncodeVarInt(maxArrayLength + 1); err == nil {
		buf.Write(varint.Bytes())
	} else {
		t.Fatalf("prepare varint failed: %v", err)
	}
	_, err := df.Decode(buf)
	if err == nil {
		t.Fatalf("expected error for too large array")
	}
}

func TestDecodeUnsupportedTypeId(t *testing.T) {
	df := NewDatatypeFactory()
	// put an unknown type id in buffer
	buf := bytes.NewBuffer(nil)
	buf.WriteByte(0xff) // unregistered type id
	_, err := df.Decode(buf)
	if err == nil {
		t.Fatalf("expected error for unsupported type id")
	}
}
