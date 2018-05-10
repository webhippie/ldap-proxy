# LDAP Proxy

[![Build Status](http://github.dronehippie.de/api/badges/webhippie/ldap-proxy/status.svg)](http://github.dronehippie.de/webhippie/ldap-proxy)
[![Stories in Ready](https://badge.waffle.io/webhippie/ldap-proxy.svg?label=ready&title=Ready)](http://waffle.io/webhippie/ldap-proxy)
[![Join the Matrix chat at https://matrix.to/#/#webhippie:matrix.org](https://img.shields.io/badge/matrix-%23webhippie%3Amatrix.org-7bc9a4.svg)](https://matrix.to/#/#webhippie:matrix.org)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/4d3a80070282417fab062d1df363f7e2)](https://www.codacy.com/app/webhippie/ldap-proxy?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=webhippie/ldap-proxy&amp;utm_campaign=Badge_Grade)
[![Go Doc](https://godoc.org/github.com/webhippie/ldap-proxy?status.svg)](http://godoc.org/github.com/webhippie/ldap-proxy)
[![Go Report](https://goreportcard.com/badge/github.com/webhippie/ldap-proxy)](https://goreportcard.com/report/github.com/webhippie/ldap-proxy)
[![](https://images.microbadger.com/badges/image/tboerger/ldap-proxy.svg)](http://microbadger.com/images/tboerger/ldap-proxy "Get your own image badge on microbadger.com")
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/1832/badge)](https://bestpractices.coreinfrastructure.org/projects/1832)

**This project is under heavy development, it's not in a working state yet!**

A reverse proxy and static file server that provides an authentication layer via LDAP to any web application that doesn't support it natively.


## Docs

Our documentation gets generated directly out of the [docs/](docs/) folder, it get's built via Drone and published to GitHub pages. You can find the documentation at [https://webhippie.github.io/ldap-proxy/](https://webhippie.github.io/ldap-proxy/).


## Install

You can download prebuilt binaries from the GitHub releases or from our [download site](http://dl.webhippie.de/misc/ldap-proxy). You are a Mac user? Just take a look at our [homebrew formula](https://github.com/webhippie/homebrew-webhippie).


## Development

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](http://golang.org/doc/install.html). This project requires Go >= v1.8.

```bash
go get -d github.com/webhippie/ldap-proxy
cd $GOPATH/src/github.com/webhippie/ldap-proxy
make retool sync clean generate build

./bin/ldap-proxy -h
```


## Security

If you find a security issue please contact thomas@webhippie.de first.


## Contributing

Fork -> Patch -> Push -> Pull Request


## Authors

* [Thomas Boerger](https://github.com/tboerger)


## License

Apache-2.0


## Copyright

```
Copyright (c) 2018 Thomas Boerger <http://www.webhippie.de>
```
