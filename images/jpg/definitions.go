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
