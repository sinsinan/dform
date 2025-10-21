package compression

import "errors"

type CompressionFactory struct {
	idCompressionMap map[CompressionType]Compression
}

func (c *CompressionFactory) DeCompress(compressedData []byte, comType CompressionType) ([]byte, error) {
	return c.idCompressionMap[comType].Decompress(compressedData)
}

func (c *CompressionFactory) Compress(data []byte, compressionType CompressionType) ([]byte, error) {
	return c.idCompressionMap[compressionType].Compress(data)
}

func (c *CompressionFactory) Register(compessionToRegister Compression) error {
	if _, ok := c.idCompressionMap[compessionToRegister.GetCompressionType()]; ok {
		return errors.New("Compression with same ID already registered")
	}
	c.idCompressionMap[compessionToRegister.GetCompressionType()] = compessionToRegister
	return nil
}

func (c *CompressionFactory) IsCompressionTypePresent(comType CompressionType) bool {
	_, ok := c.idCompressionMap[comType]
	return ok
}

func NewCompressionFactory() *CompressionFactory {
	return &CompressionFactory{
		idCompressionMap: map[CompressionType]Compression{},
	}
}
