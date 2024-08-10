package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	svgan "github.com/FileFormatInfo/svgan/lib"
)

func main() {
	var logger = slog.New(slog.NewTextHandler(os.Stderr, nil))

	if len(os.Args) < 2 {
		logger.Error("no files specified")
		os.Exit(1)
	}

	for _, fileName := range os.Args[1:] {
		logger.Info("processing file", "filename", fileName)
		results, err := svgCheckFile(logger, fileName)
		if err != nil {
			logger.Error("unable to process file", "filename", fileName, "err", err)
			os.Exit(1)
		}

		jsonData, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			logger.Error("unable to convert results to JSON", "filename", fileName, "err", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stdout, "%s\n", string(jsonData))
	}
}

func svgCheckFile(logger *slog.Logger, fileName string) (*svgan.SvgCheckResult, error) {

	raw, readErr := os.ReadFile(fileName)
	if readErr != nil {
		logger.Error("Unable to read file", "fileName", fileName, "error", readErr)
		return nil, readErr
	}

	return svgan.SvgCheck(logger, raw)
}
