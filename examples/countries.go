package main

import "github.com/mmcloughlin/globe"

func main() {
	g := globe.New()
	g.DrawGraticule(15.0)
	g.DrawCountryBoundaries()
	g.SavePNG("countries.png", 400)
}
