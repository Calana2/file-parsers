## JPG Metadata viewer and tools

`go run main_jpg.go -d image.jpeg `
```
Start of Information            
 File Name                      : image.jpeg
 Segment                        : Application 0
  Length                        : 16 Bytes
  Identifier                    : JFIF
  Version                       : 1.1
  Density units                 : 0 [Pixel Aspect Ratio]
  X Density                     : 1
  Y Density                     : 1
  Thumbnail                     : 0x0
  Thumbnail Size                : None
 Segment                        : Comment
  Length                        : 22 Bytes
  Comment                       : SGlkZGVuIGNvbW1lbnQ=
 Segment                        : Define Quantization Table
  Length                        : 132 Bytes
  Precision                     : 0 [8 bits]
  Table Index                   : 9
 Segment                        : Start Of Frame 0
  Length                        : 17 Bytes
  Precision (bits per component): 8
  Image Width                   : 217
  Image Height                  : 232
  Image Size                    : 217x232
  Number of components          : 3 [Colorized]
  Encoding Process              : Baseline DCT, Huffman coding
  Component Y sampling factor   : 2x2
  Quantization Table Index      : 0
  Component Cb sampling factor  : 1x1
  Quantization Table Index      : 1
  Component Cr sampling factor  : 1x1
  Quantization Table Index      : 1
 Segment                        : Define Huffman Table
  Length                        : 28 Bytes
  Class                         : 0 [DC]
  Table Index                   : 0
  Huffman Codes                 : [1 : 0] [2 : 2] [3 : 2] [4 : 3] [5 : 1] [6 : 1] [7 : 0] [8 : 0] [9 : 0] 
                                  [10 : 0] [11 : 0] [12 : 0] [13 : 0] [14 : 0] [15 : 0] [16 : 0] 
 Segment                        : Define Huffman Table
  Length                        : 56 Bytes
  Class                         : 1 [AC]
  Table Index                   : 0
  Huffman Codes                 : [1 : 0] [2 : 2] [3 : 1] [4 : 2] [5 : 4] [6 : 5] [7 : 2] [8 : 4] [9 : 4] 
                                  [10 : 6] [11 : 1] [12 : 5] [13 : 1] [14 : 0] [15 : 0] [16 : 0] 
 Segment                        : Define Huffman Table
  Length                        : 23 Bytes
  Class                         : 0 [DC]
  Table Index                   : 1
  Huffman Codes                 : [1 : 1] [2 : 1] [3 : 1] [4 : 1] [5 : 0] [6 : 0] [7 : 0] [8 : 0] [9 : 0] 
                                  [10 : 0] [11 : 0] [12 : 0] [13 : 0] [14 : 0] [15 : 0] [16 : 0] 
 Segment                        : Define Huffman Table
  Length                        : 31 Bytes
  Class                         : 1 [AC]
  Table Index                   : 1
  Huffman Codes                 : [1 : 1] [2 : 1] [3 : 1] [4 : 0] [5 : 2] [6 : 2] [7 : 2] [8 : 3] [9 : 0] 
                                  [10 : 0] [11 : 0] [12 : 0] [13 : 0] [14 : 0] [15 : 0] [16 : 0] 
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

`image2gray.go` --> Convert images to grayscale

`gifmaker.go`   --> Create gifs (not true colors)           
