// +build ignore

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	geojson "github.com/paulmach/go.geojson"
)

var (
	inputFilepath  string
	outputFilepath string
	variableName   string
)

func init() {
	flag.StringVar(&inputFilepath, "input", "", "Input GeoJSON file")
	flag.StringVar(&outputFilepath, "output", "", "Output Go file")
	flag.StringVar(&variableName, "var", "", "Variable name")
}

func LoadFeatureCollection(filename string) (*geojson.FeatureCollection, error) {
	b, err := ioutil.ReadFile(inputFilepath)
	if err != nil {
		return nil, err
	}

	return geojson.UnmarshalFeatureCollection(b)
}

func ExtractPathsFromFeatureCollection(collection *geojson.FeatureCollection) [][][]float64 {
	paths := [][][]float64{}
	for _, feature := range collection.Features {
		geom := feature.Geometry
		var layer4 [][][][]float64
		switch geom.Type {
		case geojson.GeometryPolygon:
			layer4 = [][][][]float64{geom.Polygon}
		case geojson.GeometryMultiPolygon:
			layer4 = geom.MultiPolygon
		case geojson.GeometryPoint, geojson.GeometryMultiPoint:
			log.Printf("discarding point geometry type %s", string(geom.Type))
		default:
			log.Fatalf("no handler for geometry type %s", string(geom.Type))
		}

		for _, layer3 := range layer4 {
			paths = append(paths, layer3...)
		}
	}

	return paths
}

func WritePathsCode(w io.Writer, varname string, paths [][][]float64) error {
	fmt.Fprint(w, "// Generated code. DO NOT EDIT.\n")
	fmt.Fprintf(w, "// Arguments: %s\n\n", strings.Join(os.Args[1:], " "))
	fmt.Fprint(w, "package globe\n")

	fmt.Fprintf(w, "var %s = [][]struct{\nlat, lng float32\n}{\n", varname)
	for _, path := range paths {
		fmt.Fprint(w, "{\n")
		for _, point := range path {
			if len(point) != 2 {
				return errors.New("point must have two coordinates")
			}
			fmt.Fprintf(w, "{%v,%v},\n", point[1], point[0])
		}
		fmt.Fprint(w, "},\n")
	}
	fmt.Fprint(w, "}\n")

	return nil
}

func main() {
	flag.Parse()

	collection, err := LoadFeatureCollection(inputFilepath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("loaded %d features", len(collection.Features))

	paths := ExtractPathsFromFeatureCollection(collection)
	log.Printf("extracted %d paths", len(paths))

	f, err := os.Create(outputFilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = WritePathsCode(f, variableName, paths)
	if err != nil {
		log.Fatal(err)
	}
}
