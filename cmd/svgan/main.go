package main

import (
	"fmt"
	"log/slog"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file name")
		os.Exit(1)
	}

	fileName := os.Args[1]
	results, err := svgCheckFile(logger, fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("SVG info: %v\n", results)
}

func svgCheckFile(Logger *slog.Logger, fileName string) (*SvgInfoResult, error) {

	raw, readErr := os.ReadFile(fileName)
	if readErr != nil {
		Logger.Error("msg", "Unable to read file", "fileName", fileName, "error", readErr)
		return nil, readErr
	}
	text := string(raw)

	return svgCheckText(Logger, text, raw)
}
