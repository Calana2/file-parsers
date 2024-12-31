package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: program [image1] [image2]... [output.gif]")
		return
	}

  // Decoding and converting the images
	var imgArray []*image.Paletted
  var delays []int
	for _, f := range os.Args[1:len(os.Args)-1] {
		file, err := os.Open(f)
    defer file.Close()
		if err != nil {
      fmt.Printf("Error opening %s: %s", f,err)
			continue
		}
		var img image.Image
		switch strings.Split(f, ".")[1] {
		case "png":
			img, err = png.Decode(file)
			if err != nil {
				panic(err)
			}
		case "jpg", "jpeg":
			img, err = jpeg.Decode(file)
			if err != nil {
				panic(err)
			}
		default:
      fmt.Printf("Error decoding %s: %s", f, err)
			continue
		}
		imgArray = append(imgArray, convertToPalleted(img))
    delays = append(delays, 100)
	}

  // Creating the gif
	outputFile, err := os.Create(os.Args[len(os.Args)-1])
	if err != nil {
    fmt.Printf("Error creating %s: %s", os.Args[len(os.Args)-1],err)
		return
	}
  defer outputFile.Close()
  err = gif.EncodeAll(outputFile, &gif.GIF{Image: imgArray, Delay: delays})
  if err != nil {
    fmt.Printf("Error creating the GIF: %s",err)
    return
  }
}

func convertToPalleted(img image.Image) *image.Paletted {
      palette := make(color.Palette, 256)

    for i := 0; i < 256; i++ {
        r := uint8((i * 15) % 256)         
        g := uint8((i * 13) % 256)        
        b := uint8((i * 17) % 256)       
        palette[i] = color.RGBA{r, g, b, 255}
    }

  palletedImg := image.NewPaletted(img.Bounds(),palette)
  draw.Draw(palletedImg,img.Bounds(),img, img.Bounds().Min,draw.Over)
	return palletedImg
}
