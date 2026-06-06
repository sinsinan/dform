# DForm

A custom, high-performance binary serialization format in \textbf{Go} utilizing LEB128/ZigZag integer encoding and LZ4 compression to optimize network payload footprint and data pipeline I/O throughput.


## Project setup

1. install go (1.20+)
2. unzip dform-main.zip - `unzip dform-main.zip`
3. navigate to dform directory - `cd dform-main`
4. install dependencies - `go mod tidy`
5. run tests - `go test ./...`

## How to use

Minimal example (encode then decode):

```go
package main

import (
  "fmt"

  "github.com/sinsinan/dform/compression"
  "github.com/sinsinan/dform/datainput"
  "github.com/sinsinan/dform/dform"
)

func main() {
  data := datainput.DataInput{"hello", int32(42)}
  // encode with LZ4
  payload, err := dform.EncodeDFormWithCompression(data, compression.CompressionTypeLZ4)
  if err != nil {
    panic(err)
  }

  // decode
  out, err := dform.DecodeDForm(payload)
  if err != nil {
    panic(err)
  }
  fmt.Printf("decoded: %#v\n", out)
}
```


## Protocol definition
```txt
|--------------------------------|
|         Header (2 bytes)       |
|--------------------------------|
|        Compressed Data         |
|--------------------------------|
```
## Header
```txt
|-------------------------------------------|
|  version - 1 byte  | compression - 1 byte |
|-------------------------------------------|
```

  * ### Version
    Version is very important when building wirelevel protocols as we might not be able to make process data encoded with a newer version to be decoded by and older version, It is benificial to give a version field, using which we can decide whether we can decode or not.

  * ### Compression
    It is good to have a support for compression in a wirelevel protocol as most often network IO is bottleneck for most processing pipelines not cpu.

  * 0x00 - none
  * 0x01 - LZ4

## Data
```txt
|--------------------------------|
| Array Length - varint | items  |
|--------------------------------|
```
### datatypes
```txt
|-----------------------------------------|
| datatype_id - 1 byte | items data | ... |
|-----------------------------------------|
```

supported datatypes:
* 0x01 - array
* 0x02 - string
* 0x03 - int(32)

1. Array
   1. array will be respresented as same as `Data`
2. string
   1. `| size - varint | string as bytes |`
3. int32
   1. `| signed varint |`

### VarInt
VarInt uses LEB128 for variable encoding.

### Signed VarInt
Uses ZigZag + LEB128 for signed int encoding


## Performance and Complexity

- Encoding:
  - Space complexity: O(n + D)
    - n - total number of items and bytes in the encoded data
    - D - maximum nesting depth of arrays
  - Time complexity: O(n)
    - n - total number of items and bytes in the encoded data
- Decoding:
  - Space complexity: O(n + D)
    - n - total number of items and bytes in the decoded data
    - D - maximum nesting depth of arrays
  - Time complexity: O(n)
    - n - total number of items and bytes in the decoded data

## How to add new datatype/compression

Add a datatype:

1. Implement `datatype.DataType` in `datatype/` and provide `NewMyType() DataType`.
2. Register it with a `DatatypeFactory` before use:

```go
df := datatype.NewDatatypeFactory()
df.Register(datatype.NewStringDataType())
df.Register(datatype.NewInt32DataType())
df.Register(datatype.NewArrayDataType(df))
df.Register(datatype.NewMyType())
```

Add a compressor:

1. Implement `compression.Compression` in `compression/` and provide `NewMyCompression()`.
2. Register it with a `CompressionFactory` before use:

```go
cf := compression.NewCompressionFactory()
cf.Register(compression.NewNoCompession())
cf.Register(compression.NewLZ4Compession())
cf.Register(compression.NewMyCompression())
```

Notes:

- Pick unique 1-byte ids for types/compression (don't collide with built-ins).
- Test: write encode/decode tests for your type/compression.
