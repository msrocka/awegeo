package main

import (
	"strconv"
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

// A FeatureCollection of features created from KML placemarks. This
// is the root document when building the corresponding GeoJson
// document.
type FeatureCollection struct {
	Type     string     `json:"type"`
	Features []*Feature `json:"features"`
}

// A Feature in the GeoJson format that we build from a KML Placemark.
type Feature struct {
	Type       string                 `json:"type"`
	Geometry   *Geometry              `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

// A Geometry of a GeoJson feature. This is a simple polygon.
type Geometry struct {
	Type        string        `json:"type"`
	Coordinates [][][]float64 `json:"coordinates"`
}

func (p *Placemark) asFeature() *Feature {
	if p == nil {
		return nil
	}
	f := Feature{
		Type: "Feature",
		Geometry: &Geometry{
			Type: "Polygon",
		},
	}

	points := make([][]float64, 0)
	for _, point := range strings.Split(p.Coordinates, " ") {
		p := strings.Split(strings.TrimSpace(point), ",")
		if len(p) < 2 {
			continue
		}
		x, err := strconv.ParseFloat(p[0], 64)
		if err != nil {
			print("WARNING: x is not numeric in coordinate:", point)
			continue
		}
		y, err := strconv.ParseFloat(p[1], 64)
		if err != nil {
			print("WARNING: y is not numeric in coordinate:", point)
			continue
		}
		points = append(points, []float64{x, y})
	}
	f.Geometry.Coordinates = make([][][]float64, 1)
	f.Geometry.Coordinates[0] = points
	return &f
}
