### WAV metadata viewer and player

` go run wav-parser.go -d -p duvet.wav ` 
```
Name                                              : duvet.wav
File Type Block ID                                : RIFF (Resource Interchange File Format)
File Size                                         : 1.46 MB (1528134 Bytes)
File Format ID                                    : WAVE
Format Block ID                                   : fmt 
Block Size                                        : 16 Bytes
Audio Format                                      : 1 [PCM (Pulse Code Modulation) – Uncompressed format]
Number of Channels                                : 2 [Stereo]
Frequence                                         : 44100 Hz
Bytes per second                                  : 172.27 KB
Bytes per block                                   : 4 Bytes
Bits per sample                                   : 16
LIST Block ID                                     : LIST
LIST Block Size                                   : 26 Bytes
Type of LIST Block                                : INFO
 Chunk                                            : ISFT
 Data Size                                        : 13
 Data                                             : Lavf61.7.100
Data Block ID                                     : data
Data Size                                         : 1.46 MB (1528056 Bytes)
Estimated audio duration                          : 8.662449 seconds

Playing duvet.wav
┌────────────────────────────────────────────────────────────────────────────────────────────────────┐
│███████████████████████████████████████████████████████████████████████████████                     │
^C───────────────────────────────────────────────────────────────────────────────────────────────────┘
interrupt received.
```
