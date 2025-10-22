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

type frame struct {
	arr       datainput.DataInput
	remaining int
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

	stack := []frame{{remaining: arraySize, arr: make(datainput.DataInput, 0, arraySize)}}

	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		if top.remaining == 0 {
			// this array is now complete, we can pop and attach this to parent or return it if possible
			completed := top.arr
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				return completed, nil
			}
			parent := &stack[len(stack)-1]
			parent.arr = append(parent.arr, completed)
			parent.remaining--
			continue
		}

		typeId, err := buffer.ReadByte()
		if err != nil {
			return nil, err
		}

		if typeId == byte(arrayDataTypeId) {
			// nested array
			nestedSize, err := utils.DecodeVarInt(buffer)
			if err != nil {
				return nil, err
			}
			if nestedSize < 0 {
				return nil, errors.New("invalid array size")
			}
			if nestedSize > maxArrayLength {
				return nil, errors.New("array length exceeds maximum limit")
			}
			if nestedSize == 0 {
				// no need to push, just append empty array
				top.arr = append(top.arr, datainput.DataInput{})
				top.remaining--
				continue
			}
			// push new frame for the non empty nested array
			stack = append(stack, frame{remaining: nestedSize, arr: make(datainput.DataInput, 0, nestedSize)})
			continue
		}

		// if not an array recode and append the item
		if proc, ok := d.Df.dataTypeIDProcessors[typeId]; ok {
			item, err := proc.Decode(buffer)
			if err != nil {
				return nil, err
			}
			top.arr = append(top.arr, item)
			top.remaining--
		} else {
			return nil, errors.New("unsupported data type")
		}
	}

	return nil, errors.New("decode error, some error occurred")
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

	if err := buffer.WriteByte(arrayDataTypeId); err != nil {
		return err
	}
	if varIntBuf, err := utils.EncodeVarInt(len(arrayData)); err != nil {
		return err
	} else {
		if _, err := buffer.Write(varIntBuf.Bytes()); err != nil {
			return err
		}
	}

	stack := []frame{{arr: arrayData, remaining: len(arrayData)}}

	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		if top.remaining == 0 {
			stack = stack[:len(stack)-1]
			continue
		}

		item := top.arr[len(top.arr)-top.remaining]
		top.remaining--

		if nested, ok := item.(datainput.DataInput); ok {
			if len(nested) > maxArrayLength {
				return errors.New("array length exceeds maximum limit")
			}
			if err := buffer.WriteByte(byte(arrayDataTypeId)); err != nil {
				return err
			}
			if varIntBuf, err := utils.EncodeVarInt(len(nested)); err != nil {
				return err
			} else {
				if _, err := buffer.Write(varIntBuf.Bytes()); err != nil {
					return err
				}
			}
			if len(nested) == 0 {
				continue
			}
			stack = append(stack, frame{arr: nested, remaining: len(nested)})
			continue
		}

		// not an array, encode normally
		if err := d.Df.Encode(item, buffer); err != nil {
			return err
		}
	}

	return nil
}

func NewArrayDataType(df *DatatypeFactory) DataType {
	return &ArrayDataType{Df: df}
}
