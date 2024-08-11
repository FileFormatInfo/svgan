package svgan

import (
	"fmt"
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/JoshVarga/svgparser"
	"github.com/mazznoer/csscolorparser"
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/css"
)

type SvgCheckResult struct {
	SvgHeight          string         `json:"svgheight"`
	SvgWidth           string         `json:"svgwidth"`
	ViewBox            string         `json:"viewbox"`
	Namespace          bool           `json:"namespace"`
	Namespaces         []string       `json:"namespaces"`
	TextCount          int            `json:"text_count"`
	ForeignObjectCount int            `json:"foreignobject_count"`
	ImageCount         int            `json:"image_count"`
	TagCountMap        map[string]int `json:"tag_counts"`
	Colors             []string       `json:"colors"`
}

func SvgCheck(logger *slog.Logger, raw []byte) (*SvgCheckResult, error) {

	// check if raw is valid utf-8
	if !utf8.Valid(raw) {
		return nil, fmt.Errorf("not valid utf-8")
	}

	text := string(raw)

	rootElement, parseErr := svgparser.Parse(strings.NewReader(text), false)
	if parseErr != nil {
		logger.Error("Unable to parse", "error", parseErr)
		return nil, parseErr
	}

	result := &SvgCheckResult{
		SvgWidth:    rootElement.Attributes["width"],
		SvgHeight:   rootElement.Attributes["height"],
		ViewBox:     rootElement.Attributes["viewBox"],
		TagCountMap: make(map[string]int),
	}

	// walk all elements and count them
	colorMap := make(map[string]bool)
	walkElements(rootElement, func(node *svgparser.Element) {
		result.TagCountMap[node.Name]++

		processColor(logger, node.Attributes["fill"], colorMap)
		processColor(logger, node.Attributes["stroke"], colorMap)
		processColor(logger, node.Attributes["stop-color"], colorMap)
		processColor(logger, node.Attributes["color"], colorMap)

		processStyle(logger, node.Attributes["style"], colorMap)
	})
	i := 0
	result.Colors = make([]string, len(colorMap))
	for k := range colorMap {
		result.Colors[i] = k
		i++
	}
	result.Namespaces = getNamespaces(rootElement)

	textNodes := rootElement.FindAll("text")
	result.TextCount = len(textNodes)

	foNodes := rootElement.FindAll("foreignObject")
	result.ForeignObjectCount = len(foNodes)

	imageNodes := rootElement.FindAll("image")
	result.ImageCount = len(imageNodes)

	return result, nil
}

func walkElements(node *svgparser.Element, f func(*svgparser.Element)) {
	f(node)
	for _, child := range node.Children {
		walkElements(child, f)
	}
}

func processColor(logger *slog.Logger, rawColor string, colorMap map[string]bool) {
	if rawColor == "" || rawColor == "none" {
		return
	}
	color, colorErr := csscolorparser.Parse(rawColor)
	if colorErr != nil {
		logger.Error("Unable to parse color", "err", colorErr, "color", rawColor)
	} else {
		colorMap[color.HexString()] = true
	}
}

func processStyle(logger *slog.Logger, rawStyle string, colorMap map[string]bool) {
	if rawStyle == "" {
		return
	}

	input := parse.NewInputString(rawStyle)
	parser := css.NewParser(input, true)
	for {
		gt, _, data := parser.Next()
		if gt == css.ErrorGrammar {
			logger.Error("Unable to parse style", "style", rawStyle, "data", data)
			break
		}
		if gt == css.DeclarationGrammar {
			attr := string(data)
			if attr == "fill" || attr == "stroke" || attr == "color" || attr == "stop-color" {
				for _, val := range parser.Values() {
					processColor(logger, string(val.Data), colorMap)
				}
			}
		}
	}
}

// this doesn't work.  Need to use the xml.decoder stuff
func getNamespaces(node *svgparser.Element) []string {
	namespaces := make([]string, 0)
	for key := range node.Attributes {
		fmt.Println(key)
		if strings.Contains(key, "xmlns") {
			namespaces = append(namespaces, key)
		}
	}
	return namespaces
}
