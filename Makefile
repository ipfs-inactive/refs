local="http://localhost:8080/ipfs/"
gway="https://ipfs.io/ipfs/"
domain="refs.ipfs.io"
record="dmca.denylists.lists"

publish:
	go run denylist.go | tail -n1 >versions/current
	cat versions/current >>versions/history
	@export hash=`cat versions/current`; \
		echo ""; \
		echo "new version:"; \
		echo "- $(local)$$hash"; \
		echo "- $(gway)$$hash"; \
		echo ""; \
		echo "next:"; \
		echo "- pin it: /ipfs/$$hash"; \
		echo "- update dnslink: make dnslink";

# Only run after publish, or there won't be a path to set.
dnslink: auth.token node_modules
	DIGITAL_OCEAN=$(shell cat auth.token) node_modules/.bin/dnslink-deploy \
		--domain=$(domain) --record=$(record) --path=/ipfs/$(shell cat versions/current)

node_modules: package.json
	npm install
	touch node_modules

.PHONY: publish dnslink
