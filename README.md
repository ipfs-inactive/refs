# lists-denylists-dmca

DMCA notices, and tools for publishing them
at https://ipfs.io/refs/lists/denylists/dmca

[![](https://img.shields.io/badge/project-IPFS-blue.svg?style=flat-square)](http://ipfs.io/) [![](https://img.shields.io/badge/freenode-%23ipfs-blue.svg?style=flat-square)](http://webchat.freenode.net/?channels=%23ipfs)

- All DMCA takedown notices for ipfs.io and gateway.ipfs.io,
  each with a list of the affected hashes.
- A golang tool for rendering them to IPFS objects.
- TODO: make tasks for adding, pinning, and publishing them via IPNS.

TODO: The code which links https://ipfs.io/refs/a/b/c to /ipns/c.b.a.refs.ipfs.io
is part of the gateways' nginx configuration:
[roles/ipfs_gateway/templates/nginx_ipfs_gateway.conf.j2](https://github.com/ipfs/infrastructure/blob/master/solarnet/roles/ipfs_gateway/templates/nginx_ipfs_gateway.conf.j2)

# Usage

```sh
$ cd lists-denylists-dmca/
$ make publish
https://ipfs.io/refs/lists/denylists/dmca
https://ipfs.io/ipfs/QmRER7erZxU63huYgSBryGhKrfHdkDkVjwQTd8RD4RdSW5
QmRER7erZxU63huYgSBryGhKrfHdkDkVjwQTd8RD4RdSW5
```

Long version:

```sh
$ cd lists-denylists-dmca/
# render the denylist, and add it
$ hash=$(go run dmca.go)
# pin the denylist
$ ipfs pin add /ipfs/$hash
# update the TXT record for IPNS
$ dnslink-deploy --domain refs.ipfs.io --record dmca.denylists.lists --path /ipfs/$hash
# wait for it to propagate
$ watch dig TXT dmca.denylists.lists.refs.ipfs.io
```
