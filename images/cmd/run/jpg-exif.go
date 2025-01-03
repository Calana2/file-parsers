package main

import (
	"flag"
	"fmt"
	"main/jpg"
)

func main() {
 metadata := flag.Bool("data",false,"")
 flag.BoolVar(metadata,"d",false,"Print metadata")
 flag.Parse()
 if len(flag.Args()) != 1 {
   flag.Usage()
 }

 jpg, err := jpg.New(flag.Arg(0))
 if err != nil {
  fmt.Println(err)
  return
 }

 if *metadata {
   jpg.ShowMetadata()
 }
}
