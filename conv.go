package main

// The Placemark KML tag. We only fetch the data into fields that
// we use in the conversion. In the current version each placemark
// only contains a single polygon where the outer boundary is given.
type Placemark struct {
	ID          string `xml:"id,attr"`
	Description string `xml:"description"`
	Coordinates string `xml:"MultiGeometry>Polygon>outerBoundaryIs>LinearRing>coordinates"`
}

// A Feature in the GeoJson format that we build from a KML Placemark.
type Feature struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}

// A Geometry of a GeoJson feature. This is a simple polygon.
type Geometry struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}
