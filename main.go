package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		printHelp()
		return
	}

	// set up the reader and writer
	reader, err := NewReader(os.Args[1])
	if err != nil {
		println("ERROR: failed to read file", os.Args[1])
		return
	}
	defer reader.Close()
	writer, err := NewWriter(os.Args[2])
	if err != nil {
		println("ERROR: failed to write file", os.Args[2])
		return
	}
	defer writer.Close()

	channel := make(chan *Placemark, 1000)

	// read the placemarks from the reader
	// and put them into the channel
	go func() {
		count := 0
		for {
			placemark, err := reader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				println("ERROR: failed to parse KML:", err.Error())
				break
			}
			channel <- placemark
			count++
			if count%1000 == 0 {
				println("parsed", count, "features")
			}
		}
		close(channel)
		println("finished", count, "features")
	}()

	// read the placemarks from the channel and put
	// them into the writer
	for {
		placemark, more := <-channel
		if !more {
			break
		}
		writer.Put(placemark)
	}
	println("all done")
}

func printHelp() {
	fmt.Println(`
usage:
  $> awegeo [input kml file] [output json file]`)
	fmt.Println()
}
