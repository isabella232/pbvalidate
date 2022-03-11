[![No Maintenance Intended](http://unmaintained.tech/badge.svg)](http://unmaintained.tech/)

# DEPRECATED

⛔️ This repository is now considered **deprecated** as is not in-use and it provides a very basic functionality that can be achieved using in-house mechanisms.

For the sake of granting a slow decomission, this repository will be archived **at the end of April**: 29th, April.

# pbvalidate

[![Build Status](https://travis-ci.com/bitnami-labs/pbvalidate.svg?branch=master)](https://travis-ci.com/bitnami-labs/pbvalidate)
[![Go Report Card](https://goreportcard.com/badge/github.com/bitnami-labs/pbvalidate)](https://goreportcard.com/report/github.com/bitnami-labs/pbvalidate)

pbvalidate validates pbjson files against a protobuf message


## Usage

```
$ pbvalidate -f some.proto -I somedir,someotherdir -m pkg1.Msg1 somefile.json
```

## Contributing

PRs accepted.
