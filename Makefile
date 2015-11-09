local="http://localhost:8080/ipfs/"
gway="https://ipfs.io/ipfs/"
domain="refs.ipfs.io"
record="@"

build:
	rm -rf dmca/notices
	git clone https://github.com/ipfs/refs-denylists-dmca.git dmca/notices
	go-bindata -pkg dmca -o dmca/bindata.go -ignore '^dmca/notices/\.git' dmca/notices/...

publish:
	go run main.go -current=$(shell cat versions/current) | tail -n1 >versions/current
	cat versions/current >>versions/history
	@export hash=`cat versions/current`; \
		echo ""; \
		echo "new version:"; \
		echo "- $(local)$$hash"; \
		echo "- $(gway)$$hash"; \
		echo ""; \
		echo "next:"; \
		echo "- pin $$hash"; \
		echo "- make dnslink";

# Only run after publish, or there won't be a path to set.
dnslink: node_modules
	DIGITAL_OCEAN=$(shell cat $(HOME)/.protocol/digitalocean.key) node_modules/.bin/dnslink-deploy \
		--domain=$(domain) --record=$(record) --path=/ipfs/$(shell cat versions/current)

node_modules: package.json
	npm install
	touch node_modules

.PHONY: publish dnslink
