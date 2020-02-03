package main

import (
	"bufio"
	"encoding/json"
	"os"
	"regexp"
	"strconv"
	"strings"
)

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

// A Writer converts KML placemarks to GeoJson features and
// writes them to a file.
type Writer struct {
	regex    *regexp.Regexp
	file     *os.File
	features []*Feature
}

// NewWriter creates a new writer to convert and write
// GeoJson features to the given file.
func NewWriter(jsonFile string) (*Writer, error) {
	file, err := os.Create(jsonFile)
	if err != nil {
		return nil, err
	}
	regex := regexp.MustCompile(
		`<tr><td>([^<]+)<\/td><td>([^<]+)<\/td><\/tr>`)
	w := &Writer{
		regex: regex,
		file:  file,
	}
	return w, nil
}

// Close closes the writer. It finally writes the converted
// to the file of the writer and closes the corresponding
// resources.
func (w *Writer) Close() {
	defer w.file.Close()
	coll := FeatureCollection{
		Type:     "FeatureCollection",
		Features: w.features,
	}
	bytes, err := json.Marshal(coll)
	if err != nil {
		println("ERROR: failed to create JSON output:", err)
		return
	}
	out := bufio.NewWriter(w.file)
	_, err = out.Write(bytes)
	if err != nil {
		println("ERROR: failed to write JSON file:", err)
	}
	out.Flush()
}

// Put converts the given KML placemark into a GeoJson feature
// and adds it to the writer.
func (w *Writer) Put(p *Placemark) {
	if p == nil {
		return
	}
	f := Feature{
		Type: "Feature",
		Geometry: &Geometry{
			Type: "Polygon",
		},
	}
	f.Geometry.Coordinates = make([][][]float64, 1)
	f.Geometry.Coordinates[0] = parseCoordinates(p.Coordinates)
	f.Properties = w.parseProps(p.Description)
	f.Properties["id"] = p.ID
	w.features = append(w.features, &f)
}

func parseCoordinates(coordinates string) [][]float64 {
	points := make([][]float64, 0)
	for _, point := range strings.Split(coordinates, " ") {
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
	return points
}

func (w *Writer) parseProps(description string) map[string]interface{} {
	props := make(map[string]interface{})
	matches := w.regex.FindAllStringSubmatch(description, -1)
	for _, match := range matches {
		key := match[1]
		val := match[2]
		num, err := strconv.ParseFloat(val, 64)
		if err == nil {
			if num != -9999 {
				props[key] = num
			}
		} else {
			props[key] = val
		}
	}
	return props
}
