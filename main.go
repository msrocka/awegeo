package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
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
	}
	defer kmlFile.Close()

	buff := bufio.NewReader(kmlFile)
	decoder := xml.NewDecoder(buff)
	placemarks := 0
	for {
		token, _ := decoder.Token()
		if token == nil {
			break
		}
		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "Placemark" {
				placemarks++
			}
		}

		if placemarks > 5 {
			break
		}
	}

	fmt.Println("found", placemarks, "placemarks")
}

func printHelp() {
	fmt.Println(`
usage:
  $> awegeo [input kml file] [output json file]`)
	fmt.Println()
}