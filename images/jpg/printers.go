package jpg

import "fmt"

func showSegment(s Segment) {
	fmt.Printf(" %-31s: \033[1;31m%s (0x%X)\033[1;0m\n",
		"Segment", "Unknown",s.Marker)
	fmt.Printf("  %-30s: %d Bytes\n",
		"Length", s.Length)
}

func showAPP(s APPSegment) {
	num := s.Marker[1] - 0xe0
	fmt.Printf(" \033[1;34m%-31s: %s %d\033[1;0m\n",
		"Segment", "Application", num)
	fmt.Printf("  %-30s: %d Bytes\n",
		"Length", s.Length)
	fmt.Printf("  %-30s: %s\n",
		"Identifier", s.Identifier[:len(s.Identifier)-1])
	fmt.Printf("  %-30s: %s\n",
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
	fmt.Printf("  %-30s: %d [%s]\n",
		"Density units", s.DensityUnits, desc)
	fmt.Printf("  %-30s: %d\n",
		"X Density", s.Xdensity)
	fmt.Printf("  %-30s: %d\n",
		"Y Density", s.Ydensity)
	fmt.Printf("  %-30s: %dx%d\n",
		"Thumbnail", s.XThumbnail, s.YThumbnail)
	if len(s.ThumbnailData) == 0 {
		fmt.Printf("  %-30s: %s\n",
			"Thumbnail Size", "None")
	} else {
		fmt.Printf("  %-30s: %d Bytes\n",
			"Thumbnail Size", len(s.ThumbnailData))
	}
}

func showEXIF(s EXIFSegment) {
	num := s.Marker[1] - 0xe0
	fmt.Printf(" \033[1;34m%-31s: %s %d\033[1;0m\n",
		"Segment", "Application", num)
	fmt.Printf("  %-30s: %d Bytes\n",
		"Length", s.Length)
	fmt.Printf("  %-30s: %s\n",
		"Identifier", s.Identifier)
	fmt.Printf("  %-30s: %s\n",
		"Format", "TIFF")
	var endianness string
	switch string(s.TIFFHeader.Alignment[:]) {
	case "II":
		endianness = "little-endian"
	case "MM":
		endianness = "big-endian"
	}
	fmt.Printf("  %-30s: %s [%s]\n",
		"Endianness", s.TIFFHeader.Alignment, endianness)
	for idx, IFD := range s.IFDs {
		fmt.Printf("  %-30s: %d (%s)\n",
			"Image File Directory", idx, IFDType[idx])
		fmt.Printf("  %-30s: %d \n",
			"Number of Entries", IFD.EntriesNum)
		for _, entry := range IFD.Entries {
      switch idx {
       case 0:
			  fmt.Printf("   %-29s: ",findIFD0Tag(entry.Tag).Name) 
       case 1:
			  fmt.Printf("   %-29s: ",findIFD1Tag(entry.Tag).Name) 
       case 2:
			  fmt.Printf("   %-29s: ",findSubIFDTag(entry.Tag).Name) 
      }
      for _,component := range entry.Data {
        fmt.Printf("%v ",component)
      }
      fmt.Printf("\n")
		}
	}
}

func showICC(s ICCSegment) {
	num := s.Marker[1] - 0xe0
	fmt.Printf(" \033[1;34m%-31s: %s %d\033[1;0m\n",
		"Segment", "Application", num)
	fmt.Printf("  %-30s: %d Bytes\n",
		"Length", s.Length)
	fmt.Printf("  %-30s: %s\n",
		"Identifier", s.Identifier)
}

func showCOM(s COMSegment) {
	fmt.Printf(" \033[1;34m%-31s: %s\033[1;0m\n",
		"Segment", "Comment")
	fmt.Printf("  %-30s: %d Bytes\n",
		"Length", s.Length)
	fmt.Printf("  %-30s: %s\n",
		"Comment", s.Data)
}

func showSOF(s SOFSegment, codingAlg string) {
	num := s.Marker[1] - 0xc0
	fmt.Printf(" \033[1;34m%-31s: %s %d\033[1;0m\n",
		"Segment", "Start Of Frame", num)
	fmt.Printf("  %-30s: %d Bytes\n",
		"Length", s.Length)
	fmt.Printf("  %-30s: %d\n",
		"Precision (bits per component)", s.Precision)
	fmt.Printf("  %-30s: %d\n",
		"Image Width", s.Samples_line)
	fmt.Printf("  %-30s: %d\n",
		"Image Height", s.LineNB)
	fmt.Printf("  %-30s: %dx%d\n",
		"Image Size", s.Samples_line, s.LineNB)
	var descImage string
	switch s.Components {
	case 1:
		descImage = "Grayscale"
	default:
		descImage = "Colorized"
	}
	fmt.Printf("  %-30s: %d [%s]\n",
		"Number of components", s.Components, descImage)
	fmt.Printf("  %-30s: %s, %s\n",
		"Encoding Process", DCTEncodingTypes[num], codingAlg)
	for i, v := range s.SOFComponents {
		horizontalSamplingFactor := (v.Sampling_X__Y >> 4) & 0xf
		verticalSamplingFactor := v.Sampling_X__Y & 0xf
		fmt.Printf("  %-30s: %dx%d\n",
			"Component "+ImageComponents[i]+" sampling factor",
			horizontalSamplingFactor, verticalSamplingFactor)
		fmt.Printf("  %-30s: %d\n",
			"Quantization Table Index", v.Quantization)
	}
}

func showDQT(s DQTSegment) {
	fmt.Printf(" \033[1;34m%-31s: %s\033[1;0m\n",
		"Segment", "Define Quantization Table")
	fmt.Printf("  %-30s: %d Bytes\n",
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
	fmt.Printf("  %-30s: %d [%s]\n",
		"Precision", higherPart, component)
	fmt.Printf("  %-30s: %d\n",
		"Table Index", lowerPart)
}

func showDHT(s DHTSegment) {
	fmt.Printf(" \033[1;34m%-31s: %s\033[1;0m\n",
		"Segment", "Define Huffman Table")
	fmt.Printf("  %-30s: %d Bytes\n",
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
	fmt.Printf("  %-30s: %d [%s]\n",
		"Class", classN, class)
	fmt.Printf("  %-30s: %d\n",
		"Table Index", idx)
	var codes string
	for i, v := range s.Bit_Codes {
		if i == 9 {
			codes += fmt.Sprintf("\n   %-30s ", "")
		}
		codes += fmt.Sprintf("[%d : %d] ", i+1, v)
	}
	fmt.Printf("  %-30s: %s\n",
		"Huffman Codes", codes)
}

func showSOS(s SOSSegment) {
	fmt.Printf(" \033[1;34m%-31s: %s\033[1;0m\n",
		"Segment", "Start Of Scan")
	fmt.Printf("  %-30s: %d Bytes\n",
		"Length", s.Length)
	fmt.Printf("  %-30s: %d\n",
		"Number of Components", s.Components)
	yDC := s.Y_AC__DC & 0xf
	yAC := (s.Y_AC__DC >> 4) & 0xf
	fmt.Printf("  %-30s: %d [DC Table Index] %d [AC Table Index]\n",
		"Luminance(Y)", yDC, yAC)
	if s.Components > 1 {
		CbDC := s.Cb_AC__DC & 0xf
		CbAC := (s.Cb_AC__DC >> 4) & 0xf
		fmt.Printf("  %-30s: %d [DC Table Index] %d [AC Table Index]\n",
			"Crominance(Cb)", CbDC, CbAC)
		if s.Components > 2 {
			CrDC := s.Cb_AC__DC & 0xf
			CrAC := (s.Cb_AC__DC >> 4) & 0xf
			fmt.Printf("  %-30s: %d [DC Table Index] %d [AC Table Index]\n",
				"Crominance(Cr)", CrDC, CrAC)
		}
	}
	fmt.Printf("  %-30s: %d\n",
		"Start of spectral selection", s.SS_Start)
	fmt.Printf("  %-30s: %d\n",
		"End of spectral selection", s.SS_End)
	fmt.Printf("  %-30s: %d\n",
		"Sucesive approximation bits", s.Sucessive_approx)
}
