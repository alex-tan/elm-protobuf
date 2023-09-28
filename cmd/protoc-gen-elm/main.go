package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/thematthopkins/elm-protobuf/pkg/elmpb"
	"github.com/thematthopkins/elm-protobuf/pkg/forwardcache"
	"github.com/thematthopkins/elm-protobuf/pkg/forwardgraph"
	"github.com/thematthopkins/elm-protobuf/pkg/forwardids"
	"github.com/thematthopkins/elm-protobuf/pkg/generationparams"
	"github.com/thematthopkins/elm-protobuf/pkg/parsepb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

const version = "0.0.2"
const docUrl = "https://github.com/thematthopkins/elm-protobuf"

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Fprintf(os.Stdout, "%v %v\n", filepath.Base(os.Args[0]), version)
		os.Exit(0)
	}
	if len(os.Args) == 2 && os.Args[1] == "--help" {
		fmt.Fprintf(os.Stdout, "See "+docUrl+" for usage information.\n")
		os.Exit(0)
	}

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Could not read request from STDIN: %v", err)
	}

	req := &pluginpb.CodeGeneratorRequest{}

	err = proto.Unmarshal(data, req)
	if err != nil {
		log.Fatalf("Could not unmarshal request: %v", err)
	}

	parameters, err := generationparams.ParseParameters(req.Parameter)
	if err != nil {
		log.Fatalf("Failed to parse parameters: %v", err)
	}

	if parameters.Debug {
		// Remove useless source code data.
		for _, inFile := range req.GetProtoFile() {
			inFile.SourceCodeInfo = nil
		}

		result, err := proto.Marshal(req)
		if err != nil {
			log.Fatalf("Failed to marshal request: %v", err)
		}

		log.Printf("Input data: %s", result)
	}

	plugins := (uint64)(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	resp := &pluginpb.CodeGeneratorResponse{
		SupportedFeatures: &plugins,
	}

	for _, inFile := range req.GetProtoFile() {
		log.Printf("Processing file %s", inFile.GetName())
		// Well Known Types.
		if elmpb.ExcludedFiles[inFile.GetName()] {
			log.Printf("Skipping well known type")
			continue
		}

		file, err := elmpb.Generate(inFile, parameters)
		if err != nil {
			log.Fatalf("Could not template elmpb: %v", err)
		}

		resp.File = append(resp.File, file)
		messages := parsepb.Messages([]string{}, inFile.GetMessageType(), parameters)
		if parameters.GenerateForwardCache {
			file, err := forwardcache.Generate(inFile, messages)
			if err != nil {
				log.Fatalf("Could not template forwardcache: %v", err)
			}

			resp.File = append(resp.File, file)
		}
		if parameters.GenerateForwardGraph {
			file, err := forwardgraph.Generate(inFile, messages)
			if err != nil {
				log.Fatalf("Could not template forwardgraph: %v", err)
			}

			resp.File = append(resp.File, file)
		}
		if parameters.GenerateForwardIds {
			file, err := forwardids.Generate(inFile, messages)
			if err != nil {
				log.Fatalf("Could not template forwardids: %v", err)
			}

			resp.File = append(resp.File, file)
		}
	}

	data, err = proto.Marshal(resp)
	if err != nil {
		log.Fatalf("Could not marshal response: %v [%v]", err, resp)
	}

	_, err = os.Stdout.Write(data)
	if err != nil {
		log.Fatalf("Could not write response to STDOUT: %v", err)
	}
}
