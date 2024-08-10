package svgan

import (
	"fmt"
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/JoshVarga/svgparser"
)

type SvgCheckResult struct {
	SvgHeight          string   `json:"svgheight"`
	SvgWidth           string   `json:"svgwidth"`
	ViewBox            string   `json:"viewbox"`
	Namespace          bool     `json:"namespace"`
	Namespaces         []string `json:"namespaces"`
	TextCount          int      `json:"text_count"`
	ForeignObjectCount int      `json:"foreignobject_count"`
	ImageCount         int      `json:"image_count"`
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
		SvgWidth:  rootElement.Attributes["width"],
		SvgHeight: rootElement.Attributes["height"],
		ViewBox:   rootElement.Attributes["viewBox"],
	}
	/*
		if svgNamespace {

			namespaces, _ := shared.GetNamespaces(bytes.NewReader(raw))
			//fmt.Printf("namespace: %v", namespaces)
			f.RecordResult("svgNamespace", namespaces.Default == "http://www.w3.org/2000/svg", map[string]interface{}{
				"namespace": namespaces.Default,
			})

			if len(svgNamespaces) == 0 {
				f.RecordResult("svgNoAdditionalNamespaces", len(namespaces.Additional) == 0, map[string]interface{}{
					"namespaces": namespaces.Additional,
				})
			} else if len(svgNamespaces) == 1 && svgNamespaces[0] == "*" {
				// no check
			} else {
				for key, value := range namespaces.Additional {
					_, keyExists := svgNamespaceSet[key]
					_, valueExists := svgNamespaceSet[value]
					f.RecordResult("svgAdditionalNamespaces", keyExists || valueExists, map[string]interface{}{
						"namespaceUrl":   value,
						"namespaceValue": key,
					})
				}
			}
		}
	*/

	textNodes := rootElement.FindAll("text")
	result.TextCount = len(textNodes)

	foNodes := rootElement.FindAll("foreignObject")
	result.ForeignObjectCount = len(foNodes)

	imageNodes := rootElement.FindAll("image")
	result.ImageCount = len(imageNodes)

	return result, nil
}
