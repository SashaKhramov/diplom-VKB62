package classifiers

import (
	"regexp"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	//"github.com/mushorg/go-dpi/types"
)

// HTTPClassifier struct
type HTTPClassifier struct{}

var httpVerbs = []string{
	"OPTIONS",
	"GET",
	"HEAD",
	"POST",
	"PUT",
	"DELETE",
	"TRACE",
	"CONNECT",
}

var regex *regexp.Regexp

func init() {
	var regexStr = "^(" + strings.Join(httpVerbs, "|") + ") [^\\s]+ " +
		"HTTP/[12](.[01])?\r\n(.*\r\n)*\r\n"
	// regex should match the first line of all HTTP requests
	regex, _ = regexp.Compile(regexStr)
}

// HeuristicClassify for HTTPClassifier
func (classifier HTTPClassifier) HeuristicClassify(flow *types.Flow) bool {
	return checkFlowLayer(flow, layers.LayerTypeTCP, func(layer gopacket.Layer) bool {
		payload := layer.LayerPayload()
		return regex.Match(payload)
	})
}

// GetProtocol returns the corresponding protocol
func (classifier HTTPClassifier) GetProtocol() types.Protocol {
	return types.HTTP
}
