# ipfs.io/refs

Tool for building and publishing the DAG at https://ipfs.io/refs

[![](https://img.shields.io/badge/project-IPFS-blue.svg?style=flat-square)](http://ipfs.io/) [![](https://img.shields.io/badge/freenode-%23ipfs-blue.svg?style=flat-square)](http://webchat.freenode.net/?channels=%23ipfs)

Renders a DAG from various data sources:

- DMCA notices served for the gateway at ipfs.io -- https://github.com/ipfs/refs-denylists-dmca
- TODO: Content archived on our storage hosts -- https://github.com/ipfs/refs-solarnet-storage

# Usage

```sh
# render, add, and publish the denylist
$ make build && make publish
$ git add versions/ && git commit -m Publish && git push

# update the TXT record for refs.ipfs.io
$ npm install && make dnslink
# wait for it to propagate
$ watch dig TXT refs.ipfs.io
```
