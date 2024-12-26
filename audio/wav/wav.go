package wav

import (
	"encoding/binary"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"main/utils"
	"math"
	"os"
	"strings"
	"time"
)

type WAV struct {
	Name string
	// RIF
	FileTypeBlockID string
	FileSize        uint32
	FileFormatID    string
	// fmt
	FormatBlockID string
	BlockSize     uint32
	AudioFormat   uint16
	NumChannels   uint16
	Frequence     uint32
	BytePerSec    uint32
	BytePerBlock  uint16
	BitsPerSample uint16
	// LIST
	LISTBlockID   string
	LISTBlockSize uint32
	LISTBlockType string
	// IART
	IARTBlockID    string
	IARTDataSize   uint32
	IARTArtistName string
	// INAM
	INAMBlockID  string
	INAMDataSize uint32
	INAMSongName string
	// ISFT
	ISFTBlockID  string
	ISFTDataSize uint32
	ISFTSoftware string
	// data
	DataBlockID string
	DataSize    uint32
}

// Create a new WAV file.
func New(filepath string) (*WAV, error) {
	file, err := os.Open(filepath)
	if err != nil {
    return nil,err
	}
	defer file.Close()
	wav := &WAV{}
	wav.Name = file.Name()
	headers := make([]byte, 202)
	if _, err := file.Read(headers); err != nil {
		return nil, err
	}
	// RIFF
	wav.FileTypeBlockID = string(headers[0:4])
	wav.FileSize = binary.LittleEndian.Uint32(headers[4:8]) + 0x8
	wav.FileFormatID = string(headers[8:12])
	// fmt
	wav.FormatBlockID = string(headers[12:16])
	wav.BlockSize = binary.LittleEndian.Uint32(headers[16:20])
	wav.AudioFormat = binary.LittleEndian.Uint16(headers[20:22])
	wav.NumChannels = binary.LittleEndian.Uint16(headers[22:24])
	wav.Frequence = binary.LittleEndian.Uint32(headers[24:28])
	wav.BytePerSec = binary.LittleEndian.Uint32(headers[28:32])
	wav.BytePerBlock = binary.LittleEndian.Uint16(headers[32:34])
	wav.BitsPerSample = binary.LittleEndian.Uint16(headers[34:36])
	// optional chunks
	offset := 36
	chunk := string(headers[offset : offset+4])
	for chunk != "data" {
		switch chunk {
		// strange error, maybe because of some padding to memory alignment ?
		case "\x00dat","\x00\x00da","\x00\x00\x00d":
			offset++
		case "LIST":
			wav.LISTBlockID = string(headers[offset : offset+4])
			wav.LISTBlockSize = binary.LittleEndian.Uint32(headers[offset+4 : offset+8])
			wav.LISTBlockID = string(headers[offset+8 : offset+12])
			offset += 12
		case "IART":
			wav.IARTBlockID = string(headers[offset : offset+4])
			wav.IARTDataSize = binary.LittleEndian.Uint32(headers[offset+4 : offset+8])
			wav.IARTArtistName = string(headers[offset+8 : offset+8+int(wav.IARTDataSize)])
			offset += 8 + int(wav.IARTDataSize)
		case "INAM":
			wav.INAMBlockID = string(headers[offset : offset+4])
			wav.INAMDataSize = binary.LittleEndian.Uint32(headers[offset+4 : offset+8])
			wav.INAMSongName = string(headers[offset+8 : offset+8+int(wav.INAMDataSize)])
			offset += 8 + int(wav.INAMDataSize)
		case "ISFT":
			wav.ISFTBlockID = string(headers[offset : offset+4])
			wav.ISFTDataSize = binary.LittleEndian.Uint32(headers[offset+4 : offset+8])
			wav.ISFTSoftware = string(headers[offset+8 : offset+8+int(wav.ISFTDataSize)])
			offset += 8 + int(wav.ISFTDataSize)
 		case "IARL", "ICMS", "ICMT", "ICOP", "ICRD",
			"ICRP", "IDIM", "IDPI", "IENG", "IGNR",
			"IKEY", "ILGT", "IMED", "IPLT", "IPRD",
			"ISBJ", "ISRC", "ISRF", "ITCH":
			// Not yet
			offset += 8 + int(binary.LittleEndian.Uint32(headers[offset+4:offset+8]))
    default:
     return nil,fmt.Errorf("Unknown chunk ID") 
	}
		chunk = string(headers[offset : offset+4])
	}
	// data
	wav.DataBlockID = string(headers[offset : offset+4])
	wav.DataSize = binary.LittleEndian.Uint32(headers[offset+4 : offset+8])
	// verify
	if err = isWav(wav.Name, wav.FileTypeBlockID, wav.FileFormatID); err != nil {
		return nil, err
	}
	return wav, nil
}

// Formatted metadata display.
func (w *WAV) PrintMetadata() {
	// RIFF
	fmt.Printf("%-50s: %s\n",
		"Name", w.Name)
	fmt.Printf("%-50s: %s (Resource Interchange File Format)\n",
		"File Type Block ID", w.FileTypeBlockID)
	fmt.Printf("%-50s: %s (%d Bytes)\n",
		"File Size", utils.ConvertSizeUint32(w.FileSize), w.FileSize)
	// fmt
	fmt.Printf("%-50s: %s\n",
		"File Format ID", w.FileFormatID)
	fmt.Printf("%-50s: %s\n",
		"Format Block ID", w.FormatBlockID)
	fmt.Printf("%-50s: %s\n",
		"Block Size", utils.ConvertSizeUint32(w.BlockSize))
	switch w.AudioFormat {
	case 0x1:
		fmt.Printf("%-50s: %d [PCM (Pulse Code Modulation) – Uncompressed format]\n",
			"Audio Format", w.AudioFormat)
	case 0x3:
		fmt.Printf("%-50s: %d [IEEE Float – Floating Point Audio]\n",
			"Audio Format", w.AudioFormat)
	case 0x6:
		fmt.Printf("%-50s: %d [A-law – A-law compression algorithm]",
			"Audio Format", w.AudioFormat)
	case 0x7, 0xFFFE:
		fmt.Printf("%-50s: %d [µ-law – µ-law compression algorithm]",
			"Audio Format", w.AudioFormat)
	default:
		fmt.Printf("%-50s: %d [Corrupted or unknown format]",
			"Audio Format", w.AudioFormat)
	}
	if w.NumChannels == 1 {
		fmt.Printf("%-50s: %d [Mono]\n",
			"Number of Channels", w.NumChannels)
	} else {
		fmt.Printf("%-50s: %d [Stereo]\n",
			"Number of Channels", w.NumChannels)
	}
	fmt.Printf("%-50s: %d Hz\n",
		"Frequence", w.Frequence)
	fmt.Printf("%-50s: %s\n",
		"Bytes per second", utils.ConvertSizeUint32(w.BytePerSec))
	fmt.Printf("%-50s: %s\n",
		"Bytes per block", utils.ConvertSizeUint16(w.BytePerBlock))
	fmt.Printf("%-50s: %d\n",
		"Bits per sample", w.BitsPerSample)
	// LIST
	if w.LISTBlockID != "" {
		fmt.Printf("%-50s: %s\n",
			"LIST Block ID", w.LISTBlockID)
		fmt.Printf("%-50s: %d Bytes\n",
			"LIST Block Size", w.LISTBlockSize)
	}
	// INFO
	if w.LISTBlockID != "" {
		fmt.Printf("%-50s: %s\n",
			"Type of LIST Block", w.LISTBlockID)
	}
	// IART
	if w.IARTBlockID != "" {
		fmt.Printf("%-50s: %s\n",
			"IART Block ID", w.IARTBlockID)
		fmt.Printf("%-50s: %d Bytes\n",
			"IART Data Size", w.IARTDataSize)
		fmt.Printf("%-50s: %s\n",
			"IART Artist Name", w.IARTArtistName)
	}
	// INAM
	if w.INAMBlockID != "" {
		fmt.Printf("%-50s: %s\n",
			"INAM Block ID", w.INAMBlockID)
		fmt.Printf("%-50s: %d Bytes\n",
			"INAM Data Size", w.INAMDataSize)
		fmt.Printf("%-50s: %s\n",
			"INAM Song Name", w.INAMSongName)
	}
	// ISFT
	if w.ISFTBlockID != "" {
		fmt.Printf("%-50s: %s\n",
			"ISFT Block ID", w.ISFTBlockID)
		fmt.Printf("%-50s: %d Bytes\n",
			"ISFT Data Size", w.ISFTDataSize)
		fmt.Printf("%-50s: %s\n",
			"ISFT Software", w.ISFTSoftware)
	}
	// data
	fmt.Printf("%-50s: %s\n",
		"Data Block ID", w.DataBlockID)
	fmt.Printf("%-50s: %s (%d Bytes)\n",
		"Data Size", utils.ConvertSizeUint32(w.DataSize), w.DataSize)
	// Extra
	fmt.Printf("%-50s: %f seconds\n",
		"Estimated audio duration", utils.CalculateAudioDuration(w.DataSize, w.BytePerSec))
	fmt.Println()
}

// Verify if it is a WAV file.
func isWav(name string, riff string, wave string) error {
	if substr := name[strings.LastIndex(name, ".")+1:]; substr != "wav" {
		print(substr)
		return fmt.Errorf("File %s does not have the .wav extension.", name)
	}
	if riff != "RIFF" || wave != "WAVE" {
		return fmt.Errorf("File %s is not a valid WAV file.\n"+
			"File Type Block ID = %v (expected %v)\n"+
			"File Format ID = %v (expected %v)", name,
			[]byte(riff), []byte("RIFF"),
			[]byte(wave), []byte("WAVE"))
	}
	return nil
}

// Reproduces the wav file
func (w *WAV) PlayAudio() error {
	file, _ := os.Open(w.Name)
	defer file.Close()
	streamer, format, err := wav.Decode(file)
	if err != nil {
		return fmt.Errorf("Error decoding the WAV: %v", err)
	}
	defer streamer.Close()
	duration_seconds := float64(streamer.Len()) / float64(w.Frequence)
	percentage := int64(duration_seconds) * int64(math.Pow(10, 9)) / 100
	go func() {
		total := 100
		upper_border := "\u250c"
		empty_middle := "\u2502"
		lower_border := "\u2514"
		for i := 0; i < 100; i++ {
			upper_border += "\u2500"
			lower_border += "\u2500"
			empty_middle += " "
		}
		upper_border += "\u2510"
		lower_border += "\u2518"
		empty_middle += "\u2502"
		fmt.Printf("Playing %s\n", w.Name)
		fmt.Println(upper_border)
		fmt.Println(empty_middle)
		fmt.Print(lower_border + "\033[?25l")
		for {
			for i := 0; i <= total; i++ {
				fmt.Printf("\033[F\r\u2502")
				for j := 0; j < total; j++ {
					if j < i {
						fmt.Printf("\u2588")
					} else {
						fmt.Print(" ")
					}
				}
				fmt.Printf("\u2502\n")
				time.Sleep(time.Duration(percentage))
			}
			break
		}
	}()
	// Initialize the speaker
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	// Play
	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() { done <- true })))
	// Wait until it ends
	<-done
	fmt.Println("\033[?25h")
	return nil
}
