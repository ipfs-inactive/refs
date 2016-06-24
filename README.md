# refs

> DMCA notices, and tools for publishing them

[![](https://img.shields.io/badge/made%20by-Protocol%20Labs-blue.svg?style=flat-square)](http://ipn.io)
[![](https://img.shields.io/badge/project-IPFS-blue.svg?style=flat-square)](http://ipfs.io/)
[![](https://img.shields.io/badge/freenode-%23ipfs-blue.svg?style=flat-square)](http://webchat.freenode.net/?channels=%23ipfs)
[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)

Tool for building and publishing the DAG at https://ipfs.io/refs

Renders a DAG from various data sources:

- DMCA notices served for the gateway at ipfs.io -- https://github.com/ipfs/refs-denylists-dmca
- TODO: Content archived on our storage hosts -- https://github.com/ipfs/refs-solarnet-storage

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Contribute](#contribute)
- [License](#license)

## Install

Clone this repo. Depends on [npm](https://npmjs.com) and [node.js](https://nodejs.com).

## Usage

```sh
# render, add, and publish the denylist
$ make build && make publish
$ git add versions/ && git commit -m Publish && git push

# update the TXT record for refs.ipfs.io
$ npm install && make dnslink
# wait for it to propagate
$ watch dig TXT refs.ipfs.io
```

## Contribute

Feel free to join in. All welcome. Open an [issue](https://github.com/ipfs/refs/issues)!

This repository falls under the IPFS [Code of Conduct](https://github.com/ipfs/community/blob/master/code-of-conduct.md).

[![](https://cdn.rawgit.com/jbenet/contribute-ipfs-gif/master/img/contribute.gif)](https://github.com/ipfs/community/blob/master/contributing.md)

## License

[MIT](LICENSE)
