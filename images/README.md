## JPG Metadata viewer and tools

`cmd/tools/image2gray.go` --> Convert images to grayscale 

`cmd/tools/gifmaker.go`   --> Create gifs (not true colors)           

`cmd/run/jpg-exif.go`     --> JPG metadata viewer

`go run jpg-exif.go -d Kodak_CX7530.jpg`
```
Start of Information            
 File Name                      : Kodak_CX7530.jpg
 Segment                        : Application 0
  Length                        : 16 Bytes
  Identifier                    : JFIF
  Version                       : 1.1
  Density units                 : 1 [Pixels per inch (2.54 cm)]
  X Density                     : 72
  Y Density                     : 72
  Thumbnail                     : 0x0
  Thumbnail Size                : None
 Segment                        : Application 1
  Length                        : 2898 Bytes
  Identifier                    : Exif
  Format                        : TIFF
  Endianness                    : II [little-endian]
  Image File Directory          : 0 (Main image)
  Number of Entries             : 11 
   Make                         : EASTMAN KODAK COMPANY
   Model                        : KODAK CX7530 ZOOM DIGITAL CAMERA
   Orientation                  : Horizontal (normal)
   XResolution                  : 72
   YResolution                  : 72
   ResolutionUnit               : Inch
   Software                     : GIMP 2.4.5
   DateTime                     : 2008:07:31 10:39:26
   YCbCrPositioning             : Centered
   ExitOffset                   : 248
   GPSInfoOffset                : 816
  Image File Directory          : 1 (Thumbnail image)
  Number of Entries             : 3 
   Compression                  : JPEG (old-style)
   JpegIFOffset                 : 972
   JpegIFByteCount              : 1918
  Image File Directory          : 2 (IFD Formatted Data -- SubIFD)
  Number of Entries             : 35 
   ExposureTime                 : 0.004
   FNumber                      : 4.6
   ExposureProgram              : Program Normal
   ExifVersion                  : 0221
   DateTimeOriginal             : 2005:08:13 09:47:23
   DateTimeDigitized            : 2005:08:13 09:47:23
   ComponentConfiguration       : YCbCr
   ShutterSpeedValue            : 8
   ApertureValue                : 4.4
   ExposureBiasValue            : 0
   MaxApertureValue             : 4.4
   MeteringMode                 : Multi-Segment
   LightSource                  : Unknown
   Flash                        : Auto, Did not fire
   FocalLength                  : 16.8
   FlashPixVersion              : 0100
   ColorSpace                   : sRGB
   ExifImageWidth               : 1000
   ExifImageHeight              : 780
   ExifInteroperabilityOffset   : 786
   ExposureIndex                : 80
   SensingMethod                : One-chip color area
   FileSource                   : Cr
   SceneType                    : Y
   CustomRendered               : Normal
   ExposureMode                 : Auto
   WhiteBalance                 : Auto
   DigitalZoomRatio             : 0
   FocalLength35mmFormat        : 102
   SceneCaptureType             : Standard
   GainControl                  : None
   Contrast                     : Normal
   Saturation                   : Normal
   Sharpness                    : Normal
   SubjectDistanceRange         : 0
  Image File Directory          : 3 (GPS IFD)
  Number of Entries             : 5 
   GPSVersionID                 : 2200
   GPSLatitudeRef               : S
   GPSLatitude                  : 022.2780
   GPSLongitudeRef              : E
   GPSLongitude                 : 363.3850
   GPS Position                 : 2 deg 2' 16.68" S 36 deg 3' 23.10" E
 Segment                        : Define Quantization Table
  Length                        : 67 Bytes
  Precision                     : 0 [8 bits]
  Table Index                   : 8
 Segment                        : Define Quantization Table
  Length                        : 67 Bytes
  Precision                     : 0 [8 bits]
  Table Index                   : 8
 Segment                        : Start Of Frame 0
  Length                        : 17 Bytes
  Precision (bits per component): 8
  Image Width                   : 100
  Image Height                  : 78
  Image Size                    : 100x78
  Number of components          : 3 [Colorized]
  Encoding Process              : Baseline DCT, Huffman coding
  Component Y sampling factor   : 1x1
  Quantization Table Index      : 0
  Component Cb sampling factor  : 1x1
  Quantization Table Index      : 1
  Component Cr sampling factor  : 1x1
  Quantization Table Index      : 1
  Y Cb Cr Sub Sampling          : YCbCr4:4:4 (1 1)
 Segment                        : Define Huffman Table
  Length                        : 26 Bytes
  Class                         : 0 [DC]
  Table Index                   : 0
  Huffman Codes                 : [1 : 0] [2 : 3] [3 : 1] [4 : 1] [5 : 1] [6 : 1] [7 : 0] [8 : 0] [9 : 0] 
                                  [10 : 0] [11 : 0] [12 : 0] [13 : 0] [14 : 0] [15 : 0] [16 : 0] 
 Segment                        : Define Huffman Table
  Length                        : 53 Bytes
  Class                         : 1 [AC]
  Table Index                   : 0
  Huffman Codes                 : [1 : 0] [2 : 2] [3 : 1] [4 : 2] [5 : 5] [6 : 2] [7 : 4] [8 : 4] [9 : 4] 
                                  [10 : 5] [11 : 5] [12 : 0] [13 : 0] [14 : 0] [15 : 0] [16 : 0] 
 Segment                        : Define Huffman Table
  Length                        : 26 Bytes
  Class                         : 0 [DC]
  Table Index                   : 1
  Huffman Codes                 : [1 : 0] [2 : 3] [3 : 1] [4 : 1] [5 : 1] [6 : 1] [7 : 0] [8 : 0] [9 : 0] 
                                  [10 : 0] [11 : 0] [12 : 0] [13 : 0] [14 : 0] [15 : 0] [16 : 0] 
 Segment                        : Define Huffman Table
  Length                        : 39 Bytes
  Class                         : 1 [AC]
  Table Index                   : 1
  Huffman Codes                 : [1 : 0] [2 : 2] [3 : 2] [4 : 2] [5 : 2] [6 : 2] [7 : 1] [8 : 3] [9 : 5] 
                                  [10 : 1] [11 : 0] [12 : 0] [13 : 0] [14 : 0] [15 : 0] [16 : 0] 
 Segment                        : Start Of Scan
  Length                        : 12 Bytes
  Number of Components          : 3
  Luminance(Y)                  : 0 [DC Table Index] 0 [AC Table Index]
  Crominance(Cb)                : 1 [DC Table Index] 1 [AC Table Index]
  Crominance(Cr)                : 1 [DC Table Index] 1 [AC Table Index]
  Start of spectral selection   : 0
  End of spectral selection     : 63
  Sucesive approximation bits   : 0
End of Information
```

