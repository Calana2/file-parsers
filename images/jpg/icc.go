package jpg

type ICCSegment struct {
  Marker [2]byte
  Length uint16
  Identifier string
  Data []byte
}
