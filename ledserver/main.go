package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"time"

	"github.com/mikelikespie/go-led-spi/spi"
)

type Color struct {
	A       uint8
	R, G, B uint8
}

var spiDevicePath = flag.String("d", "/dev/spidev0.0", "spidev device to connect to")
var spiFrequency = flag.Int("S", 12000000, "SPI Frequency")

var ledCount = flag.Int("c", 540, "LED count")
var listenPort = flag.Int("p", 7890, "TCP listen port")

var reusableBuffers = make(chan []Color, 10)

func getBuffer() []Color {
	select {
	case existingBuffer := <-reusableBuffers:
		return existingBuffer
	default:
		log.Println("making buffer")
		return make([]Color, *ledCount)
	}
}

func returnBuffer(buffer []Color) {
	select {
	case reusableBuffers <- buffer:
	default:
	}
}

func main() {
	var file *os.File
	var err error

	flag.Parse()

	file, err = spi.SPIOpen(*spiDevicePath, *spiFrequency)

	if err != nil {
		log.Panicln("Error:", err)
	}

	ledCount := *ledCount

	buffer := make([]byte, 4+ledCount*4+ledCount/16+2)

	// Set the end frame
	for i := 4 + ledCount*4; i < len(buffer); i++ {
		buffer[i] = 0xFF
	}

	const pendingBufferCount = 2

	/// It will buffer one frame of color
	renderFrameChan := make(chan []Color, pendingBufferCount)

	// go runDemo(renderFrameChan, ledCount)
	go listen(*listenPort, renderFrameChan, ledCount)

	lastReportTime := time.Now()
	framesSinceReport := 0

	for colors := range renderFrameChan {
		for i, color := range colors {
			offset := i*4 + 4
			var brightness uint8 = 31

			if color.A < 32 {
				brightness = color.A
			}

			buffer[offset+0] = 0xE0 | brightness
			buffer[offset+1] = color.B
			buffer[offset+2] = color.G
			buffer[offset+3] = color.R
		}
		bytesWritten, err := file.Write(buffer)

		if err != nil {
			log.Panicln("Error", err)
		}

		if bytesWritten != len(buffer) {
			log.Panicln("Didn't write all the bytes")
		}

		framesSinceReport += 1

		timeSinceLastReport := time.Since(lastReportTime)
		if timeSinceLastReport > 1.0*time.Second && framesSinceReport > 10 {
			log.Println("FPS:", float64(framesSinceReport)/float64(timeSinceLastReport)*float64(time.Second))
			framesSinceReport = 0
			lastReportTime = time.Now()
		}
	}
}

func handleConnection(conn net.Conn, renderFrameChan chan<- []Color, ledCount int) {
	headerBuffer := make([]byte, 4)
	var readBuffer []byte
	defer conn.Close()
	for true {
		count, err := conn.Read(headerBuffer)

		if err != nil {
			log.Println("Reading header failed. continuing.", err)
			return
		}

		if count != 4 {
			log.Println("Didn't read full header. Continuing")
			return
		}

		var size = int(headerBuffer[2])<<8 | int(headerBuffer[3])

		switch {
		case readBuffer == nil || len(readBuffer) < size:
			readBuffer = make([]byte, size)
		case len(readBuffer) > size:
			readBuffer = readBuffer[0:size]
		}

		var bytesRead = 0
		for bytesRead < size {
			count, err = conn.Read(readBuffer[bytesRead:len(readBuffer)])
			if err != nil {
				log.Println("Reading colors failed. continuing.", err)
				return
			}

			bytesRead += count
		}

		var colors = getBuffer()
		for i, _ := range colors {
			offset := i * 3

			if offset+3 > len(readBuffer) {
				break
			}

			colors[i] = Color{A: 2, R: readBuffer[offset+0], G: readBuffer[offset+1], B: readBuffer[offset+2]}
		}

		renderFrameChan <- colors

		returnBuffer(colors)
	}
}

func listen(port int, renderFrameChan chan<- []Color, ledCount int) {
	ln, err := net.Listen("tcp", fmt.Sprint(":", port))
	if err != nil {
		log.Panicln("Listen failed.aborting", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Accept failed. continuing.", err)
			continue
		}
		go handleConnection(conn, renderFrameChan, ledCount)
	}
}

func runDemo(renderFrameChan chan<- []Color, ledCount int) {
	start := time.Now()
	for true {
		/// TODO: reuse these buffers
		var colors = getBuffer()
		for i, _ := range colors {
			duration := time.Since(start)
			v := math.Sin(float64(duration)/float64(time.Second)*10 + math.Pi*2*float64(i%18)/18)
			unitV := (v*.5 + .5)
			unitV *= unitV
			adjustedV := uint8(unitV * 255)
			a := uint8(31)
			switch {
			case unitV == 0:
				break
			case adjustedV == 1:
				adjustedV = 2
				// a has to be somewhere betwen 16 and 32
				a = uint8((unitV-1.0/255)*255*16) + 16
			case adjustedV < 1 && unitV > 0:
				adjustedV = 1
				a = uint8(unitV * 255 * 32)
			}
			colors[i] = Color{A: a, R: adjustedV, G: adjustedV, B: adjustedV}
		}
		renderFrameChan <- colors
		returnBuffer(colors)
	}

}
