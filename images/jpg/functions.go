package jpg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func (jpg *JPG) ShowMetadata() {
	fmt.Printf("%-50s: %s\n",
		"File Name:", jpg.Name)
	fmt.Printf("%-50s\n", "Start of Information")
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
		case DHTSegment:
			showDHT(s)
		case SOSSegment:
			showSOS(s)
		case Segment:
			showSegment(s)
		}
	}
	fmt.Printf("%-50s\n", "End of Information")
}

func New(filepath string) (*JPG, error) {
	// Initialization
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, fileInfo.Size())
	if _, err = io.ReadFull(file, buffer); err != nil {
		return nil, err
	}
	soi := [2]byte{buffer[0], buffer[1]}
	eoi := [2]byte{buffer[len(buffer)-1], buffer[len(buffer)-2]}
	isJPG(soi, eoi)
	// JPG base
	jpg := &JPG{Segments: make([]interface{}, 0), Name: file.Name(),
		SOI: SOISegment, EOI: EOISegment}
	offset := 2
	length := 0
	// Segments
L:
	for {
		length = int(binary.BigEndian.Uint16(buffer[offset+2 : offset+4]))
		switch {
		// APP
		case
			bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xe0}) ||
				bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xe1}) ||
				bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xe2}) ||
				bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xe3}) ||
				bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xe4}) ||
				bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xe5}) ||
				bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xe6}) ||
				bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xe7}) ||
				bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xe8}) ||
				bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xe9}) ||
				bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xea}) ||
				bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xeb}) ||
				bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xec}) ||
				bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xed}) ||
				bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xee}) ||
				bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xef}):
			if string(buffer[offset+4:offset+10]) == "Exif\x00\x00" {
				var exifSegment = parseEXIF(buffer[offset:offset+length+2], length)
				jpg.Segments = append(jpg.Segments, exifSegment)
			} else {
				var appSegment = parseAPP(buffer[offset:offset+length+2], length)
				jpg.Segments = append(jpg.Segments, appSegment)
			}
		// COM
		case bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xfe}):
			var comSegment = parseCOM(buffer[offset:offset+length+2], length)
			jpg.Segments = append(jpg.Segments, comSegment)
			// DQT
		case bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xdb}):
			var dqtSegment = parseDQT(buffer[offset:offset+length+2], length)
			jpg.Segments = append(jpg.Segments, dqtSegment)
			// DHT
		case bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xc4}):
			var dhtSegment = parseDHT(buffer[offset:offset+length+2], length)
			jpg.Segments = append(jpg.Segments, dhtSegment)
			// SOS
		case bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xda}):
			var sosSegment = parseSOS(buffer[offset:offset+length+2], length)
			jpg.Segments = append(jpg.Segments, sosSegment)
			jpg.Data = buffer[offset+length+2:]
		// Unparsed Segment
		case buffer[offset] == 0xff:
			var segment = parseSegment(buffer[offset:offset+length+2], length)
			jpg.Segments = append(jpg.Segments, segment)
		default:
			break L
		}
		offset += length + 2
	}
	return jpg, nil
}

func isJPG(soi [2]byte, eoi [2]byte) error {
	if soi != SOISegment && eoi != EOISegment {
		return fmt.Errorf("It is not a JPEG file, SOI and/or DOI does not found.")
	}
	return nil
}

func parseSegment(data []byte, size int) Segment {
	var s Segment
	s.Marker = [2]byte{data[0], data[1]}
	s.Length = uint16(size)
	s.Data = data[4:]
	return s
}

func parseAPP(data []byte, size int) APPSegment {
	var s APPSegment
	s.Marker = [2]byte{data[0], data[1]}
	s.Length = uint16(size)
	s.Identifier = string(data[4:9])
	s.Version = string(data[9]+0x30) + "." + string(data[10]+0x30)
	s.DensityUnits = uint8(data[11])
	s.Xdensity = binary.BigEndian.Uint16(data[12:14])
	s.Ydensity = binary.BigEndian.Uint16(data[14:16])
	s.XThumbnail = uint8(data[16])
	s.XThumbnail = uint8(data[17])
	if s.XThumbnail != 0 {
		copy(s.ThumbnailData, data[18:])
	}
	return s
}

func parseEXIF(data []byte, size int) EXIFSegment {
	var s EXIFSegment
	s.Marker = [2]byte{data[0], data[1]}
	s.Length = uint16(size)
	s.Identifier = string(data[4:10])
	s.TIFFHeader.Alignment = string(data[10:12])
  s.TIFFHeader.FixedBytes = [2]byte{data[12],data[13]}
  s.TIFFHeader.IFDOffset = binary.BigEndian.Uint32(data[14:18])
	return s
}

func parseCOM(data []byte, size int) COMSegment {
	var s COMSegment
	s.Marker = [2]byte{data[0], data[1]}
	s.Length = uint16(size)
	s.Data = string(data[4:])
	return s
}

func parseDQT(data []byte, size int) DQTSegment {
	var s DQTSegment
	s.Marker = [2]byte{data[0], data[1]}
	s.Length = uint16(size)
	s.Destination = uint8(data[5])
	s.Data = data[4:]
	return s
}

func parseDHT(data []byte, size int) DHTSegment {
	var s DHTSegment
	s.Marker = [2]byte{data[0], data[1]}
	s.Length = uint16(size)
	s.Class__Idx = uint8(data[4])
	offset := 5
	for i := 0; offset < 21; i++ {
		s.Bit_Codes[i] = uint8(data[offset])
		offset++
	}
	s.Real_Huffman_Codes = data[offset:]
	return s
}

func parseSOS(data []byte, size int) SOSSegment {
	var s SOSSegment
	s.Marker = [2]byte{data[0], data[1]}
	s.Length = uint16(size)
	s.Components = uint8(data[4])
	offset := 0
	switch s.Components {
	case 1:
		s.YIndex = uint8(data[5])
		s.Y_AC__DC = uint8(data[6])
		offset = 7
	case 2:
		s.YIndex = uint8(data[5])
		s.Y_AC__DC = uint8(data[6])
		s.CbIndex = uint8(data[7])
		s.Cb_AC__DC = uint8(data[8])
		offset = 9
	case 3:
		s.YIndex = uint8(data[5])
		s.Y_AC__DC = uint8(data[6])
		s.CbIndex = uint8(data[7])
		s.Cb_AC__DC = uint8(data[8])
		s.CrIndex = uint8(data[9])
		s.Cr_AC__DC = uint8(data[10])
		offset = 11
	}
	s.SS_Start = uint8(data[offset])
	s.SS_End = uint8(data[offset+1])
	s.Sucessive_approx = uint8(data[offset+2])
	return s
}

func showSegment(s Segment) {
	fmt.Printf(" %-49s: \033[1;31m%s\033[1;0m\n",
		"Segment", "Unknown")
	fmt.Printf("  %-48s: %d Bytes\n",
		"Length", s.Length)
}

func showAPP(s APPSegment) {
	num := s.Marker[1] - 0xe0
	fmt.Printf(" \033[1;34m%-49s: %s %d\033[1;0m\n",
		"Segment", "Application", num)
	fmt.Printf("  %-48s: %d Bytes\n",
		"Length", s.Length)
	fmt.Printf("  %-48s: %s\n",
		"Identifier", s.Identifier[:len(s.Identifier)-1])
	fmt.Printf("  %-48s: %s\n",
		"Version", s.Version)
	var desc string
	switch s.DensityUnits {
	case 0:
		desc = "Pixel Aspect Ratio"
	case 1:
		desc = "Pixels per inch (2.54 cm)"
	case 2:
		desc = "Pixels per centimeter"
	}
	fmt.Printf("  %-48s: %d [%s]\n",
		"Density units", s.DensityUnits, desc)
	fmt.Printf("  %-48s: %d\n",
		"X Density", s.Xdensity)
	fmt.Printf("  %-48s: %d\n",
		"Y Density", s.Ydensity)
	fmt.Printf("  %-48s: %dx%d\n",
		"Thumbnail", s.XThumbnail, s.YThumbnail)
	if len(s.ThumbnailData) == 0 {
		fmt.Printf("  %-48s: %s\n",
			"Thumbnail Size", "None")
	} else {
		fmt.Printf("  %-48s: %d Bytes\n",
			"Thumbnail Size", len(s.ThumbnailData))
	}
}

func showEXIF(s EXIFSegment) {
	num := s.Marker[1] - 0xe0
	fmt.Printf(" \033[1;34m%-49s: %s %d\033[1;0m\n",
		"Segment", "Application", num)
	fmt.Printf("  %-48s: %d Bytes\n",
		"Length", s.Length)
	fmt.Printf("  %-48s: %s\n",
		"Identifier", s.Identifier)
	fmt.Printf("  %-48s: %s\n",
		"Format", "TIFF")
  var endianness string
  switch string(s.TIFFHeader.Alignment[:]) {
    case "II":
     endianness = "little-endian"
    case "MM":
     endianness = "big-endian"
  }
	fmt.Printf("  %-48s: %s [%s]\n",
		"Endianness",s.TIFFHeader.Alignment,endianness)
}

func showCOM(s COMSegment) {
	fmt.Printf(" \033[1;34m%-49s: %s\033[1;0m\n",
		"Segment", "Comment")
	fmt.Printf("  %-48s: %d Bytes\n",
		"Length", s.Length)
	fmt.Printf("  %-48s: %s\n",
		"Comment", s.Data)
}

func showDQT(s DQTSegment) {
	fmt.Printf(" \033[1;34m%-49s: %s\033[1;0m\n",
		"Segment", "Define Quantization Table")
	fmt.Printf("  %-48s: %d Bytes\n",
		"Length", s.Length)
	higherPart := (s.Destination >> 4) & 0x0f
	lowerPart := s.Destination & 0x0f
	var component string
	switch higherPart {
	case 0:
		component = "8 bits"
	case 1:
		component = "16 bits (extended mode)"
	}
	fmt.Printf("  %-48s: %d [%s]\n",
		"Precision", higherPart, component)
	fmt.Printf("  %-48s: %d\n",
		"Table Index", lowerPart)
}

func showDHT(s DHTSegment) {
	fmt.Printf(" \033[1;34m%-49s: %s\033[1;0m\n",
		"Segment", "Define Huffman Table")
	fmt.Printf("  %-48s: %d Bytes\n",
		"Length", s.Length)
	classN := (s.Class__Idx >> 4) & 0xf
	var class string
	idx := s.Class__Idx & 0xf
	switch classN {
	case 0:
		class = "DC"
	case 1:
		class = "AC"
	}
	fmt.Printf("  %-48s: %d [%s]\n",
		"Class", classN, class)
	fmt.Printf("  %-48s: %d\n",
		"Table Index", idx)
	var codes string
	for i, v := range s.Bit_Codes {
		if i == 9 {
			codes += fmt.Sprintf("\n   %-48s ", "")
		}
		codes += fmt.Sprintf("[%d : %d] ", i+1, v)
	}
	fmt.Printf("  %-48s: %s\n",
		"Huffman Codes", codes)
}

func showSOS(s SOSSegment) {
	fmt.Printf(" \033[1;34m%-49s: %s\033[1;0m\n",
		"Segment", "Start Of Scan")
	fmt.Printf("  %-48s: %d Bytes\n",
		"Length", s.Length)
	fmt.Printf("  %-48s: %d\n",
		"Number of Components", s.Components)
	yDC := s.Y_AC__DC & 0xf
	yAC := (s.Y_AC__DC >> 4) & 0xf
	fmt.Printf("  %-48s: %d [DC Table Index] %d [AC Table Index]\n",
		"Luminance(Y)", yDC, yAC)
	if s.Components > 1 {
		CbDC := s.Cb_AC__DC & 0xf
		CbAC := (s.Cb_AC__DC >> 4) & 0xf
		fmt.Printf("  %-48s: %d [DC Table Index] %d [AC Table Index]\n",
			"Crominance(Cb)", CbDC, CbAC)
		if s.Components > 2 {
			CrDC := s.Cb_AC__DC & 0xf
			CrAC := (s.Cb_AC__DC >> 4) & 0xf
			fmt.Printf("  %-48s: %d [DC Table Index] %d [AC Table Index]\n",
				"Crominance(Cr)", CrDC, CrAC)
		}
	}
	fmt.Printf("  %-48s: %d\n",
		"Start of spectral selection", s.SS_Start)
	fmt.Printf("  %-48s: %d\n",
		"End of spectral selection", s.SS_End)
	fmt.Printf("  %-48s: %d\n",
		"Sucesive approximation bits", s.Sucessive_approx)
}
