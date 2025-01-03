package jpg

import "fmt"

type SOFComponentInfo struct {
	Component     uint8
	Sampling_X__Y uint8
	Quantization  uint8
}

var DCTEncodingTypes = []string{
	"Baseline DCT",
	"Extended Sequential DCT",
	"Progressive DCT",
	"Lossless DCT",
	"Differential Progressive DCT",
	"Hierarchical DCT",
}

var ImageComponents = [6]string{
	"Y",
	"Cb",
	"Cr",
  "R",
  "G",
  "B",
}

type UnsignedRational struct {
	Numerator   uint32
	Denominator uint32
}

type SignedRational struct {
	Numerator   int32
	Denominator int32
}

func (ur *UnsignedRational) Representation() string {
	if ur.Denominator != 0 {
		return fmt.Sprintf("%v", float32(ur.Numerator)/float32(ur.Denominator))
	}
	return fmt.Sprintf("Undefined / undetermined : Denominator 0")
}

func (sr *SignedRational) Representation() string {
	if sr.Denominator != 0 {
		return fmt.Sprintf("%v", float32(sr.Numerator)/float32(sr.Denominator))
	}
	return fmt.Sprintf("Undefined / undetermined : Denominator 0")
}

