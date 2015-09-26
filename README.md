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
# render, add, and publish the denylist
$ make publish
$ git add versions/ && git commit -m 'Publish' && git push
# update the TXT record for dmca.denylists.lists.refs.ipfs.io
$ make dnslink
# wait for it to propagate
$ watch dig TXT dmca.denylists.lists.refs.ipfs.io
```
