geodata: land.geodata.go countries.geodata.go

world.topojson:
	wget -nv -O $@ https://unpkg.com/world-atlas@1.1.4/world/110m.json

%.world.geojson: world.topojson
	topo2geo --in $< $*=$@

%.geodata.go: %.world.geojson buildgeodata.go
	go run buildgeodata.go -input $< -output $@ -var $*
	gofmt -s -w $@

%.md: %.md.j2
	j2 $< > $@

tools:
	pip3 install j2cli==v0.3.2.post0
	npm install -g topojson

testimages:
	go test -images

testimagehashes: testimages
	md5sum Test*.png
