package main

import (
	"image/color"

	"github.com/mmcloughlin/globe"
)

func main() {
	g := globe.New()
	g.DrawGraticule(10.0)
	g.DrawLandBoundaries()
	g.DrawRect(
		41.897209, 12.500285,
		55.782693, 37.615993,
		globe.Color(color.NRGBA{255, 0, 0, 255}),
	)
	g.CenterOn(48, 25)
	g.SavePNG("rect.png", 400)
}
