package main

import (
	"flag"
	"log"
	"main/wavk"
	"os"
)

func main() {
 metadata := flag.Bool("data",false,"Print metadata")
 flag.BoolVar(metadata,"d",false,"Print metadata")

 flag.Parse()

 if len(flag.Args()) != 1 {
   flag.Usage()
   return 
 }

 file, err := os.Open(flag.Arg(0))
 if err != nil {
   log.Fatalln("Error reading the file: ",err)
 }
 defer file.Close()
 wav,err := wavk.New(file)
 if err != nil {
   log.Fatalln("Error parsing the file: ",err)
 }

 switch {
  case *metadata:
   wav.PrintMetadata()
   fallthrough
  default:
 }
}



