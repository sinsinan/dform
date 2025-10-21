package dform

import (
	"bytes"
	"errors"

	"github.com/sinsinan/dform/compression"
	"github.com/sinsinan/dform/datainput"
	"github.com/sinsinan/dform/datatype"
)

var df = datatype.NewDatatypeFactory()
var cf = compression.NewCompressionFactory()

// version takes 1 byte for now
const dfromVersion byte = 0x01

func init() {
	if err := df.Register(datatype.NewStringDataType()); err != nil {
		panic(err)
	}
	if err := df.Register(datatype.NewInt32DataType()); err != nil {
		panic(err)
	}
	if err := df.Register(datatype.NewArrayDataType(df)); err != nil {
		panic(err)
	}
	if err := cf.Register(compression.NewNoCompession()); err != nil {
		panic(err)
	}

	if err := cf.Register(compression.NewLZ4Compession()); err != nil {
		panic(err)
	}
}

func EncodeDForm(toSend datainput.DataInput) (string, error) {
	return EncodeDFormWithCompression(toSend, compression.CompressionTypeLZ4)
}

func EncodeDFormWithCompression(toSend datainput.DataInput, compressionType compression.CompressionType) (string, error) {
	if !cf.IsCompressionTypePresent(compressionType) {
		return "", errors.New("compression type not found")
	}
	var buf bytes.Buffer
	if err := df.Encode(toSend, &buf); err != nil {
		return "", err
	}
	data := buf.Bytes()
	compressedData, err := cf.Compress(data, compressionType)
	if err != nil {
		return "", err
	}

	var outputBuf bytes.Buffer
	outputBuf.WriteByte(dfromVersion)
	outputBuf.WriteByte(byte(compressionType))
	outputBuf.Write(compressedData)

	return outputBuf.String(), nil

}

func DecodeDForm(data string) (datainput.DataInput, error) {
	buf := bytes.NewBufferString(data)
	version, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}
	if version != dfromVersion {
		return nil, errors.New("version missmatch")
	}

	comTypeByte, err := buf.ReadByte()
	comType := compression.CompressionType(comTypeByte)
	if err != nil {
		return nil, err
	}
	if !cf.IsCompressionTypePresent(comType) {
		return nil, errors.New("compression type not found")
	}

	compressedData := buf.Bytes()
	dataBytes, err := cf.DeCompress(compressedData, comType)
	if err != nil {
		return nil, err
	}

	outputBuf := bytes.NewBuffer(dataBytes)
	dii, err := df.Decode(outputBuf)
	if err != nil {
		return nil, err
	}

	di, ok := dii.(datainput.DataInput)
	if !ok {
		return nil, errors.New("could not cast output to DataInput")
	}
	return di, nil
}
