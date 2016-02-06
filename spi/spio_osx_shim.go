// +build !linux

// Shim for not having SPIO to ease development

package spi

/// Noops on OS X
func spiSetWrMode(fd uintptr, spiMode uint8) {
}

func spiSetRdMode(fd uintptr, spiMode uint8) {
}

func spiSetWrBitsPerWord(fd uintptr, bitsPerWord uint8) {
}

func spiSetRdBitsPerWord(fd uintptr, bitsPerWord uint8) {
}

func spiSetWrMaxSpeedHz(fd uintptr, maxSpeed uint32) {
}

func spiSetRdMaxSpeedHz(fd uintptr, maxSpeed uint32) {
}
