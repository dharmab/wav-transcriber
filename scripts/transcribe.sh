#!/bin/sh

export C_INCLUDE_PATH=third_party/whisper.cpp/
export LIBRARY_PATH=third_party/whisper.cpp/

go run cmd/transcribe/main.go \
  --model models/ggml-model-whisper-medium.en-q5_0.bin \
  --csv out.csv \
  "$@"

