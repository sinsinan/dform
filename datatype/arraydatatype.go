package datatype

import (
	"bytes"
	"errors"
	"reflect"

	"github.com/sinsinan/dform/datainput"
	"github.com/sinsinan/dform/utils"
)

const arrayDataTypeId = 0x03
const maxArrayLength = 1000

type ArrayDataType struct {
	Df *DatatypeFactory
}

// Decode implements DataType.
func (d *ArrayDataType) Decode(buffer *bytes.Buffer) (interface{}, error) {
	arraySize, err := utils.DecodeVarInt(buffer)
	if err != nil {
		return nil, err
	}
	if arraySize < 0 {
		return nil, errors.New("invalid array size")
	}
	if arraySize > maxArrayLength {
		return nil, errors.New("array length exceeds maximum limit")
	}

	arrayData := make(datainput.DataInput, arraySize)
	for i := 0; i < arraySize; i++ {
		item, err := d.Df.Decode(buffer)
		if err != nil {
			return nil, err
		}
		arrayData[i] = item
	}

	return arrayData, nil
}

// GetType implements DataType.
func (d *ArrayDataType) GetType() reflect.Type {
	return reflect.TypeOf(datainput.DataInput{})
}

// GetTypeId implements DataType.
func (d *ArrayDataType) GetTypeId() byte {
	return arrayDataTypeId
}

func (d ArrayDataType) Encode(buffer *bytes.Buffer, data interface{}) error {
	arrayData, ok := data.(datainput.DataInput)
	if !ok {
		return errors.New("invalid data type for arrayDataType")
	}

	if len(arrayData) > maxArrayLength {
		return errors.New("array length exceeds maximum limit")
	}

	if err := buffer.WriteByte(byte(arrayDataTypeId)); err != nil {
		return err
	}

	if varIntBuf, err := utils.EncodeVarInt(len(arrayData)); err != nil {
		return err
	} else {
		_, err := buffer.Write(varIntBuf.Bytes())
		if err != nil {
			return err
		}
	}

	for _, item := range arrayData {
		if err := d.Df.Encode(item, buffer); err != nil {
			return err
		}
	}

	return nil
}
