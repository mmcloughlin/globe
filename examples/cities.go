// +build cities

package main

import (
	"encoding/json"
	"image/color"
	"io/ioutil"
	"log"

	"github.com/mmcloughlin/globe"
)

type City struct {
	Lat float64
	Lng float64
}

func LoadCities(filename string) ([]City, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	raw := []struct {
		Lat json.Number `json:"lat"`
		Lng json.Number `json:"lng"`
	}{}
	err = json.Unmarshal(data, &raw)
	if err != nil {
		return nil, err
	}

	cities := make([]City, 0, len(raw))
	for _, r := range raw {
		lat, latErr := r.Lat.Float64()
		lng, lngErr := r.Lng.Float64()
		if latErr != nil || lngErr != nil {
			continue
		}
		cities = append(cities, City{
			Lat: lat,
			Lng: lng,
		})
	}

	return cities, nil
}

func main() {
	cities, err := LoadCities("./cities.json")
	if err != nil {
		log.Fatal(err)
	}

	c := color.NRGBA{0x11, 0x2F, 0x56, 128}
	g := globe.New()
	g.DrawGraticule(10.0)
	for _, s := range cities {
		g.DrawDot(s.Lat, s.Lng, 0.02, globe.Color(c))
	}
	g.CenterOn(51.453349, -2.588323)
	err = g.SavePNG("cities.png", 400)
	if err != nil {
		log.Fatal(err)
	}
}
