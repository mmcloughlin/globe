all: starbucks.png cities.png

%.png: %.go
	go run -tags $* main.go $< -filename $@ -side 600

starbucks.png: starbucks.json

starbucks.json:
	wget -O $@ 'https://github.com/mmcloughlin/starbucks/blob/77b395b047373fdb430d971f47560735d554e575/locations.json?raw=true'

cities.png: cities.json

cities.json:
	wget -O cities.json 'https://unpkg.com/cities.json@1.0.2/cities.json'

clean:
	$(RM) *.json *.png
