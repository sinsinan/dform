package datatype

import (
	"bytes"
	"errors"
	"reflect"

	"github.com/sinsinan/dform/utils"
)

const int32DataTypeId = 0x02

type Int32DataType struct {
}

// Decode implements DataType.
func (d *Int32DataType) Decode(buffer *bytes.Buffer) (interface{}, error) {
	uint32Value, err := utils.DecodeVarUInt32(buffer)
	if err != nil {
		return nil, err
	}
	int32Value := utils.DecodeZigZagInt32(uint32Value)
	return int32Value, nil
}

// GetType implements DataType.
func (d *Int32DataType) GetType() reflect.Type {
	return reflect.TypeOf(int32(0))
}

// GetTypeId implements DataType.
func (d *Int32DataType) GetTypeId() byte {
	return int32DataTypeId
}

func (d *Int32DataType) Encode(buffer *bytes.Buffer, data interface{}) error {
	idata, ok := data.(int32)
	if !ok {
		return errors.New("invalid data type for int32DataType")
	}
	if err := buffer.WriteByte(byte(int32DataTypeId)); err != nil {
		return err
	}

	if varIntBuff, err := utils.EncodeVarUInt32(utils.EncodeZigZagInt32(idata)); err != nil {
		return err
	} else {
		_, err := buffer.Write(varIntBuff.Bytes())
		return err
	}
}

func NewInt32DataType() DataType {
	return &Int32DataType{}
}
