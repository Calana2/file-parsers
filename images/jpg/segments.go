package jpg


var SOISegment = [2]byte{0xff,0xd8}
var EOISegment = [2]byte{0xff,0xd9}

type Segment struct {
  Marker [2]byte
  Length uint16
  Data []byte
}

type APPSegment struct {
  // Application 
  // ffeN
  Marker [2]byte
  Length uint16
  Identifier string
  // JFIF
  Version string
  DensityUnits uint8
  Xdensity uint16
  Ydensity uint16
  XThumbnail uint8
  YThumbnail uint8
  ThumbnailData []byte
}

type COMSegment struct {
  // Comment
  // fffe
  Marker [2]byte
  Length uint16
  Data string
}

type DQTSegment struct {
 // Define Quantization Table
 // ffdb
 Marker [2]byte
 Length uint16
 Destination uint8
 Data []byte
}

type SOFSegment struct {
 // Start Of DCT Frame
 // ffcN
 Marker [2]byte
 Length uint16
 Precision uint8
 LineNB uint16
 Samples_line uint16
 Components uint8
 SOFComponents []SOFComponentInfo
}

type DHTSegment struct {
// Define Huffman Table
// ffc4
 Marker [2]byte
 Length uint16
 Class__Idx uint8
 Bit_Codes [16]uint8
 Real_Huffman_Codes []byte
}

type SOSSegment struct {
 // Start Of Scan
 // ffda
 Marker [2]byte
 Length uint16
 Components uint8
 YIndex uint8
 Y_AC__DC uint8
 CbIndex uint8
 Cb_AC__DC uint8
 CrIndex uint8
 Cr_AC__DC uint8
 SS_Start uint8
 SS_End uint8
 Sucessive_approx uint8
}
