# globe

Globe wireframe visualizations in Golang backed by
[pinhole](https://github.com/tidwall/pinhole).

[![GoDoc Reference](http://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](http://godoc.org/github.com/mmcloughlin/globe)
[![Build status](https://img.shields.io/travis/mmcloughlin/globe.svg?style=flat-square)](https://travis-ci.org/mmcloughlin/globe)





<p align="center"><img src="http://i.imgur.com/1xkWn2v.png" /></p>

## Getting Started

Install `globe` with

```sh
$ go get -u github.com/mmcloughlin/globe
```

Start with a blank globe with a graticule at 10 degree intervals.

```go
g := globe.New()
g.DrawGraticule(10.0)
g.SavePNG("graticule.png", 400)
```
<p align="center"><img src="http://i.imgur.com/XcfK1VI.png" /></p>

Add some land boundaries and center it on a point. Alternatively
[`DrawCountryBoundaries`](https://godoc.org/github.com/mmcloughlin/globe#Globe.DrawCountryBoundaries)
will give you countries.

```go
g := globe.New()
g.DrawGraticule(10.0)
g.DrawLandBoundaries()
g.CenterOn(51.453349, -2.588323)
g.SavePNG("land.png", 400)
```
<p align="center"><img src="http://i.imgur.com/VNsNySI.png" /></p>

Here's all the [Starbucks
locations](https://github.com/mmcloughlin/starbucks). Note `color.NRGBA`
recommended to [avoid
artifacts](https://github.com/mmcloughlin/globe/issues/6).

```go
shops, err := LoadCoffeeShops("./starbucks.json")
if err != nil {
	log.Fatal(err)
}

green := color.NRGBA{0x00, 0x64, 0x3c, 192}
g := globe.New()
g.DrawGraticule(10.0)
for _, s := range shops {
	g.DrawDot(s.Lat, s.Lng, 0.05, globe.Color(green))
}
g.CenterOn(40.645423, -73.903879)
err = g.SavePNG("starbucks.png", 400)
if err != nil {
	log.Fatal(err)
}
```
<p align="center"><img src="http://i.imgur.com/FyZABLE.png" /></p>

You can also do lines along great circles.

```go
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
```
<p align="center"><img src="http://i.imgur.com/6WLCWv8.png" /></p>

Also rectangles.

```go
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
```
<p align="center"><img src="http://i.imgur.com/KJfIl9l.png" /></p>

See [examples](examples/) and
[godoc](https://godoc.org/github.com/mmcloughlin/globe) for more.

## License

`globe` is available under the ISC [License](/LICENSE).