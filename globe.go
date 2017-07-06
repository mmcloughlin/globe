// Package globe builds 3D visualizations on the earth.
package globe

import (
	"image"
	"image/color"
	"math"

	"github.com/tidwall/pinhole"
)

// Precision constants.
const (
	// graticuleLineStep is the gap between nodes of a parallel or meridian line
	// in degrees.
	graticuleLineStep = 1.0

	// linePointInterval is the max distance (in km) between line segments when
	// drawing along a great circle.
	linePointInterval = 500.0
)

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
	for lng := -180.0; lng < 180.0; lng += graticuleLineStep {
		x1, y1, z1 := cartestian(lat, lng)
		x2, y2, z2 := cartestian(lat, lng+graticuleLineStep)
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
	for lat := -90.0; lat < 90.0; lat += graticuleLineStep {
		x1, y1, z1 := cartestian(lat, lng)
		x2, y2, z2 := cartestian(lat+graticuleLineStep, lng)
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

// DrawLine draws a line between (lat1, lng1) and (lat2, lng2) along the great
// circle.
// Uses the default LineColor unless overridden by style Options.
func (g *Globe) DrawLine(lat1, lng1, lat2, lng2 float64, style ...Option) {
	defer g.styled(Color(g.style.LineColor), style...)()

	d := haversine(lat1, lng1, lat2, lng2)
	step := d / math.Ceil(d/linePointInterval)
	fx, fy, fz := cartestian(lat1, lng1)
	for p := step; p < d-step/2; p += step {
		tlat, tlng := intermediate(lat1, lng1, lat2, lng2, p/d)
		tx, ty, tz := cartestian(tlat, tlng)
		g.p.DrawLine(fx, fy, fz, tx, ty, tz)
		fx, fy, fz = tx, ty, tz
	}

	tx, ty, tz := cartestian(lat2, lng2)
	g.p.DrawLine(fx, fy, fz, tx, ty, tz)
}

// DrawRect draws the rectangle with the given corners. Sides are drawn along
// great circles, as in DrawLine.
// Uses the default LineColor unless overridden by style Options.
func (g *Globe) DrawRect(minlat, minlng, maxlat, maxlng float64, style ...Option) {
	g.DrawLine(minlat, minlng, maxlat, minlng, style...)
	g.DrawLine(maxlat, minlng, maxlat, maxlng, style...)
	g.DrawLine(maxlat, maxlng, minlat, maxlng, style...)
	g.DrawLine(minlat, maxlng, minlat, minlng, style...)
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

// Image renders an image object for the visualization with dimensions
// (side, side).
func (g *Globe) Image(side int) *image.RGBA {
	opts := g.style.imageOptions()
	return g.p.Image(side, side, opts)
}

// SavePNG writes the visualization to filename in PNG format with dimensions
// (side, side).
func (g *Globe) SavePNG(filename string, side int) error {
	opts := g.style.imageOptions()
	return g.p.SavePNG(filename, side, side, opts)
}

// cartestian maps (lat, lng) to pinhole cartestian space.
func cartestian(lat, lng float64) (x, y, z float64) {
	x = cos(lat) * cos(lng)
	y = cos(lat) * sin(lng)
	z = -sin(lat)
	return
}

// earthRadius is the radius of the earth.
const earthRadius = 6371.0

// haversine returns the distance (in km) between the points (lat1, lng1) and
// (lat2, lng2).
func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	dlat := lat2 - lat1
	dlng := lng2 - lng1
	a := sin(dlat/2)*sin(dlat/2) + cos(lat1)*cos(lat2)*sin(dlng/2)*sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadius * c
}

// intermediate returns the point that is fraction f between (lat1, lng1) and
// (lat2, lng2).
func intermediate(lat1, lng1, lat2, lng2, f float64) (float64, float64) {
	dr := haversine(lat1, lng1, lat2, lng2) / earthRadius
	a := math.Sin((1-f)*dr) / math.Sin(dr)
	b := math.Sin(f*dr) / math.Sin(dr)
	x := a*cos(lat1)*cos(lng1) + b*cos(lat2)*cos(lng2)
	y := a*cos(lat1)*sin(lng1) + b*cos(lat2)*sin(lng2)
	z := a*sin(lat1) + b*sin(lat2)
	phi := math.Atan2(z, math.Sqrt(x*x+y*y))
	lambda := math.Atan2(y, x)
	return radToDeg(phi), radToDeg(lambda)
}

// destination computes the destination point reached when travelling distance d
// from (lat, lng) at bearing brng.
func destination(lat, lng, d, brng float64) (float64, float64) {
	dr := d / earthRadius
	phi := math.Asin(sin(lat)*math.Cos(dr) + cos(lat)*math.Sin(dr)*cos(brng))
	lambda := degToRad(lng) + math.Atan2(sin(brng)*math.Sin(dr)*cos(lat), math.Cos(dr)-sin(lat)*math.Sin(phi))
	return radToDeg(phi), math.Mod(radToDeg(lambda)+540, 360) - 180
}

// sin is math.Sin for degrees.
func sin(d float64) float64 { return math.Sin(degToRad(d)) }

// cos is math.Cos for degrees
func cos(d float64) float64 { return math.Cos(degToRad(d)) }

// degToRad converts d degrees to radians.
func degToRad(d float64) float64 {
	return math.Pi * d / 180.0
}

// radToDeg converts r radians to degrees.
func radToDeg(r float64) float64 {
	return 180.0 * r / math.Pi
}
