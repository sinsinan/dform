package compression

const CompressionTypeNone CompressionType = 0x00

type NoCompession struct{}

// Compress implements Compression.
func (n *NoCompession) Compress(data []byte) ([]byte, error) {
	return data, nil
}

// Decompress implements Compression.
func (n *NoCompession) Decompress(data []byte) ([]byte, error) {
	return data, nil
}

// GetCompressionType implements Compression.
func (n *NoCompession) GetCompressionType() CompressionType {
	return CompressionTypeNone
}

func NewNoCompession() Compression {
	return &NoCompession{}
}
