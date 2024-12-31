package jpg

type EXIFSegment struct {
 // Application
 Marker [2]byte
 Length uint16
 Identifier string
 TIFFHeader TIFFHeader
}

type TIFFHeader struct {
  Alignment string
  FixedBytes [2]byte
  IFDOffset uint32
}
