package utils

func EncodeZigZagInt32(val int32) uint32 {
	return uint32(val<<1) ^ uint32(val>>31)
}

func DecodeZigZagInt32(val uint32) int32 {
	return int32((val >> 1) ^ uint32((int32(val&1) * -1)))
}
