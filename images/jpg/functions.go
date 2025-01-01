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
			identifier := string(buffer[offset+4 : offset+15])
			if identifier[:6] == "Exif\x00\x00" {
				var exifSegment = parseEXIF(buffer[offset:offset+length+2], length)
				jpg.Segments = append(jpg.Segments, exifSegment)
			} else if identifier[:5] == "JFIF\x00" {
				var appSegment = parseAPP(buffer[offset:offset+length+2], length)
				jpg.Segments = append(jpg.Segments, appSegment)
			} else if identifier == "ICC_PROFILE" {
				var appSegment = parseICC(buffer[offset:offset+length+2], length)
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

func findEXIFTag(idx uint16) IFDTag {
	IFDTags := map[uint16]IFDTag{
		// IFD 0 (Main image)
		0x010e: {Name: "Image Description",
			Description: "Describes image"},
		0x010f: {Name: "Make",
			Description: "Show manufacturer of digicam"},
		0x0110: {Name: "Model",
			Description: "Shows model number of digicam"},
		0x0112: {Name: "Orientation",
			Description: "The orientation of the camera relative to the scene, when the image wa captured. The start point of stored data is,'1' means upper left, '3' lower right, '6' upper left, '8' lower left, '9' undefined."},
		0x011a: {Name: "XResolution",
			Description: "Display/Print resolution of image"},
		0x011b: {Name: "YResolution",
			Description: "Display/Print resolution of image"},
		0x0128: {Name: "ResolutionUnit",
			Description: "Unit of XResolution/YResolution. Same meaning that JFIF image density."},
		0x0131: {Name: "Software",
			Description: "Shows firmware(internal software of digicam) version number."},
		0x0132: {Name: "DateTime",
			Description: "Date/Time of image was last modified. Format is 'YYYY:MM:DD'."},
		0x013e: {Name: "WhitePoint",
			Description: "Defines chromaticity of white point of the image"},
		0x013f: {Name: "PrimaryChromaticities",
			Description: "Defines chromaticity of the primaries of the image"},
		0x0211: {Name: "YCbCrCoefficients",
			Description: "When image is YCbCr, this vaule shows a constant to translate it to RGB format"},
		0x0213: {Name: "YCbCrPositioning",
			Description: "When image is YCbCr and uses Subsampling, defines the chroma sample pint of subsamplng pixe array. '1' means the center pixel array, '2' means the datum point"},
		0x0214: {Name: "ReferenceBlackWhite",
			Description: "Shows reference value of black point/white point"},
		0x8298: {Name: "CopyRight",
			Description: "Shows copyright information"},
		0x8769: {Name: "ExitOffset",
			Description: "offset to Exit Sub IFD"},
    // IFD1 (Thumbnail image)
		0x0100: {Name: "ImageWidth",
			Description: "Shows size of thumbnail image"},
		0x0101: {Name: "ImageLength",
			Description: "Shows size of thumbnail image"},
		0x0102: {Name: "BitsPerSample",
			Description: "When image format is no compression, this value shows the number of bits per component for each pixel"}, 
		0x0103: {Name: "Compression",
			Description: "Shows compression method. 1 means no compression, 6 means JPEG compression."},
		0x0106: {Name: "PhotometricInterpretation",
			Description: "Shows the color space of the image components. 1 means monochrome, 2 means RGB, 6 means YCbCr"},
		0x0111: {Name: "StripOffsets",
			Description: "When image format is no compression, this value shows offset to image data. In some case image data is striped and this value is plural."},
		0x0115: {Name: "SamplesPerPixel",
			Description: "When image format is no compression, this value shows the number of omponents stored for each pixel"},
		0x0116: {Name: "RowsPerStrip",
			Description: "When image format is no compression and stored as strip show how many rows stored to each strip. If image has not striped, this value is the same as ImageLength."},
	}
	return IFDTags[idx]
}

// Exif Entry Data Parser
func EntryDataOf(data []byte, df DataFormat, endianness binary.ByteOrder) interface{} {
	switch df.Format {
	case "unsigned byte":
		return data
	case "ascii strings":
		return string(data)
	case "unsigned short":
		if endianness == binary.BigEndian {
			return binary.BigEndian.Uint16(data)
		}
		return binary.LittleEndian.Uint16(data)
	case "unsigned long":
		if endianness == binary.BigEndian {
			return binary.BigEndian.Uint32(data)
		}
		return binary.LittleEndian.Uint32(data)
	case "unsigned rational":
		var value UnsignedRational
		if endianness == binary.BigEndian {
			binary.Read(bytes.NewReader(data[:4]), binary.BigEndian, &value.Denominator)
			binary.Read(bytes.NewReader(data[4:8]), binary.BigEndian, &value.Numerator)
			return value
		} else {
			binary.Read(bytes.NewReader(data[:4]), binary.BigEndian, &value.Denominator)
			binary.Read(bytes.NewReader(data[4:8]), binary.BigEndian, &value.Numerator)
			return value
		}
	case "signed byte":
		return int8(data[0])
	case "undefined":
		return data
	case "signed short":
		var value int16
		if endianness == binary.BigEndian {
			binary.Read(bytes.NewReader(data), binary.BigEndian, &value)
			return value
		} else {
			binary.Read(bytes.NewReader(data), binary.LittleEndian, &value)
			return value
		}
	case "signed long":
		var value int32
		if endianness == binary.BigEndian {
			binary.Read(bytes.NewReader(data), binary.BigEndian, &value)
			return value
		} else {
			binary.Read(bytes.NewReader(data), binary.LittleEndian, &value)
			return value
		}
	case "signed rational":
		var value SignedRational
		if endianness == binary.BigEndian {
			binary.Read(bytes.NewReader(data[:4]), binary.BigEndian, &value.Denominator)
			binary.Read(bytes.NewReader(data[4:8]), binary.BigEndian, &value.Numerator)
			return value
		} else {
			binary.Read(bytes.NewReader(data[:4]), binary.BigEndian, &value.Denominator)
			binary.Read(bytes.NewReader(data[4:8]), binary.BigEndian, &value.Numerator)
			return value
		}
	case "signed float":
		var value float32
		if endianness == binary.BigEndian {
			binary.Read(bytes.NewReader(data), binary.BigEndian, &value)
			return value
		} else {
			binary.Read(bytes.NewReader(data), binary.LittleEndian, &value)
			return value
		}
	case "double float":
		var value float64
		if endianness == binary.BigEndian {
			binary.Read(bytes.NewReader(data), binary.BigEndian, &value)
			return value
		} else {
			binary.Read(bytes.NewReader(data), binary.LittleEndian, &value)
			return value
		}
	default:
		return data
	}
}
