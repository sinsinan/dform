package utils

import "bytes"

const lsb7maskuint32 = uint32(0b01111111)
const lsb7maskint = 0b01111111
const msb7maskint = 0b10000000

func EncodeVarUInt32(data uint32) (bytes.Buffer, error) {
	buffer := bytes.Buffer{}
	for data >= 0x80 {
		lsb7 := byte(data & lsb7maskuint32)
		data = data >> 7
		if err := buffer.WriteByte(lsb7 | msb7maskint); err != nil {
			return buffer, err
		}
	}
	err := buffer.WriteByte(byte(data))
	return buffer, err
}

func EncodeVarInt(data int) (bytes.Buffer, error) {
	buffer := bytes.Buffer{}
	for data >= 0x80 {
		lsb7 := byte(data & lsb7maskint)
		data = data >> 7
		if err := buffer.WriteByte(lsb7 | msb7maskint); err != nil {
			return buffer, err
		}
	}
	err := buffer.WriteByte(byte(data))
	return buffer, err
}

func DecodeVarInt(buffer *bytes.Buffer) (int, error) {
	result := 0
	shift := 0
	for {
		b, err := buffer.ReadByte()
		if err != nil {
			return 0, err
		}
		result = result | (int(b&lsb7maskint) << shift)
		if (b & msb7maskint) == 0 {
			break
		}
		shift += 7
	}
	return result, nil
}

func DecodeVarUInt32(buffer *bytes.Buffer) (uint32, error) {
	result := uint32(0)
	shift := 0
	for {
		b, err := buffer.ReadByte()
		if err != nil {
			return 0, err
		}
		result = result | ((uint32(b) & lsb7maskuint32) << shift)
		if (b & msb7maskint) == 0 {
			break
		}
		shift += 7
	}
	return result, nil
}
