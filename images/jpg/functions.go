package jpg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

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
	isHuffman := false
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
			identifier := string(buffer[offset+4 : offset+10])
			if identifier == "Exif\x00\x00" {
				var exifSegment = parseEXIF(buffer[offset:offset+length+2], length)
				jpg.Segments = append(jpg.Segments, exifSegment)
			} else if identifier[:len(identifier)-1] == "JFIF\x00" {
				var appSegment = parseAPP(buffer[offset:offset+length+2], length)
				jpg.Segments = append(jpg.Segments, appSegment)
			} else {
				var segment = parseSegment(buffer[offset:offset+length+2], length)
				jpg.Segments = append(jpg.Segments, segment)
			}
		// COM
		case bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xfe}):
			var comSegment = parseCOM(buffer[offset:offset+length+2], length)
			jpg.Segments = append(jpg.Segments, comSegment)
			// DQT
		case bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xdb}):
			var dqtSegment = parseDQT(buffer[offset:offset+length+2], length)
			jpg.Segments = append(jpg.Segments, dqtSegment)
			// SOF
		case
			bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xc0}) ||
				bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xc1}) ||
				bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xc2}) ||
				bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xc3}) ||
				bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xc5}) ||
				bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xc6}):
			var sofSegment = parseSOF(buffer[offset:offset+length+2], length)
			jpg.Segments = append(jpg.Segments, sofSegment)
			// DHT
		case bytes.Equal(buffer[offset:offset+2], []byte{0xff, 0xc4}):
			if !isHuffman {
				isHuffman = true
			}
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
	if isHuffman {
		jpg.EncodingAlgorithm = "Huffman coding"
	}
	return jpg, nil
}

func isJPG(soi [2]byte, eoi [2]byte) error {
	if soi != SOISegment && eoi != EOISegment {
		return fmt.Errorf("It is not a JPEG file, SOI and/or DOI does not found.")
	}
	return nil
}
