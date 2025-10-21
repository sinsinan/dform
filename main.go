package main

import (
	"fmt"
	"log"

	"github.com/sinsinan/dform/compression"
	"github.com/sinsinan/dform/datainput"
	"github.com/sinsinan/dform/dform"
)

func main() {
	data := datainput.DataInput{"sinan", int32(3), datainput.DataInput{"Kar", "edat"}}
	log.Printf("original data %+v", data)
	s, err := dform.EncodeDFormWithCompression(data, compression.CompressionTypeLZ4)
	if err != nil {
		panic(err)
	}

	log.Printf("data %s", s)
	// string builder for the hex output
	sb := ""
	for i := 0; i < len(s); i++ {
		sb += fmt.Sprintf("%02x ", s[i])
	}
	log.Printf("hex data %s", sb)

	decodedData, err := dform.DecodeDForm(s)
	if err != nil {
		panic(err)
	}

	log.Printf("decoded data %+v", decodedData)
}
