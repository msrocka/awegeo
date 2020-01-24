package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		printHelp()
		return
	}

	kmlFile, err := os.Open(os.Args[1])
	if err != nil {
		println("ERROR: failed to read file", kmlFile)
		return
	}
	defer kmlFile.Close()

	jsonFile, err := os.Create(os.Args[2])
	if err != nil {
		println("ERROR: failed to write file", jsonFile)
		return
	}
	defer jsonFile.Close()

	buff := bufio.NewReader(kmlFile)
	decoder := xml.NewDecoder(buff)

	var features []*Feature
	for {
		token, err := decoder.Token()
		if err != nil && err != io.EOF {
			println("ERROR: failed to parse KML:", err.Error())
			break
		}
		if token == nil {
			break
		}

		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "Placemark" {
				var p Placemark
				err = decoder.DecodeElement(&p, &t)
				if err != nil {
					break
				}
				f := p.asFeature()
				if f != nil {
					features = append(features, f)
					if len(features)%1000 == 0 {
						println("parsed", len(features), "features")
					}
				}

			}
		}

		if err != nil {
			println("ERROR: failed to parse placemark", err)
			break
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
