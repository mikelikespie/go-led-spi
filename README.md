# OPC LED SPI Drivers in Go

Go implementation to drive APA102 LED strips.

Server-side counterpart to [swiftled](https://github.com/mikelikespie/swiftled). Probably compatible with other things that talk [OPC](http://openpixelcontrol.org/).

MIT license.

## Installation

1. Install go
2. `go get github.com/mikelikespie/go-led-spi/spi`

## ./spi

Implementation to drive SPI drivers

## ./ledserver

[OPC](http://openpixelcontrol.org/) server that uses the spi library.
