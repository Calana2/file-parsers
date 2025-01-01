package jpg

type EXIFSegment struct {
 // Application
 Marker [2]byte
 Length uint16
 Identifier string
 TIFFHeader TIFFHeader
 IFDs []IFD
}

type TIFFHeader struct {
  Alignment string
  FixedBytes [2]byte
  IFDOffset uint32
}

type IFD struct {
 EntriesNum uint16
 Entries []IFDEntry
 OffsetNextIFD uint32
}

type IFDEntry struct {
 Tag uint16
 Format uint16
 ComponentsNum uint32
 Data_Offset uint32
 Data []interface{}
}
