package globe

import (
	"math"

	"github.com/tidwall/pinhole"
)

const precision = 1.0

type Globe struct {
	p *pinhole.Pinhole
}

func New() *Globe {
	return &Globe{
		p: pinhole.New(),
	}
}

func (g *Globe) DrawParallel(lat float64) {
	for lng := -180.0; lng < 180.0; lng += precision {
		x1, y1, z1 := cartestian(lat, lng)
		x2, y2, z2 := cartestian(lat, lng+precision)
		g.p.DrawLine(x1, y1, z1, x2, y2, z2)
	}
}

func (g *Globe) DrawParallels(interval float64) {
	g.DrawParallel(0)
	for lat := interval; lat < 90.0; lat += interval {
		g.DrawParallel(lat)
		g.DrawParallel(-lat)
	}
}

func (g *Globe) DrawMeridian(lng float64) {
	for lat := -90.0; lat < 90.0; lat += precision {
		x1, y1, z1 := cartestian(lat, lng)
		x2, y2, z2 := cartestian(lat+precision, lng)
		g.p.DrawLine(x1, y1, z1, x2, y2, z2)
	}
}

func (g *Globe) DrawMeridians(interval float64) {
	for lng := -180.0; lng < 180.0; lng += interval {
		g.DrawMeridian(lng)
	}
}

func (g *Globe) DrawGraticules(interval float64) {
	g.DrawParallels(interval)
	g.DrawMeridians(interval)
}

func (g *Globe) SavePNG(filename string, side int, opts *pinhole.ImageOptions) error {
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
