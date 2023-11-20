package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"log/slog"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"github.com/go-audio/wav"
)

func processSample(ctx context.Context, model whisper.Model, path string) (text string, err error) {
	slog.Info("instantiating model context")
	wCtx, err := model.NewContext()
	if err != nil {
		return
	}
	slog.Info("loading sample", "path", path)
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	decoder := wav.NewDecoder(file)
	buffer, err := decoder.FullPCMBuffer()
	f32Buffer := buffer.AsFloat32Buffer()

	slog.Info("processing sample", "path", path)
	err = wCtx.Process(f32Buffer.Data, nil, nil)
	if err != nil {
		return
	}
	for {
		segment, err := wCtx.NextSegment()
		if err == io.EOF {
			break
		}
		if err != nil {
			return text, err
		}

		slog.Info(
			"processed segment",
			"start", segment.Start,
			"end", segment.End,
			"text", segment.Text,
		)
		text = strings.TrimSpace(fmt.Sprintf("%s %s", text, segment.Text))
	}
	return
}

func write(ctx context.Context, records [][]string, outPath string) error {
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	outWriter := csv.NewWriter(outFile)
	err = outWriter.WriteAll(records)
	if err != nil {
		return err
	}
	outWriter.Flush()
	if err := outWriter.Error(); err != nil {
		return err
	}
	return nil
}

func run(ctx context.Context, modelPath string, samplePaths []string, outPath string) error {
	slog.Info("loading whisper model", "path", modelPath)
	model, err := whisper.New(modelPath)
	if err != nil {
		return err
	}
	defer model.Close()

	records := make([][]string, len(samplePaths))
	for i := range records {
		records[i] = make([]string, 2)
	}

	for i, path := range samplePaths {
		text, err := processSample(ctx, model, path)
		if err != nil {
			return err
		}
		records[i] = []string{path, text}
	}

	err = write(ctx, records, outPath)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	ctx := context.Background()

	var modelPath string
	var csvPath string
	flag.StringVar(&modelPath, "model", "", "path to whisper model")
	flag.StringVar(&csvPath, "csv", "", "path to CSV")
	flag.Parse()
	samplePaths := flag.CommandLine.Args()

	err := run(ctx, modelPath, samplePaths, csvPath)
	if err != nil {
		slog.Error("", "error", err)
	}
}
