package main

import (
	"encoding/xml"
)

// The Placemark tag. We only fetch the data into fields that
// we use in the conversion. In the current version each
// placemark only contains a single polygon where the outer
// boundary is given.
type Placemark struct {
	XMLName     xml.Name `xml:"Placemark"`
	ID          string   `xml:"id,attr"`
	Description string   `xml:"description"`
	Coordinates string   `xml:"MultiGeometry>Polygon>outerBoundaryIs"`
}
