# DForm

## Protocol definition

| header | data |

## Header
| version - 1 byte | compression - 1 byte | data length - varint |

### Version
Version is very important when building wirelevel protocols as we might not be able to make process data encoded with a newer version to be decoded by and older version, It is benificial to give a version field, using which we can decide whether we can decode or not.

### Compression
It is good to have a support for compression in a wirelevel protocol as most often network IO is bottleneck for most processing pipelines not cpu.

* 0x00 - none
* 0x01 - ZSTD

### Data
| Array Length - varint | items |

#### items
| datatype_id - 1 byte | items data | ...

datatype_ids:
* 0x01 - array
* 0x02 - string
* 0x03 - int(32)

1. Array
   1. array will be respresented as same as `Data`
2. string
   1. | size - varint | data |
3. int32
   1. | signed varint |

### VarInt
VarInt uses LEB128 for variable encoding.

### Signed VarInt
Uses ZigZag + LEB128 for signed int encoding