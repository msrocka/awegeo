package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		printHelp()
		return
	}

	reader, err := NewReader(os.Args[1])
	if err != nil {
		println("ERROR: failed to read file", os.Args[1])
		return
	}
	defer reader.Close()

	jsonFile, err := os.Create(os.Args[2])
	if err != nil {
		println("ERROR: failed to write file", jsonFile)
		return
	}
	defer jsonFile.Close()

	var features []*Feature
	for {
		placemark, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			println("ERROR: failed to parse KML:", err.Error())
			break
		}

		feature := placemark.asFeature()
		if feature != nil {
			features = append(features, feature)
			if len(features)%1000 == 0 {
				println("parsed", len(features), "features")
			}
		}
	}
	println("finished", len(features), "features")

	coll := FeatureCollection{
		Type:     "FeatureCollection",
		Features: features,
	}
	bytes, err := json.Marshal(coll)
	if err != nil {
		println("ERROR: failed to create JSON output:", err)
		return
	}

	out := bufio.NewWriter(jsonFile)
	_, err = out.Write(bytes)
	if err != nil {
		println("ERROR: failed to write JSON file:", err)
	}
	out.Flush()
	println("all done")
}

func printHelp() {
	fmt.Println(`
usage:
  $> awegeo [input kml file] [output json file]`)
	fmt.Println()
}
