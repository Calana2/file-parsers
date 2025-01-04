package main

import (
	"flag"
	"fmt"
	"main/wav"
	"os"
	"os/signal"
	"syscall"
)

func cleanup() {
 fmt.Print("\033[?25h")
}

func main() {
  sigchan := make(chan os.Signal, 1)
  signal.Notify(sigchan,syscall.SIGTERM,syscall.SIGINT)
  go func(){
   sig := <- sigchan
   fmt.Println()
   fmt.Printf("%s received.\n", sig)
   cleanup()
   os.Exit(0)
  }()

	metadata := flag.Bool("data", false, "")
	flag.BoolVar(metadata, "d", false, "Print metadata")

	play := flag.Bool("play", false, "")
	flag.BoolVar(play, "p", false, "Play audio")

	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		return
	}

	wav, err := wav.New(flag.Arg(0))
	if err != nil {
		fmt.Println("Error parsing the file: ", err)
    os.Exit(1)
  }

	if *metadata {
	 wav.PrintMetadata() 
  }
	if *play {
		err := wav.PlayAudio()
		if err != nil {
			fmt.Println("Error playing the audio: ", err)
      os.Exit(1)
		}
  }
}
