package spi

import (
	"fmt"
	"os"
)

func SPIOpen(device string, speedHz int) (*os.File, error) {
	file, err := os.OpenFile(device, os.O_WRONLY, 0x600)

	fd := file.Fd()

	if err != nil {
		return nil, fmt.Errorf("Unable to open SPI device %s withe error %+v", device, err)
	}

	spiSetWrMode(fd, 0)
	spiSetRdMode(fd, 0)

	spiSetWrBitsPerWord(fd, 8)
	spiSetRdBitsPerWord(fd, 8)

	spiSetWrMaxSpeedHz(fd, uint32(speedHz))
	spiSetRdMaxSpeedHz(fd, uint32(speedHz))

	return file, nil
}
