package datatype

import (
	"bytes"
	"errors"
	"reflect"

	log "github.com/sirupsen/logrus"
)

type DatatypeFactory struct {
	dataTypeProcessors   map[reflect.Type]DataType
	dataTypeIDProcessors map[byte]DataType
}

func NewDatatypeFactory() *DatatypeFactory {
	return &DatatypeFactory{
		dataTypeProcessors:   make(map[reflect.Type]DataType),
		dataTypeIDProcessors: make(map[byte]DataType),
	}
}

func (f *DatatypeFactory) Register(datatypeToRegister DataType) error {
	if _, ok := f.dataTypeIDProcessors[datatypeToRegister.GetTypeId()]; ok {
		return errors.New("data type id already registered")
	}

	if _, ok := f.dataTypeProcessors[datatypeToRegister.GetType()]; ok {
		return errors.New("data type already registered")
	}

	f.dataTypeIDProcessors[datatypeToRegister.GetTypeId()] = datatypeToRegister
	f.dataTypeProcessors[datatypeToRegister.GetType()] = datatypeToRegister

	return nil
}

func (f *DatatypeFactory) Encode(data interface{}, buffer *bytes.Buffer) error {
	reflectType := reflect.TypeOf(data)

	if datatypeProcessor, ok := f.dataTypeProcessors[reflectType]; ok {
		return datatypeProcessor.Encode(buffer, data)
	}

	log.Errorf("unsupported data type: %v", reflectType)

	return errors.New("unsupported data type")
}

func (f *DatatypeFactory) Decode(buffer *bytes.Buffer) (interface{}, error) {
	if dataTypeId, err := buffer.ReadByte(); err == nil {
		if datatypeProcessor, ok := f.dataTypeIDProcessors[dataTypeId]; ok {
			return datatypeProcessor.Decode(buffer)
		}
		log.Errorf("unsupported data type id: %v", dataTypeId)
		return nil, errors.New("unsupported data type")
	} else {
		return nil, err
	}
}
