package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image/png"
	"net"
	"os"
	"time"

	"github.com/nfnt/resize"
)

func main() {
	fmt.Fprintln(os.Stderr, "loading image..")

	// Open the original image file
	file, err := os.Open("input.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Fprintln(os.Stderr, "decoding image..")
	// Decode the image
	img, err := png.Decode(file)
	if err != nil {
		panic(err)
	}

	// // Define the desired size
	newWidth := 256
	newHeight := 0

	x_offset := 0
	y_offset := 0

	fmt.Fprintln(os.Stderr, "scaling image..")
	// // Scale the image to the desired size
	scaledImg := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)

	packetByteBuffer := new(bytes.Buffer)
	binary.Write(packetByteBuffer, binary.BigEndian, uint8(128)) // type
	binary.Write(packetByteBuffer, binary.BigEndian, uint8(0))   // code
	binary.Write(packetByteBuffer, binary.BigEndian, uint8(0))   // checksum
	binary.Write(packetByteBuffer, binary.BigEndian, uint8(0))   // identifier
	binary.Write(packetByteBuffer, binary.BigEndian, uint8(0))   // sequencenumber
	packetBytes := packetByteBuffer.Bytes()
	_ = packetBytes

	local, _ := net.ResolveIPAddr("ip6", "2a0b:f4c0:c8:6d:39c:10d3:7afc:c352")
	remote, _ := net.ResolveIPAddr("ip6", "2400:8902:e001:233:32::")
	_ = remote
	conn, err := net.ListenIP("ip6:ipv6-icmp", local)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stderr, "preperation finished, sending pings..")
	bounds := scaledImg.Bounds()
	for {
		start := time.Now()
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				color := scaledImg.At(x, y)
				r, g, b, a := color.RGBA()
				addr := fmt.Sprintf("2400:8902:e001:233:%02x%02x:%02x:%02x:%02x", x_offset+x, y_offset+y, a*r/65536, a*g/65536, a*b/65536)
				dst, _ := net.ResolveIPAddr("ip6", addr)

				for {
					_, err := conn.WriteToIP(packetBytes, dst)
					if err == nil {
						break
					}
				}
			}
		}
		elapsed := time.Since(start)
		fmt.Fprintln(os.Stderr, "sent whole image in %s", elapsed)
		time.Sleep(1 * time.Second)
	}
	conn.Close()
}
