geodata: land.geodata.go

world.topojson:
	wget -O $@ https://unpkg.com/world-atlas@1.1.4/world/110m.json

%.world.geojson: world.topojson
	topo2geo --in $< $*=$@

%.geodata.go: %.world.geojson includegeojson.go
	go run includegeojson.go -input $< -output $@ -var $*
	gofmt -s -w $@
