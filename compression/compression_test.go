package compression

import (
	"testing"
)

func TestNoCompressionRoundTrip(t *testing.T) {
	cf := NewCompressionFactory()
	if err := cf.Register(NewNoCompession()); err != nil {
		t.Fatalf("register no compression failed: %v", err)
	}
	data := []byte("hello world")
	compressed, err := cf.Compress(data, CompressionTypeNone)
	if err != nil {
		t.Fatalf("compress failed: %v", err)
	}
	if string(compressed) != string(data) {
		t.Fatalf("no compression changed data")
	}
	decompressed, err := cf.DeCompress(compressed, CompressionTypeNone)
	if err != nil {
		t.Fatalf("decompress failed: %v", err)
	}
	if string(decompressed) != string(data) {
		t.Fatalf("decompress mismatch")
	}
}

func TestLZ4CompressionRegistered(t *testing.T) {
	cf := NewCompressionFactory()
	if err := cf.Register(NewLZ4Compession()); err != nil {
		t.Fatalf("register lz4 failed: %v", err)
	}
	if !cf.IsCompressionTypePresent(CompressionTypeLZ4) {
		t.Fatalf("lz4 not present after register")
	}
}

func TestRegisterDuplicateCompression(t *testing.T) {
	cf := NewCompressionFactory()
	if err := cf.Register(NewNoCompession()); err != nil {
		t.Fatalf("initial register failed: %v", err)
	}
	// second register should return an error
	if err := cf.Register(NewNoCompession()); err == nil {
		t.Fatalf("expected error when registering duplicate compression type")
	}
}
