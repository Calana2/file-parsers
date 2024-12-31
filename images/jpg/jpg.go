package jpg

import "fmt"

type JPG struct {
 Name string
 SOI [2]byte
 Segments []interface{} 
 Data []byte
 EOI [2] byte
 EncodingAlgorithm string
}

func (jpg *JPG) ShowMetadata() {
	fmt.Printf("%-32s\n", "Start of Information")
	fmt.Printf("%   -32s: %s\n",
		" File Name", jpg.Name)
	for _, segment := range jpg.Segments {
		switch s := segment.(type) {
		case APPSegment:
			showAPP(s)
		case EXIFSegment:
			showEXIF(s)
		case COMSegment:
			showCOM(s)
		case DQTSegment:
			showDQT(s)
		case SOFSegment:
			showSOF(s, jpg.EncodingAlgorithm)
		case DHTSegment:
			showDHT(s)
		case SOSSegment:
			showSOS(s)
		case Segment:
			showSegment(s)
		}
	}
	fmt.Printf("%-32s\n", "End of Information")
}
