examples=graticule land countries line rect starbucks cities

images: $(addsuffix .png, $(examples))

upload: $(addsuffix .imgur, $(examples))

snippets: $(addsuffix .snip, $(examples))

%.png: %.go
	go run $<

%.imgur: %.png tools
	imgur -a=true $< > $@

%.snip: %.go
	awk '/func main/ {flag=1; next} /^}/ {flag=0} flag {print substr($$0,2)}' < $< > $@

starbucks.png: starbucks.json

starbucks.json:
	wget -nv -O $@ 'https://github.com/mmcloughlin/starbucks/blob/77b395b047373fdb430d971f47560735d554e575/locations.json?raw=true'

cities.png: cities.json

cities.json:
	wget -nv -O $@ 'https://unpkg.com/cities.json@1.1.2/cities.json'

tools:
	go install github.com/mattn/imgur@v0.0.0-20201112131930-001c86d92ed2

clean:
	$(RM) *.json *.png
