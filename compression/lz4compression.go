package compression

import "github.com/gocql/gocql/lz4"

const CompressionTypeLZ4 CompressionType = 0x01

type LZ4Compession struct {
	lz4 lz4.LZ4Compressor
}

// Compress implements Compression.
func (l *LZ4Compession) Compress(data []byte) ([]byte, error) {
	return l.lz4.Encode(data)
}

// Decompress implements Compression.
func (l *LZ4Compession) Decompress(data []byte) ([]byte, error) {
	return l.lz4.Decode(data)
}

// GetCompressionType implements Compression.
func (l *LZ4Compession) GetCompressionType() CompressionType {
	return CompressionTypeLZ4
}

func NewLZ4Compession() Compression {
	return &LZ4Compession{
		lz4: lz4.LZ4Compressor{},
	}
}
