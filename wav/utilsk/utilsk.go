package utilsk

import (
  "fmt"
)

func ConvertSizeUint32(size uint32,toBits ...bool) string {
  _toBits := false
 if len(toBits) > 0 {
  _toBits = toBits[0] 
 }
 const (
  KB = 1024
  MB = 1024 * KB
  GB = 1024 * MB
 )
 switch {
  case _toBits:
   return fmt.Sprintf("%.2f bits",float64(size)*0x8)
  case size >= 800*MB:
   return fmt.Sprintf("%.2f GB", float64(size)/GB)
  case size >= 800*KB:
   return fmt.Sprintf("%.2f MB", float64(size)/MB)
  case size >= 800:
   return fmt.Sprintf("%.2f KB", float64(size)/KB)
  default: 
   return fmt.Sprintf("%d Bytes", size)
 }
}


func ConvertSizeUint16(size uint16, toBits ...bool) string { 
 if len(toBits) > 0 {
  return ConvertSizeUint32(uint32(size),toBits[0]) 
 }
 return ConvertSizeUint32(uint32(size))
}


func CalculateAudioDuration(dataSize uint32, byteRate uint32) float64 { 
  return float64(dataSize) / float64(byteRate) 
} 
