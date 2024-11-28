### WAV metadata viewer

` go run wav/cmd/run.go -d file.wav`
```
'RIFF' chunk descriptor
-----------------------
File Type Block ID: RIFF (Resource Interchange File Format)
File Size: 1.00 MB (1048558 Bytes)
File Format ID: WAVE


'fmt' sub-chunk
----------------
Format Block ID: fmt 
Block Size: 24 Bytes
Audio Format: 1 [PCM (Pulse Code Modulation) – Uncompressed format]
Number of Channels: 2 [Stereo]
Frequence: 44100 Hz
Bytes per second: 172.27 KB
Bytes per block: 4 Bytes
Bits per sample: 16


'data' sub-chunk
----------------
Data Block ID: data
Data Size: 1.00 MB (1048376 Bytes) 


Estimated audio duration: 5.94 seconds
```
