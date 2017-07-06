package main

import (
	"image/color"

	"github.com/mmcloughlin/globe"
)

func main() {
	g := globe.New()
	g.DrawGraticule(10.0)
	g.DrawLandBoundaries()
	g.DrawLine(
		51.453349, -2.588323,
		40.645423, -73.903879,
		globe.Color(color.NRGBA{255, 0, 0, 255}),
	)
	g.CenterOn(50.244440, -37.207949)
	g.SavePNG("line.png", 400)
}
