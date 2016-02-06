// +build linux
package spi

// #cgo CFLAGS: -std=c99
// #include <sys/ioctl.h>
// #include <linux/spi/spidev.h>
import "C"

import (
	"syscall"
	"unsafe"
)

const (
	SPI_IOC_WR_MODE          = uintptr(C.SPI_IOC_WR_MODE)
	SPI_IOC_RD_MODE          = uintptr(C.SPI_IOC_RD_MODE)
	SPI_IOC_WR_BITS_PER_WORD = uintptr(C.SPI_IOC_WR_BITS_PER_WORD)
	SPI_IOC_RD_BITS_PER_WORD = uintptr(C.SPI_IOC_RD_BITS_PER_WORD)
	SPI_IOC_WR_MAX_SPEED_HZ  = uintptr(C.SPI_IOC_WR_MAX_SPEED_HZ)
	SPI_IOC_RD_MAX_SPEED_HZ  = uintptr(C.SPI_IOC_RD_MAX_SPEED_HZ)
)

// Functions for doing SPI syscalls. Return value isn't checked. Don't use these on real systems

func spiSetWrMode(fd uintptr, spiMode uint8) {
	syscall.Syscall(syscall.SYS_IOCTL, fd, SPI_IOC_WR_MODE, uintptr(unsafe.Pointer(&spiMode)))
}

func spiSetRdMode(fd uintptr, spiMode uint8) {
	syscall.Syscall(syscall.SYS_IOCTL, fd, SPI_IOC_RD_MODE, uintptr(unsafe.Pointer(&spiMode)))
}

func spiSetWrBitsPerWord(fd uintptr, bitsPerWord uint8) {
	syscall.Syscall(syscall.SYS_IOCTL, fd, SPI_IOC_WR_BITS_PER_WORD, uintptr(unsafe.Pointer(&bitsPerWord)))
}

func spiSetRdBitsPerWord(fd uintptr, bitsPerWord uint8) {
	syscall.Syscall(syscall.SYS_IOCTL, fd, SPI_IOC_RD_BITS_PER_WORD, uintptr(unsafe.Pointer(&bitsPerWord)))
}

func spiSetWrMaxSpeedHz(fd uintptr, maxSpeed uint32) {
	syscall.Syscall(syscall.SYS_IOCTL, fd, SPI_IOC_WR_MAX_SPEED_HZ, uintptr(unsafe.Pointer(&maxSpeed)))
}

func spiSetRdMaxSpeedHz(fd uintptr, maxSpeed uint32) {
	syscall.Syscall(syscall.SYS_IOCTL, fd, SPI_IOC_RD_MAX_SPEED_HZ, uintptr(unsafe.Pointer(&maxSpeed)))
}
