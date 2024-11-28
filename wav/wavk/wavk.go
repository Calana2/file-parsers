package wavk

import (
  "fmt"
  "os"
  "encoding/binary"
  "main/utilsk"
)

const WAV_HEADER_SIZE = 44 

type WAV struct {
 // Master RIFF chunk
 FileTypeBlockID [4]byte
 FileSize [4]byte
 FileFormatID [4]byte
 // Chunk describing the data format 
 FormatBlockID [4]byte
 BlockSize [4]byte
 AudioFormat [2]byte
 NumChannels [2]byte
 Frequence [4]byte
 BytePerSec [4]byte 
 BytePerBlock [2]byte
 BitsPerSample [2]byte
 DataBlockID [4]byte
 DataSize [4]byte
 //  Chunk containing the sample data
};


func New(file *os.File) (*WAV, error) {
  header := make([]byte,WAV_HEADER_SIZE)
  // Read file and handle errors
  _,err := file.Read(header)
  if err != nil {
   fmt.Println(err)
   return nil, err
  }
  // Create WAV object
  wav := &WAV{}
  copy(wav.FileTypeBlockID[:],header[0:4])
  copy(wav.FileSize[:], header[4:8])
  copy(wav.FileFormatID[:],header[8:12])
  copy(wav.FormatBlockID[:],header[12:16])
  copy(wav.BlockSize[:],header[16:20])
  copy(wav.AudioFormat[:],header[20:22])
  copy(wav.NumChannels[:],header[22:24])
  copy(wav.Frequence[:],header[24:28])
  copy(wav.BytePerSec[:],header[28:32])
  copy(wav.BytePerBlock[:],header[32:34])
  copy(wav.BitsPerSample[:],header[34:36])
  copy(wav.DataBlockID[:],header[36:40])
  copy(wav.DataSize[:],header[40:44])
  return wav, nil
}


func (w *WAV) PrintMetadata() {
 fmt.Printf("'RIFF' chunk descriptor\n")
 fmt.Printf("-----------------------\n")
 fmt.Printf("File Type Block ID: %s (Resource Interchange File Format)\n", w.FileTypeBlockID) 
 fileSize := binary.LittleEndian.Uint32(w.FileSize[:]) + 0x8 
 fmt.Printf("File Size: %s (%d Bytes)\n",utilsk.ConvertSizeUint32(fileSize),fileSize)
 fmt.Printf("File Format ID: %s\n", w.FileFormatID)


 fmt.Printf("\n\n'fmt' sub-chunk\n")
 fmt.Printf("----------------\n")
 fmt.Printf("Format Block ID: %s\n", w.FormatBlockID)
 blockSize := binary.LittleEndian.Uint32(w.BlockSize[:]) + 0x8
 fmt.Printf("Block Size: %s\n",utilsk.ConvertSizeUint32(blockSize))
 audioFormatInt := binary.LittleEndian.Uint16(w.AudioFormat[:])
 switch audioFormatInt {
  case 0x1:
   fmt.Printf("Audio Format: %d [PCM (Pulse Code Modulation) – Uncompressed format]\n", audioFormatInt)
  case 0x3: 
   fmt.Printf("Audio Format: %d [IEEE Float – Floating Point Audio]\n", audioFormatInt)
  case 0x6:
   fmt.Printf("Audio Format: %d [A-law – A-law compression algorithm]\n", audioFormatInt)
  case 0x7:
   fmt.Printf("Audio Format: %d [µ-law – µ-law compression algorithm]\n", audioFormatInt)
  case 0xFFFE:
   fmt.Printf("Audio Format: %d [µ-law – µ-law compression algorithm]\n", audioFormatInt)
  default:
   fmt.Printf("Audio Format: %d [Corrupted or unknown format]\n", audioFormatInt)
 }
  numChannels := binary.LittleEndian.Uint16(w.NumChannels[:])
  if numChannels == 1 {
   fmt.Printf("Number of Channels: %d [Mono]\n", numChannels)
  } else {
   fmt.Printf("Number of Channels: %d [Stereo]\n", numChannels)
  }
  fmt.Printf("Frequence: %d Hz\n", binary.LittleEndian.Uint16(w.Frequence[:]))
  bytesPerSec := binary.LittleEndian.Uint32(w.BytePerSec[:])
  fmt.Printf("Bytes per second: %s\n",utilsk.ConvertSizeUint32(bytesPerSec))
  bytesPerBlock := binary.LittleEndian.Uint16(w.BytePerBlock[:])
  fmt.Printf("Bytes per block: %s\n", utilsk.ConvertSizeUint16(bytesPerBlock))
  bitsPerSample := binary.LittleEndian.Uint16(w.BitsPerSample[:])
  fmt.Printf("Bits per sample: %d\n", bitsPerSample)


  fmt.Printf("\n\n'data' sub-chunk\n")
  fmt.Printf("----------------\n")
  fmt.Printf("Data Block ID: %s\n", w.DataBlockID)
  dataSize := binary.LittleEndian.Uint32(w.DataSize[:])
  fmt.Printf("Data Size: %s (%d Bytes) \n", utilsk.ConvertSizeUint32(dataSize), dataSize)

  fmt.Printf("\n\nEstimated audio duration: %.2f seconds\n", utilsk.CalculateAudioDuration(dataSize, bytesPerSec))
}


