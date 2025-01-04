package jpg

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
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

func findIFD0Tag(idx uint16) IFDTag {
	// IFD 0 (Main image)
	IFDTags := map[uint16]IFDTag{
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
   // Misc
   0x00fe: {Name: "NewSubfileType",
      Description: "",},
   0x00ff: {Name: "SubfileType",
      Description: "",},
   0x012d: {Name: "TransferFunction",
      Description: "",},
   0x013b: {Name: "Artist",
      Description: "",},
   0x013d: {Name: "Predictor",
      Description: "",},
   0x0142: {Name: "TileWidth",
      Description: "",},
   0x0143: {Name: "TileLength",
      Description: "",},
   0x0144: {Name: "TileOffsets",
      Description: "",},
   0x0145: {Name: "TileByteCounts",
      Description: "",},
   0x014a: {Name: "SubIFDs",
      Description: "",},
   0x015b: {Name: "JPEGTables",
      Description: "",},
   0x828d: {Name: "CFARepeatPatternDim",
      Description: "",},
   0x828e: {Name: "CFAPattern",
      Description: "",},
   0x828f: {Name: "BatteryLevel",
      Description: "",},
   0x83bb: {Name: "IPTC/NAA",
      Description: "",},
   0x8773: {Name: "InterColorProfile",
      Description: "",},
   0x8824: {Name: "SpectralSensitivity",
      Description: "",},
   0x8825: {Name: "GPSInfoOffset",
      Description: "",},
   0x8828: {Name: "OECF",
      Description: "",},
   0x8829: {Name: "Interlace",
      Description: "",},
   0x882a: {Name: "TimeZoneOffset",
      Description: "",},
   0x882b: {Name: "SelfTimerMode",
      Description: "",},
   0x920b: {Name: "FlashEnergy",
      Description: "",},
   0x920c: {Name: "SpatialFrequencyResponse",
      Description: "",},
   0x920d: {Name: "Noise",
      Description: "",},
   0x9211: {Name: "ImageNumber",
      Description: "",},
   0x9212: {Name: "SecurityClassification",
      Description: "",},
   0x9213: {Name: "ImageHistory",
      Description: "",},
   0x9214: {Name: "SubjectLocation",
      Description: "",},
   0x9215: {Name: "ExposureIndex",
      Description: "",},
   0x9216: {Name: "TIFF/EPStandardID",
      Description: "",},
   0x9290: {Name: "SubSecTime",
      Description: "",},
   0x9291: {Name: "SubSecTimeOriginal",
      Description: "",},
   0x9292: {Name: "SubSecTimeDigitized",
      Description: "",},
   0xa20b: {Name: "FlashEnergy",
      Description: "",},
   0xa20c: {Name: "SpatialFrequencyResponse",
      Description: "",},
   0xa214: {Name: "SubjectLocation",
      Description: "",},
   0xa215: {Name: "ExposureIndex",
      Description: "",},
   0xa302: {Name: "CFAPattern",
      Description: "", },
  }
  return IFDTags[idx]
}

func findIFD1Tag(idx uint16) IFDTag {
	IFDTags := map[uint16]IFDTag{
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
    0x0117: {Name: "StripByteCounts",
      Description: "When image format is no compression and stored as strip, this value shows how many bytes used for each strip and this value is plural. If image has not stripped, this value is single and means whole data size of image."},
		0x011a: {Name: "XResolution",
			Description: "Display/Print resolution of image. Large number of digicam uses 1/72inch, but it has no mean because personal computer doesn't use this value to display/print out."},
		0x011b: {Name: "YResolution",
			Description: "Display/Print resolution of image. Large number of digicam uses 1/72inch, but it has no mean because personal computer doesn't use this value to display/print out."},
		0x011c: {Name: "PlanarConfiguration",
			Description: "When image format is no compression YCbCr, this value shows byte aligns of YCbCr data. If value is '1', Y/Cb/Cr value is chunky format, contiguous for each subsampling pixel. If value is '2', Y/Cb/Cr value is separated and stored to Y plane/Cb plane/Cr plane format."},
    0x0128: {Name: "ResolutionUnit",
      Description: "Unit of XResolution(0x011a)/YResolution(0x011b). '1' means inch, '2' means centimeter."},
    0x0201: {Name: "JpegIFOffset",
      Description: "When image format is JPEG, this value show offset to JPEG data stored."},
    0x0202: {Name: "JpegIFByteCount",
      Description: "When image format is JPEG, this value shows data size of JPEG image."},
    0x0211: {Name: "YCbCrCoefficients",
      Description: "When image format is YCbCr, this value shows constants to translate it to RGB format. In usual, '0.299/0.587/0.114' are used."},
    0x0212: {Name: "YCbCrSubSampling",
      Description: "When image format is YCbCr and uses subsampling(cropping of chroma data, all the digicam do that), this value shows how many chroma data subsampled. First value shows horizontal, next value shows vertical subsample rate."},
    0x0213: {Name: "YCbCrPositioning",
      Description: "When image format is YCbCr and uses 'Subsampling'(cropping of chroma data, all the digicam do that), this value defines the chroma sample point of subsampled pixel array. '1' means the center of pixel array, '2' means the datum point(0,0)."},
    0x0214: {Name: "ReferenceBlackWhite",
      Description: "Shows reference value of black point/white point. In case of YCbCr format, first 2 show black/white of Y, next 2 are Cb, last 2 are Cr. In case of RGB format, first 2 show black/white of R, next 2 are G, last 2 are B."},
   // Misc
   0x00fe: {Name: "NewSubfileType",
      Description: "",},
   0x00ff: {Name: "SubfileType",
      Description: "",},
   0x012d: {Name: "TransferFunction",
      Description: "",},
   0x013b: {Name: "Artist",
      Description: "",},
   0x013d: {Name: "Predictor",
      Description: "",},
   0x0142: {Name: "TileWidth",
      Description: "",},
   0x0143: {Name: "TileLength",
      Description: "",},
   0x0144: {Name: "TileOffsets",
      Description: "",},
   0x0145: {Name: "TileByteCounts",
      Description: "",},
   0x014a: {Name: "SubIFDs",
      Description: "",},
   0x015b: {Name: "JPEGTables",
      Description: "",},
   0x828d: {Name: "CFARepeatPatternDim",
      Description: "",},
   0x828e: {Name: "CFAPattern",
      Description: "",},
   0x828f: {Name: "BatteryLevel",
      Description: "",},
   0x83bb: {Name: "IPTC/NAA",
      Description: "",},
   0x8773: {Name: "InterColorProfile",
      Description: "",},
   0x8824: {Name: "SpectralSensitivity",
      Description: "",},
   0x8825: {Name: "GPSInfoOffset",
      Description: "",},
   0x8828: {Name: "OECF",
      Description: "",},
   0x8829: {Name: "Interlace",
      Description: "",},
   0x882a: {Name: "TimeZoneOffset",
      Description: "",},
   0x882b: {Name: "SelfTimerMode",
      Description: "",},
   0x920b: {Name: "FlashEnergy",
      Description: "",},
   0x920c: {Name: "SpatialFrequencyResponse",
      Description: "",},
   0x920d: {Name: "Noise",
      Description: "",},
   0x9211: {Name: "ImageNumber",
      Description: "",},
   0x9212: {Name: "SecurityClassification",
      Description: "",},
   0x9213: {Name: "ImageHistory",
      Description: "",},
   0x9214: {Name: "SubjectLocation",
      Description: "",},
   0x9215: {Name: "ExposureIndex",
      Description: "",},
   0x9216: {Name: "TIFF/EPStandardID",
      Description: "",},
   0x9290: {Name: "SubSecTime",
      Description: "",},
   0x9291: {Name: "SubSecTimeOriginal",
      Description: "",},
   0x9292: {Name: "SubSecTimeDigitized",
      Description: "",},
   0xa20b: {Name: "FlashEnergy",
      Description: "",},
   0xa20c: {Name: "SpatialFrequencyResponse",
      Description: "",},
   0xa214: {Name: "SubjectLocation",
      Description: "",},
   0xa215: {Name: "ExposureIndex",
      Description: "",},
   0xa302: {Name: "CFAPattern",
      Description: "", },
  }
	return IFDTags[idx]
}

func findSubIFDTag(idx uint16) IFDTag {
  // SubIFD
	IFDTags := map[uint16]IFDTag{ 
   0x829a: {Name: "ExposureTime",
      Description: "Exposure time (reciprocal of shutter speed). Unit is second."},
   0x829d: {Name: "FNumber",
      Description: "The actual F-number(F-stop) of lens when the image was taken."},
   0x8822: {Name: "ExposureProgram",
      Description: "Exposure program that the camera used when image was taken. '1' means manual control, '2' program normal, '3' aperture priority, '4' shutter priority, '5' program creative (slow program), '6' program action(high-speed program), '7' portrait mode, '8' landscape mode."},
   0x8827: {Name: "ISOSpeedRatings",
      Description: "CCD sensitivity equivalent to Ag-Hr film speedrate."},
   0x9000: {Name: "ExifVersion",
      Description: "Exif version number. Stored as 4bytes of ASCII character (like \"0210\")"},
   0x9003: {Name: "DateTimeOriginal",
      Description: "Date/Time of original image taken. This value should not be modified by user program."},
   0x9004: {Name: "DateTimeDigitized",
      Description: "Date/Time of image digitized. Usually, it contains the same value of DateTimeOriginal(0x9003)."},
   0x9101: {Name: "ComponentConfiguration",
      Description: "Unknown. It seems value 0x00,0x01,0x02,0x03 always."},
   0x9102: {Name: "CompressedBitsPerPixel",
      Description: "The average compression ratio of JPEG."},
   0x9201: {Name: "ShutterSpeedValue",
      Description: "Shutter speed. To convert this value to ordinary 'Shutter Speed'; calculate this value's power of 2, then reciprocal. For example, if value is '4', shutter speed is 1/(2^4)=1/16 second."},
   0x9202: {Name: "ApertureValue",
      Description: "The actual aperture value of lens when the image was taken. To convert this value to ordinary F-number(F-stop), calculate this value's power of root 2 (=1.4142). For example, if value is '5', F-number is 1.4142^5 = F5.6."},
   0x9203: {Name: "BrightnessValue",
      Description: "Brightness of taken subject, unit is EV."},
   0x9204: {Name: "ExposureBiasValue",
      Description: "Exposure bias value of taking picture. Unit is EV."},
   0x9205: {Name: "MaxApertureValue ",
      Description: "Maximum aperture value of lens. You can convert to F-number by calculating power of root 2 (same process of ApertureValue(0x9202)."},
   0x9206: {Name: "SubjectDistance",
      Description: "Distance to focus point, unit is meter."},
   0x9207: {Name: "MeteringMode",
      Description: "Exposure metering method. '1' means average, '2' center weighted average, '3' spot, '4' multi-spot, '5' multi-segment."},
   0x9208: {Name: "LightSource",
      Description: "Light source, actually this means white balance setting. '0' means auto, '1' daylight, '2' fluorescent, '3' tungsten, '10' flash."},
   0x9209: {Name: "Flash",
      Description: "'1' means flash was used, '0' means not used."},
   0x920a: {Name: "FocalLength",
      Description: "Focal length of lens used to take image. Unit is millimeter."},
   0x927c: {Name: "MakerNote",
      Description: "Maker dependent internal data. Some of maker such as Olympus/Nikon/Sanyo etc. uses IFD format for this area."},
   0x9286: {Name: "UserComment",
      Description: "Stores user comment."},
   0xa000: {Name: "FlashPixVersion",
      Description: "Stores FlashPix version. Unknown but 4bytes of ASCII characters \"0100\"exists."},
   0xa001: {Name: "ColorSpace",
      Description: "Unknown, value is '1'."},
   0xa002: {Name: "ExifImageWidth",
      Description: "Size of main image."},
   0xa003: {Name: "ExifImageHeight",
      Description: "Size of main image."},
   0xa004: {Name: "RelatedSoundFile",
      Description: "If this digicam can record audio data with image, shows name of audio data."},
   0xa005: {Name: "ExifInteroperabilityOffset",
      Description: "Extension of \"ExifR98\", detail is unknown. This value is offset to IFD format data. Currently there are 2 directory entries, first one is Tag0x0001, value is \"R98\", next is Tag0x0002, value is \"0100\"."},
   0xa20e: {Name: "FocalPlaneXResolution",
      Description: "CCD's pixel density."},
   0xa20f: {Name: "FocalPlaneYResolution",
      Description: "CCD's pixel density."},
   0xa210: {Name: "FocalPlaneResolutionUnit",
      Description: "Unit of FocalPlaneXResoluton/FocalPlaneYResolution. '1' means no-unit, '2' inch, '3' centimeter."},
   0xa217: {Name: "SensingMethod",
      Description: "Shows type of image sensor unit. '2' means 1 chip color area sensor, most of all digicam use this type."},
   0xa300: {Name: "FileSource",
      Description: "Unknown but value is '3'."},
   0xa301: {Name: "SceneType",
      Description: "Unknown but value is '1'."},
   // Misc
   0x00fe: {Name: "NewSubfileType",
      Description: "",},
   0x00ff: {Name: "SubfileType",
      Description: "",},
   0x012d: {Name: "TransferFunction",
      Description: "",},
   0x013b: {Name: "Artist",
      Description: "",},
   0x013d: {Name: "Predictor",
      Description: "",},
   0x0142: {Name: "TileWidth",
      Description: "",},
   0x0143: {Name: "TileLength",
      Description: "",},
   0x0144: {Name: "TileOffsets",
      Description: "",},
   0x0145: {Name: "TileByteCounts",
      Description: "",},
   0x014a: {Name: "SubIFDs",
      Description: "",},
   0x015b: {Name: "JPEGTables",
      Description: "",},
   0x828d: {Name: "CFARepeatPatternDim",
      Description: "",},
   0x828e: {Name: "CFAPattern",
      Description: "",},
   0x828f: {Name: "BatteryLevel",
      Description: "",},
   0x83bb: {Name: "IPTC/NAA",
      Description: "",},
   0x8773: {Name: "InterColorProfile",
      Description: "",},
   0x8824: {Name: "SpectralSensitivity",
      Description: "",},
   0x8825: {Name: "GPSInfoOffset",
      Description: "",},
   0x8828: {Name: "OECF",
      Description: "",},
   0x8829: {Name: "Interlace",
      Description: "",},
   0x882a: {Name: "TimeZoneOffset",
      Description: "",},
   0x882b: {Name: "SelfTimerMode",
      Description: "",},
   0x920b: {Name: "FlashEnergy",
      Description: "",},
   0x920c: {Name: "SpatialFrequencyResponse",
      Description: "",},
   0x920d: {Name: "Noise",
      Description: "",},
   0x9211: {Name: "ImageNumber",
      Description: "",},
   0x9212: {Name: "SecurityClassification",
      Description: "",},
   0x9213: {Name: "ImageHistory",
      Description: "",},
   0x9214: {Name: "SubjectLocation",
      Description: "",},
   0x9215: {Name: "ExposureIndex",
      Description: "",},
   0x9216: {Name: "TIFF/EPStandardID",
      Description: "",},
   0x9290: {Name: "SubSecTime",
      Description: "",},
   0x9291: {Name: "SubSecTimeOriginal",
      Description: "",},
   0x9292: {Name: "SubSecTimeDigitized",
      Description: "",},
   0xa20b: {Name: "FlashEnergy",
      Description: "",},
   0xa20c: {Name: "SpatialFrequencyResponse",
      Description: "",},
   0xa214: {Name: "SubjectLocation",
      Description: "",},
   0xa215: {Name: "ExposureIndex",
      Description: "",},
   0xa302: {Name: "CFAPattern",
      Description: "",},
   0xa401: {Name: "CustomRendered",
      Description: "",},
   0xa402: {Name: "ExposureMode",
      Description: "",},
   0xa403: {Name: "WhiteBalance",
      Description: "",},
   0xa404: {Name: "DigitalZoomRatio",
      Description: "",},
   0xa405: {Name: "FocalLength35mmFormat",
      Description: "",},
   0xa406: {Name: "SceneCaptureType",
      Description: "",},
   0xa407: {Name: "GainControl",
      Description: "",},
   0xa408: {Name: "Contrast",
      Description: "",},
   0xa409: {Name: "Saturation",
      Description: "",},
   0xa40A: {Name: "Sharpness",
      Description: "",},
   0xa40B: {Name: "DeviceSettingDescription",
      Description: "",},
   0xa40C: {Name: "SubjectDistanceRange",
      Description: "",},
  }
  return IFDTags[idx] 
}

func findGPSIFDTag(idx uint16) IFDTag {
	IFDTags := map[uint16]IFDTag{ 
    0x0000: {Name: "GPSVersionID", Description:""},
    0x0001: {Name: "GPSLatitudeRef", Description:""},
    0x0002: {Name: "GPSLatitude", Description:""},
    0x0003: {Name: "GPSLongitudeRef", Description:""},
    0x0004: {Name: "GPSLongitude", Description:""},
    0x0005: {Name: "GPSAltitudeRef", Description:""},
    0x0006: {Name: "GPSAltitude", Description:""},
    0x0007: {Name: "GPSTimeStamp", Description:""},
    0x0008: {Name: "GPSSatellites", Description:""},
    0x0009: {Name: "GPSStatus", Description:""},
    0x000a: {Name: "GPSMeasureMode", Description:""},
    0x000b: {Name: "GPSDOP", Description:""},
    0x000c: {Name: "GPSSpeedRef ", Description:""},
    0x000d: {Name: "GPSSpeed", Description:""},
    0x000e: {Name: "GPSTrackRef", Description:""},
    0x000f: {Name: "GPSTrack", Description:""},
    0x0010: {Name: "GPSImgDirectionRef", Description:""},
    0x0011: {Name: "GPSImgDirection", Description:""},
    0x0012: {Name: "GPSMapDatum", Description:""},
    0x0013: {Name: "GPSDestLatitudeRef", Description:""},
    0x0014: {Name: "GPSDestLatitude", Description:""},
    0x0015: {Name: "GPSDestLongitudeRef", Description:""},
    0x0016: {Name: "GPSDestLongitude", Description:""},
    0x0017: {Name: "GPSDestBearingRef", Description:""},
    0x0018: {Name: "GPSDestBearing", Description:""},
    0x0019: {Name: "GPSDestDistanceRef", Description:""},
    0x001a: {Name: "GPSDestDistance", Description:""},
    0x001b: {Name: "GPSProcessingMethod", Description:""},
    0x001c: {Name: "GPSAreaInformation", Description:""},
    0x001d: {Name: "GPSDateStamp", Description:""},
    0x001e: {Name: "GPSDifferential", Description:""},
    0x001f: {Name: "GPSHPositioningError", Description:""},
  }
  return IFDTags[idx]
}


// Exif Entry Data Parser
func EntryDataOf(data []byte, df DataFormat, endianness binary.ByteOrder) interface{} {
	switch df.Format {
	case "unsigned byte":
		return uint8(data[0])
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
      value.Numerator = binary.BigEndian.Uint32(data[0:4])
      value.Denominator = binary.BigEndian.Uint32(data[4:8])
			return value
    // Little Endian
		} else {
      value.Numerator = binary.LittleEndian.Uint32(data[0:4])
      value.Denominator = binary.LittleEndian.Uint32(data[4:8])
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
      binary.Read(bytes.NewReader(data[0:4]), binary.BigEndian, &value.Numerator)
      binary.Read(bytes.NewReader(data[4:8]), binary.BigEndian, &value.Denominator)
			return value
    // Little Endian
		} else {
      binary.Read(bytes.NewReader(data[0:4]), binary.LittleEndian, &value.Numerator)
      binary.Read(bytes.NewReader(data[4:8]), binary.LittleEndian, &value.Denominator)
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

func parseAndPrintGPS(latRef, lat, lonRef, lon string) string {
  proceseCoordinate := func(coordStr string) (int, int, float64) {
   parts := strings.Split(coordStr,".")
   degreesPart,_ := strconv.Atoi(parts[0])
   decimalPart,_ := strconv.ParseFloat("0."+parts[1],64)
   deg := degreesPart / 10
   min := degreesPart % 10
   sec := decimalPart * 60
   return deg, min, sec
  }
 latDeg, latMin, latSec := proceseCoordinate(lat)
 lonDeg, lonMin, lonSec := proceseCoordinate(lon)
 return fmt.Sprintf("%d deg %d' %.2f\" %s %d deg %d' %.2f\" %s",
   latDeg, latMin, latSec, latRef, lonDeg, lonMin, lonSec, lonRef)
}


