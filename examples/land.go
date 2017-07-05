package main

import "github.com/mmcloughlin/globe"

func main() {
	g := globe.New()
	g.DrawGraticule(10.0)
	g.DrawLandBoundaries()
	g.CenterOn(51.453349, -2.588323)
	g.SavePNG("land.png", 400)
}
