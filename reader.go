package main

import (
	"bufio"
	"encoding/xml"
	"os"
)

// The Placemark KML tag. We only fetch the data into fields that
// we use in the conversion. In the current version each placemark
// only contains a single polygon where the outer boundary is given.
type Placemark struct {
	ID          string `xml:"id,attr"`
	Description string `xml:"description"`
	Coordinates string `xml:"MultiGeometry>Polygon>outerBoundaryIs>LinearRing>coordinates"`
}

// Reader reads placemarks from a KML file.
type Reader struct {
	file    *os.File
	decoder *xml.Decoder
}

// Close the underlying KML file of the reader.
func (r *Reader) Close() {
	if r.file != nil {
		r.file.Close()
	}
}

// NewReader creates a new reader for the given KML file.
func NewReader(kmlFile string) (*Reader, error) {
	file, err := os.Open(kmlFile)
	if err != nil {
		return nil, err
	}
	buff := bufio.NewReader(file)
	decoder := xml.NewDecoder(buff)
	r := &Reader{
		file:    file,
		decoder: decoder,
	}
	return r, nil
}

// Next reads the next placemark from the reader. It returns
// (nil, io.EOL) if there are no more placemarks to read.
func (r *Reader) Next() (*Placemark, error) {
	for {
		token, err := r.decoder.Token()
		if err != nil {
			return nil, err
		}
		switch t := token.(type) {
		case xml.StartElement:
			if t.Name.Local == "Placemark" {
				var p Placemark
				err = r.decoder.DecodeElement(&p, &t)
				if err != nil {
					return nil, err
				}
				return &p, nil
			}
		}
	}
}
