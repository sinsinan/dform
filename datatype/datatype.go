package datatype

import (
	"bytes"
	"reflect"
)

type DataType interface {
	Encode(*bytes.Buffer, interface{}) error
	Decode(*bytes.Buffer) (interface{}, error)
	GetTypeId() byte
	GetType() reflect.Type
}
