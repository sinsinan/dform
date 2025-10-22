package main

import (
	"fmt"
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/sinsinan/dform/compression"
	"github.com/sinsinan/dform/datainput"
	"github.com/sinsinan/dform/dform"
)

func testCase(t *testing.T, name string, data datainput.DataInput, expectedError bool) {

	// Test both with and without compression
	compressionTypes := []compression.CompressionType{
		compression.CompressionTypeNone,
		compression.CompressionTypeLZ4,
	}

	for _, compressionType := range compressionTypes {

		// Encode
		encoded, err := dform.EncodeDFormWithCompression(data, compressionType)
		if err != nil {
			if expectedError {
				return
			}
			log.Errorf("Unexpected encode error: %v", err)
			t.Fatalf("encode error: %v", err)
		}

		// Print encoded data in hex
		hexOutput := ""
		for i := 0; i < len(encoded); i++ {
			hexOutput += fmt.Sprintf("%02x ", encoded[i])
		}
		// log.Printf("Encoded (hex): %s", hexOutput)

		// Decode
		decoded, err := dform.DecodeDForm(encoded)
		if err != nil {
			if expectedError {
				return
			}
			log.Errorf("Unexpected decode error: %v", err)
			t.Fatalf("decode error: %v", err)
		}

		// Verify
		if !reflect.DeepEqual(data, decoded) {
			log.Errorf("Data mismatch: expected %#v got %#v", data, decoded)
			t.Fatalf("data mismatch: expected %#v got %#v", data, decoded)
		}
	}
}

func TestArrayDataTypeVariousCases(t *testing.T) {
	// Test 1: Basic integer arrays
	testCase(t, "Basic Integer Array",
		datainput.DataInput{int32(-1), int32(-11), int32(-1111111111)},
		false)

	// Test 2: Empty array
	testCase(t, "Empty Array",
		datainput.DataInput{},
		false)

	// Test 3: Array with empty string
	testCase(t, "Array with Empty String",
		datainput.DataInput{""},
		false)

	// Test 4: Nested empty array
	testCase(t, "Nested Empty Array",
		datainput.DataInput{datainput.DataInput{}},
		false)

	// Test 5: Complex nested structure
	testCase(t, "Complex Nested Structure",
		datainput.DataInput{
			"hello",
			int32(42),
			datainput.DataInput{
				"nested",
				int32(-999),
				datainput.DataInput{"deep", "nesting"},
			},
		},
		false)

	// Test 6: Array with mixed types
	testCase(t, "Mixed Types",
		datainput.DataInput{
			"string",
			int32(123),
			datainput.DataInput{"nested"},
			"",
			int32(-456),
		},
		false)

	// Test 7: Large array (should be within limits)
	largeArray := make(datainput.DataInput, 999)
	for i := range largeArray {
		largeArray[i] = int32(i)
	}
	testCase(t, "Large Array (within limits)",
		largeArray,
		false)

	// Test 8: Too large array (should error)
	tooLargeArray := make(datainput.DataInput, 1001)
	testCase(t, "Too Large Array",
		tooLargeArray,
		true)

	// Test 9: Deep nesting
	deepNested := datainput.DataInput{"level1"}
	current := &deepNested
	for i := 2; i <= 5; i++ {
		newLevel := datainput.DataInput{fmt.Sprintf("level%d", i)}
		*current = append(*current, newLevel)
		current = &newLevel
	}
	testCase(t, "Deep Nesting",
		deepNested,
		false)
}
