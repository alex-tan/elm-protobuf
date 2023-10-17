package generationparams

import (
	"fmt"
	"strings"
)

type Parameters struct {
	Version              bool
	Debug                bool
	RemoveDeprecated     bool
	GenerateForwardCache bool
	GenerateForwardGraph bool
	GenerateForwardIds   bool
	GenerateApi          bool
}

func ParseParameters(input *string) (Parameters, error) {
	var result Parameters
	var err error

	if input == nil {
		return result, nil
	}

	for _, i := range strings.Split(*input, ",") {
		switch i {
		case "remove-deprecated":
			result.RemoveDeprecated = true
		case "debug":
			result.Debug = true
		case "generate-forward-cache":
			result.GenerateForwardCache = true
		case "generate-forward-ids":
			result.GenerateForwardIds = true
		case "generate-forward-graph":
			result.GenerateForwardGraph = true
		case "generate-api":
			result.GenerateApi = true
		default:
			err = fmt.Errorf("unknown parameter: \"%s\"", i)
		}
	}

	return result, err
}
