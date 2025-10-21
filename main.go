package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/sinsinan/dform/compression"
	"github.com/sinsinan/dform/datainput"
	"github.com/sinsinan/dform/dform"
)

func main() {
	datas := []datainput.DataInput{
		{"sinan", int32(3), datainput.DataInput{"Kar", "edat"}},
		{},
		{""},
		{int32(0)},
		{int32(1)},
		{datainput.DataInput{}},
		{datainput.DataInput{"a", int32(-1), datainput.DataInput{"b", int32(2)}}},
		{datainput.DataInput{"nested", datainput.DataInput{"array", datainput.DataInput{"string1", "string12", int32(42)}}}},
	}
	for i, data := range datas {
		log.Printf("%d - original data %#v", i, data)
		s, err := dform.EncodeDFormWithCompression(data, compression.CompressionTypeLZ4)
		if err != nil {
			panic(err)
		}

		sb := ""
		for i := 0; i < len(s); i++ {
			sb += fmt.Sprintf("%02x ", s[i])
		}
		log.Printf("%d - hex data %s", i, sb)

		decodedData, err := dform.DecodeDForm(s)
		if err != nil {
			panic(err)
		}

		log.Printf("%d - decoded data %#v", i, decodedData)
	}
}
