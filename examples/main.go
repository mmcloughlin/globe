package main

import (
	"flag"
	"log"
)

var (
	filename = flag.String("filename", "example.png", "image filename")
	side     = flag.Int("side", 500, "output image size")
)

func main() {
	flag.Parse()

	g, err := Build()
	if err != nil {
		log.Fatal(err)
	}

	err = g.SavePNG(*filename, *side)
	if err != nil {
		log.Fatal(err)
	}
}
