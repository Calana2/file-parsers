package jpg

type EXIFSegment struct {
	// Application
	Marker     [2]byte
	Length     uint16
	Identifier string
	TIFFHeader TIFFHeader
	IFDs       []IFD
}

type TIFFHeader struct {
	Alignment  string
	FixedBytes [2]byte
	IFDOffset  uint32
}

type IFD struct {
	EntriesNum    uint16
	Entries       []IFDEntry
	OffsetNextIFD uint32
}

type IFDEntry struct {
	Tag           uint16
	Format        uint16
	ComponentsNum uint32
	Data_Offset   uint32
	Data          []interface{}
}

type IFDTag struct {
	Name        string
	Description string
}

var IFDType = []string{
	"Main image",
	"Thumbnail image",
	"IFD Formatted Data -- SubIFD",
  "GPS IFD",
}

type DataFormat struct {
	Format              string
	Bytes_per_component uint8
}

var DataFormatIndex = []DataFormat{
	DataFormat{Format: "unknown", Bytes_per_component: 0}, // Nope
	DataFormat{Format: "unsigned byte", Bytes_per_component: 1},
	DataFormat{Format: "ascii strings", Bytes_per_component: 1},
	DataFormat{Format: "unsigned short", Bytes_per_component: 2},
	DataFormat{Format: "unsigned long", Bytes_per_component: 4},
	DataFormat{Format: "unsigned rational", Bytes_per_component: 8},
	DataFormat{Format: "signed byte", Bytes_per_component: 1},
	DataFormat{Format: "undefined", Bytes_per_component: 1},
	DataFormat{Format: "signed short", Bytes_per_component: 2},
	DataFormat{Format: "signed long", Bytes_per_component: 4},
	DataFormat{Format: "signed rational", Bytes_per_component: 8},
	DataFormat{Format: "signed float", Bytes_per_component: 4},
	DataFormat{Format: "double float", Bytes_per_component: 8},
}

// 0x012
var Orientation = []string{
	"None",
	"Horizontal (normal)",
	"Mirror Horizontal",
	"Rotate 180",
	"Mirror vertical",
	"Mirror horizontal and rotate 270 CW",
	"Rotate 90 CW",
	"Mirror horizontal and rotate 90 CW",
	"Rotate 270 CW",
}

// 0x0128 - IFD0
var ResolutionUnit_1 = []string{
	"None",
	"No unit",
	"Inch",
	"Centimeter",
}

// 0x0128 - IFD1
var ResolutionUnit_2 = []string{
	"None",
	"Inch",
	"Centimeter",
}

// 0x0213 - IFD0
var YCbCrPositioning = []string{
	"None",
	"Centered",
	"Datum Point (0,0)",
}

// 0x0103 - IFD1
var Compression = map[uint16]string{
	1:     "Uncompressed",
	6:     "JPEG (old-style)",
	7:     "JPEG",
}

// 0x0106 - subIFD
var PhotometricInterpretation = []string{
	1: "Monochrome",
	2: "RGB",
	6: "YCbCr",
}

// 0x011c - subIFD
var PlanarConfiguration = []string{
	1: "Chunky",
	2: "Planar",
}

// 0x8822 - subIFD
var ExposureProgram = []string{
	"Not defined",
	"Manual Control",
	"Program Normal",
	"Aperture Priority",
	"Shutter Priority",
	"Program Creative",
	"Program Action",
	"Portrait Mode",
	"Landscape Mode",
	"Bulb",
}

// 0x9207 - subIFD
var MeteringMode = map[uint16]string{
  0:"Unknown",
  1:"Average",
  2:"Center Weighted Average",
  3:"Spot",
  4:"Multi-Spot",
  5:"Multi-Segment",
  6:"Partial",
	255: "Other",
}

// 0x9208 - subIFD
var LightSource = map[uint16]string{
  0:"Unknown",
  1:"Daylight",
  2:"Fluorescent",
  3:"Tungsten (incandescent",
  4:"Flash",
	9: "Fine Weather",
  10:"Cloudy",
  11:"Shade",
  12:"Daylight",
  13:"Daylight Fluorescent",
  14:"Day white Fluorescent",
  15:"Cool white Fluorescent",
  16:"White Fluorescent",
  17:"Warm White Fluorescent",
  18:"Standart Light A",
  19:"Standart Light B",
  20:"Standart Light C",
  21:"D55",
  22:"D65",
  23:"D75",
  24:"D50",
  25:"ISO Studio Tungsten",
	255: "Other",
}

// 0x9101 - subIFD
var ComponentConfiguration = [7]string{
	"-",
	"Y",
	"Cb",
	"Cr",
	"R",
	"G",
	"B",
}

// 0x09209 - subIFD
var Flash = []string{
	0x0:  "No Flash",
	0x1:  "Fired",
	0x5:  "Fired, Return not detected",
	0x7:  "Fired, Return detected",
	0x8:  "On, Did not fire",
	0x9:  "On, Fired",
	0xD:  "On, Return not detected",
	0xF:  "On, Return detected",
	0x10: "Off, Did not fire",
	0x14: "Off, Did not fire, Return not detected",
	0x18: "Auto, Did not fire",
	0x19: "Auto, Fired",
	0x1D: "Auto, Fired, Return not detected",
	0x1F: "Auto,Fired, Return detected",
	0x20: "No flash function",
	0x30: "Off, No flash function",
	0x41: "Fired, Red-eye reduction",
	0x45: "Fired, Red-eye reduction, Return not detected",
	0x47: "Fired, Red-eye reuction, Return detected",
	0x49: "On, Red-eye reduction",
	0x4D: "On, Red-eye reduction, Return not detected",
	0x4F: "On, Red-eye reuction, Return detected",
	0x50: "Off, Red-eye reduction",
	0x58: "Auto, Did not fire, Red-eye reduction",
	0x59: "Auto, Fired, Red-eye reduction",
	0x5D: "Auto, Fired, Red-eye reduction, Return not detected",
	0x5F: "Auto, Fired, Red-eye reduction, Return detected",
}

// 0xa001 - subIFD
var ColorSpace = map[uint16]string{
	1:      "sRGB",
	2:      "Adobe RGB",
	0xfffd: "Wide Gamut RGB",
	0xfffe: "ICC Profile",
	0xffff: "Uncalibrated",
}

// 0xa210 - subIFD
var FocalPlaneResolutionUnit = []string{
	1: "No-unit",
	2: "Inches",
	3: "cm",
	4: "mm",
	5: "um",
}

// 0xa217 - subIFD
var SensingMethod = []string{
	1: "Not defined",
	2: "One-chip color area",
	3: "Two-chip color area",
	4: "Three-chip color area",
	5: "Color sequential area",
	7: "Trilinear",
	8: "Color sequential linear",
}

// 0xa300 - subIFD
var FileSource = map[uint16]string{
	1: "Film Scanner",
	2: "Reflection Print Scanner",
	3: "Digital Camera",
}

// 0xa301 - subIFD
var SceneType = []string{
	1: "Directly Photographed",
}

// 0xa401 - subIFD
var CustomRendered = [9]string{
	0: "Normal",
	1: "Custom",
	2: "HDR (no original saved)",
	3: "HDR (original saved)",
	4: "Original (for HDR)",
	6: "Panorama",
	7: "Portrait HDR",
	8: "Portrait",
}

// 0xa402 - subIFD
var ExposureMode = [3]string{
	0: "Auto",
	1: "Manual",
	2: "Auto bracket",
}

// 0xa403 - subIFD
var WhiteBalance = [2]string{
	0: "Auto",
	1: "Manual",
}

// 0xa406 - subIFD
var SceneCaptureType = [5]string{
	0: "Standard",
	1: "Landscape",
	2: "Portrait",
	3: "Night",
	4: "Other",
}

// 0xa407 - subIFD
var GainControl = [5]string{
	0: "None",
	1: "Low gain up",
	2: "High gain up",
	3: "Low gain down",
	4: "High gain down",
}

// 0xa408 - subIFD
var Contrast = [3]string{
	0: "Normal",
	1: "Low",
	2: "High",
}

// 0xa409 - subIFD
var Saturation = [3]string{
	0: "Normal",
	1: "Low",
	2: "High",
}

// 0xa40a - subIFD
var Sharpness = [3]string{
	0: "Normal",
	1: "Soft",
	2: "Hard",
}

// 0xa40c - subIFD
var SubjectDistanceRange = [4]string{
	0: "Unknown",
	1: "Macro",
	2: "Close",
	3: "Distant",
}

type GPSData struct {
 LongitudeValue string
 LongitudeRef string
 LatitudeValue string
 LatitudeRef string
}

