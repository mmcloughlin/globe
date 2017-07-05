package main

import "github.com/mmcloughlin/globe"

func main() {
	g := globe.New()
	g.DrawGraticule(10.0)
	g.SavePNG("graticule.png", 400)
}
