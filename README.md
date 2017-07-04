# globe

Globe visualizations in Golang backed by
[pinhole](https://github.com/tidwall/pinhole).





<p align="center"><img src="http://i.imgur.com/VuDcbKB.png" /></p>

## Getting Started

Install `globe` with

```sh
$ go get -u github.com/mmcloughlin/globe
```

Here's how to plot all the Starbucks locations.

```go
shops, err := LoadCoffeeShops("./starbucks.json")
if err != nil {
	return nil, err
}

green := color.NRGBA{0x00, 0x64, 0x3c, 192}
g := globe.New()
g.DrawGraticule(10.0)
for _, s := range shops {
	g.DrawDot(s.Lat, s.Lng, 0.05, globe.Color(green))
}
g.CenterOn(40.645423, -73.903879)
```

<p align="center"><img src="http://i.imgur.com/do3m4Bj.png" /></p>

## License

`globe` is available under the ISC [License](/LICENSE).