package datatype

import (
	"bytes"
	"errors"
	"reflect"
	"unicode/utf8"

	"github.com/sinsinan/dform/utils"
)

const stringDataTypeId = 0x01
const maxStringLength = 1000000

type StringDataType struct {
}

// decode implements dataType.
func (d StringDataType) Decode(buffer *bytes.Buffer) (interface{}, error) {
	length, err := utils.DecodeVarInt(buffer)
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, errors.New("invalid string length")
	}
	if length > maxStringLength {
		return nil, errors.New("string length limit exceeded")
	}

	strBytes := make([]byte, length)
	n, err := buffer.Read(strBytes)
	if err != nil {
		return nil, err
	}
	if n != length {
		return nil, errors.New("incomplete string data")
	}

	if !utf8.Valid(strBytes) {
		return nil, errors.New("non utf 8 strings found")
	}

	return string(strBytes), nil

}

// getType implements dataType.
func (d StringDataType) GetType() reflect.Type {
	return reflect.TypeOf("")
}

// getTypeId implements dataType.
func (d StringDataType) GetTypeId() byte {
	return stringDataTypeId
}

func (d StringDataType) Encode(buffer *bytes.Buffer, data interface{}) error {
	sdata, ok := data.(string)
	if !ok {
		return errors.New("invalid data type for stringDataType")
	}
	if len(sdata) > maxStringLength {
		return errors.New("string length limit exceeded")
	}

	if !utf8.ValidString(sdata) {
		return errors.New("non utf 8 strings found")
	}

	if err := buffer.WriteByte(byte(stringDataTypeId)); err != nil {
		return err
	}

	if varint, err := utils.EncodeVarInt(len(sdata)); err == nil {
		_, err := buffer.Write(varint.Bytes())
		if err != nil {
			return err
		}
	} else {
		return err
	}

	_, err := buffer.WriteString(sdata)
	return err
}

func NewStringDataType() DataType {
	return &StringDataType{}
}
