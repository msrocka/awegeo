package main

import (
	"archive/zip"
	"bufio"
	"encoding/xml"
	"errors"
	"io"
	"os"
	"strings"
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
	file    io.ReadCloser
	zipFile *zip.ReadCloser
	decoder *xml.Decoder
}

// Close the underlying KML file of the reader.
func (r *Reader) Close() {
	if r.zipFile != nil {
		r.zipFile.Close()
	}
	if r.file != nil {
		r.file.Close()
	}
}

// NewReader creates a new reader for the given KML file.
func NewReader(kmlFile string) (*Reader, error) {

	reader := &Reader{}

	if !strings.HasSuffix(kmlFile, ".kmz") {

		// try to read the file as KML file
		kml, err := os.Open(kmlFile)
		if err != nil {
			return nil, err
		}
		reader.file = kml

	} else {

		// try to read the KML document from a KMZ file
		zipFile, err := zip.OpenReader(kmlFile)
		if err != nil {
			return nil, err
		}
		reader.zipFile = zipFile

		// search for the *kml entry in the zip
		var zipEntry *zip.File
		for i := range zipFile.File {
			f := zipFile.File[i]
			if strings.HasSuffix(f.Name, ".kml") {
				zipEntry = f
				break
			}
		}
		if zipEntry == nil {
			return nil, errors.New(
				"Could not find KML file in " + kmlFile)
		}

		kmz, err := zipEntry.Open()
		if err != nil {
			return nil, err
		}
		reader.file = kmz
	}

	buff := bufio.NewReader(reader.file)
	reader.decoder = xml.NewDecoder(buff)
	return reader, nil
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
