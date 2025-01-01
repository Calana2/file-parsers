package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: program [image_path]")
		return
	}
	var filename = os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

  var img image.Image
  switch strings.Split(filename,".")[1] {
   case "png":
    img,err = png.Decode(file)
    if err != nil { panic(err) }
   case "jpg","jpeg":
    img,err = jpeg.Decode(file)
    if err != nil { panic(err) }
   default:
    return
  }

	bounds := img.Bounds()
	bwImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, _ := c.RGBA()
			gray := uint8(int(0.299*float32(r) + 0.587*float32(g) + 0.114*float32(b)) % 256)
			bwImg.Set(x, y, color.Gray{Y: gray})
		}
	}

	outFile, err := os.Create("bw_" + filename)
  fmt.Println("bw_"+filename+" created.")
	if err != nil {
		panic(err)
	}
	err = jpeg.Encode(outFile, bwImg, nil)
	if err != nil {
		panic(err)
	}
}
