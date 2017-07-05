// Package globe builds 3D visualizations on the earth.
package globe

import (
	"image/color"
	"math"

	"github.com/tidwall/pinhole"
)

// precision is the gap between nodes of a line in degrees.
const precision = 1.0

// Style encapsulates globe display options.
type Style struct {
	GraticuleColor color.Color
	LineColor      color.Color
	DotColor       color.Color
	Background     color.Color
	LineWidth      float64
	Scale          float64
}

// imageOptions builds the pinhole ImageOptions object for this Style.
func (s Style) imageOptions() *pinhole.ImageOptions {
	return &pinhole.ImageOptions{
		BGColor:   s.Background,
		LineWidth: s.LineWidth,
		Scale:     s.Scale,
	}
}

// DefaultStyle specifies out-of-the box style options.
var DefaultStyle = Style{
	GraticuleColor: color.Gray{192},
	LineColor:      color.Gray{32},
	DotColor:       color.NRGBA{255, 0, 0, 255},
	Background:     color.White,
	LineWidth:      0.1,
	Scale:          0.7,
}

// Globe is a globe visualization.
type Globe struct {
	p     *pinhole.Pinhole
	style Style
}

// New constructs an empty globe with the default style.
func New() *Globe {
	return &Globe{
		p:     pinhole.New(),
		style: DefaultStyle,
	}
}

// Option is a function that stylizes a globe.
type Option func(*Globe)

// Color uses the given color.
func Color(c color.Color) Option {
	return func(g *Globe) {
		g.p.Colorize(c)
	}
}

// styled is an internal convenience for applying style Options within a pinhole
// Begin/End context.
func (g *Globe) styled(base Option, options ...Option) func() {
	g.p.Begin()
	return func() {
		base(g)
		for _, option := range options {
			option(g)
		}
		g.p.End()
	}
}

// DrawParallel draws the parallel of latitude lat.
// Uses the default GraticuleColor unless overridden by style Options.
func (g *Globe) DrawParallel(lat float64, style ...Option) {
	defer g.styled(Color(g.style.GraticuleColor), style...)()
	for lng := -180.0; lng < 180.0; lng += precision {
		x1, y1, z1 := cartestian(lat, lng)
		x2, y2, z2 := cartestian(lat, lng+precision)
		g.p.DrawLine(x1, y1, z1, x2, y2, z2)
	}
}

// DrawParallels draws parallels at the given interval.
// Uses the default GraticuleColor unless overridden by style Options.
func (g *Globe) DrawParallels(interval float64, style ...Option) {
	g.DrawParallel(0, style...)
	for lat := interval; lat < 90.0; lat += interval {
		g.DrawParallel(lat, style...)
		g.DrawParallel(-lat, style...)
	}
}

// DrawMeridian draws the meridian at longitude lng.
// Uses the default GraticuleColor unless overridden by style Options.
func (g *Globe) DrawMeridian(lng float64, style ...Option) {
	defer g.styled(Color(g.style.GraticuleColor), style...)()
	for lat := -90.0; lat < 90.0; lat += precision {
		x1, y1, z1 := cartestian(lat, lng)
		x2, y2, z2 := cartestian(lat+precision, lng)
		g.p.DrawLine(x1, y1, z1, x2, y2, z2)
	}
}

// DrawMeridians draws meridians at the given interval.
// Uses the default GraticuleColor unless overridden by style Options.
func (g *Globe) DrawMeridians(interval float64, style ...Option) {
	for lng := -180.0; lng < 180.0; lng += interval {
		g.DrawMeridian(lng, style...)
	}
}

// DrawGraticule draws a latitude/longitude grid at the given interval.
// Uses the default GraticuleColor unless overridden by style Options.
func (g *Globe) DrawGraticule(interval float64, style ...Option) {
	g.DrawParallels(interval, style...)
	g.DrawMeridians(interval, style...)
}

// DrawDot draws a dot at (lat, lng) with the given radius.
// Uses the default DotColor unless overridden by style Options.
func (g *Globe) DrawDot(lat, lng float64, radius float64, style ...Option) {
	defer g.styled(Color(g.style.DotColor), style...)()
	x, y, z := cartestian(lat, lng)
	g.p.DrawDot(x, y, z, radius)
}

// DrawLandBoundaries draws land boundaries on the globe.
// Uses the default LineColor unless overridden by style Options.
func (g *Globe) DrawLandBoundaries(style ...Option) {
	g.drawPreparedPaths(land, style...)
}

// DrawCountryBoundaries draws country boundaries on the globe.
// Uses the default LineColor unless overridden by style Options.
func (g *Globe) DrawCountryBoundaries(style ...Option) {
	g.drawPreparedPaths(countries, style...)
}

func (g *Globe) drawPreparedPaths(paths [][]struct{ lat, lng float32 }, style ...Option) {
	defer g.styled(Color(g.style.LineColor), style...)()
	for _, path := range paths {
		n := len(path)
		for i := 0; i+1 < n; i++ {
			p1, p2 := path[i], path[i+1]
			x1, y1, z1 := cartestian(float64(p1.lat), float64(p1.lng))
			x2, y2, z2 := cartestian(float64(p2.lat), float64(p2.lng))
			g.p.DrawLine(x1, y1, z1, x2, y2, z2)
		}
	}
}

// CenterOn rotates the globe to center on (lat, lng).
func (g *Globe) CenterOn(lat, lng float64) {
	g.p.Rotate(0, 0, -degToRad(lng)-math.Pi/2)
	g.p.Rotate(math.Pi/2-degToRad(lat), 0, 0)
}

// SavePNG writes the visualization to filename in PNG format with dimensions
// (side, side).
func (g *Globe) SavePNG(filename string, side int) error {
	opts := g.style.imageOptions()
	return g.p.SavePNG(filename, side, side, opts)
}

// cartestian maps (lat, lng) to pinhole cartestian space.
func cartestian(lat, lng float64) (x, y, z float64) {
	phi := degToRad(lat)
	lambda := degToRad(lng)
	x = math.Cos(phi) * math.Cos(lambda)
	y = math.Cos(phi) * math.Sin(lambda)
	z = -math.Sin(phi)
	return
}

// degToRad converts d degrees to radians.
func degToRad(d float64) float64 {
	return math.Pi * d / 180.0
}
