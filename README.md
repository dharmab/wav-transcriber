# WAV Transcriber

This is a small program for transcribing the content of WAV files. You pass it a Whisper model, a list of WAV files, and a CSV path. It writes a CSV file containing the filenames and the transcription of the file.

I made this to automatically transcribe a rip of the WAV files from Ace Combat 5 and Ace Combat Zero that someone posted onto the internet years ago. It's close enough that you can use it to find the file containing the line you need. However, it needs manual correction for use as in-game/in-video subtitles.

## Usage

To use this, you'll need the [Go programming language tools](https://go.dev)

Download a ggml whisper.cpp model from https://huggingface.co/ggerganov/whisper.cpp. I found that the medium model worked pretty well, but you should experiment on your data.

Then run:

```
C_INCLUDE_PATH=third_party/whisper.cpp/ \
LIBRARY_PATH=third_party/whisper.cpp/ \
go run cmd/main.go \
--model $path_to_model \
--csv out.csv \
file1.wav file2.wav file3.wav
```

`C_INCLUDE_PATH` and `LIBRARY_PATH` should be a directory containing `whisper.h` and `libwhisper.a`, respectively. `--model` should be the whisper.cpp model downloaded from huggingface. `--csv` is the CSV file that will be created. The remaining arguments are the WAV files to transcribe.

A sample script is included in `scripts`.
