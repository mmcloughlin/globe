all: starbucks.png

%.png: %.go
	go run main.go $< -filename $@ -side 500

starbucks.png: starbucks.json

starbucks.json:
	wget -O $@ 'https://github.com/mmcloughlin/starbucks/blob/77b395b047373fdb430d971f47560735d554e575/locations.json?raw=true'

clean:
	$(RM) *.json *.png
