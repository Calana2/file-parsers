package utils

import "encoding/binary"

func ExtractUint16(data []byte, endianness binary.ByteOrder) uint16 {
 if endianness == binary.LittleEndian {
  return binary.LittleEndian.Uint16(data)
 } else {
  return binary.BigEndian.Uint16(data)
 }
}

func ExtractUint32(data []byte, endianness binary.ByteOrder) uint32 {
 if endianness == binary.LittleEndian {
  return binary.LittleEndian.Uint32(data)
 } else {
  return binary.BigEndian.Uint32(data)
 }
}
