package jpg

type SOFComponentInfo struct {
 Component uint8
 Sampling_X__Y uint8
 Quantization uint8
}

var DCTEncodingTypes = []string{
  "Baseline DCT",
  "Extended Sequential DCT",
  "Progressive DCT",
  "Lossless DCT",
  "Differential Progressive DCT",
  "Hierarchical DCT",
}

var ImageComponents = [3]string{
 "Y",
 "Cb",
 "Cr",
}

type DataFormat struct {
  Format string
  Bytes_per_component uint8
}

var DataFormatIndex = []DataFormat {
 DataFormat{Format: "unknown", Bytes_per_component: 0}, // ?
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

type IFDTag struct {
 Name string
 Description string
}

type UnsignedRational struct {
 Denominator int32 
 Numerator int32
}

type SignedRational struct {
 Denominator uint32
 Numerator uint32
}

var IFDType = []string {
 "Main image",
 "Thumbnail image",
 "IFD Formatted Data (SubIFD)",
}
