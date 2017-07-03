package globe

import (
	"image/color"
	"math"

	"github.com/tidwall/pinhole"
)

const precision = 1.0

type Style struct {
	GraticuleColor color.Color
	DotColor       color.Color
	Background     color.Color
	LineWidth      float64
	Scale          float64
}

func (s Style) imageOptions() *pinhole.ImageOptions {
	return &pinhole.ImageOptions{
		BGColor:   s.Background,
		LineWidth: s.LineWidth,
		Scale:     s.Scale,
	}
}

var DefaultStyle = Style{
	GraticuleColor: color.Gray{192},
	DotColor:       color.NRGBA{255, 0, 0, 255},
	Background:     color.White,
	LineWidth:      0.1,
	Scale:          0.7,
}

type Globe struct {
	p     *pinhole.Pinhole
	style Style
}

func New() *Globe {
	return &Globe{
		p:     pinhole.New(),
		style: DefaultStyle,
	}
}

type Option func(*Globe)

func Color(c color.Color) Option {
	return func(g *Globe) {
		g.p.Colorize(c)
	}
}

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

func (g *Globe) DrawParallel(lat float64, style ...Option) {
	defer g.styled(Color(g.style.GraticuleColor), style...)()
	for lng := -180.0; lng < 180.0; lng += precision {
		x1, y1, z1 := cartestian(lat, lng)
		x2, y2, z2 := cartestian(lat, lng+precision)
		g.p.DrawLine(x1, y1, z1, x2, y2, z2)
	}
}

func (g *Globe) DrawParallels(interval float64, style ...Option) {
	g.DrawParallel(0, style...)
	for lat := interval; lat < 90.0; lat += interval {
		g.DrawParallel(lat, style...)
		g.DrawParallel(-lat, style...)
	}
}

func (g *Globe) DrawMeridian(lng float64, style ...Option) {
	defer g.styled(Color(g.style.GraticuleColor), style...)()
	for lat := -90.0; lat < 90.0; lat += precision {
		x1, y1, z1 := cartestian(lat, lng)
		x2, y2, z2 := cartestian(lat+precision, lng)
		g.p.DrawLine(x1, y1, z1, x2, y2, z2)
	}
}

func (g *Globe) DrawMeridians(interval float64, style ...Option) {
	for lng := -180.0; lng < 180.0; lng += interval {
		g.DrawMeridian(lng, style...)
	}
}

func (g *Globe) DrawGraticules(interval float64, style ...Option) {
	g.DrawParallels(interval, style...)
	g.DrawMeridians(interval, style...)
}

func (g *Globe) DrawDot(lat, lng float64, radius float64, style ...Option) {
	defer g.styled(Color(g.style.DotColor), style...)()
	x, y, z := cartestian(lat, lng)
	g.p.DrawDot(x, y, z, radius)
}

func (g *Globe) CenterOn(lat, lng float64) {
	g.p.Rotate(0, 0, -degToRad(lng)-math.Pi/2)
	g.p.Rotate(degToRad(lat)+math.Pi/2, 0, 0)
}

func (g *Globe) SavePNG(filename string, side int) error {
	opts := g.style.imageOptions()
	return g.p.SavePNG(filename, side, side, opts)
}

func cartestian(lat, lng float64) (x, y, z float64) {
	phi := degToRad(lat)
	lambda := degToRad(lng)
	x = math.Cos(phi) * math.Cos(lambda)
	y = math.Cos(phi) * math.Sin(lambda)
	z = math.Sin(phi)
	return
}

func degToRad(d float64) float64 {
	return math.Pi * d / 180.0
}
