package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/draw"
	"image/gif"
	"net"
	"os"
	"time"

	"github.com/nfnt/resize"
)

func main() {
	fmt.Fprintln(os.Stderr, "loading image..")

	// Open the original image file
	file, err := os.Open("input.gif")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Fprintln(os.Stderr, "decoding image..")
	// Decode the image
	gifimg, err := gif.DecodeAll(file)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stderr, "socket setup..")
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

	// // Define the desired size
	newWidth := 32
	newHeight := 0

	x_offset := 64
	y_offset := 16

	// var canvas [][]color.Color = make([][]color.Color, gifimg.Config.Width)
	// for i := range canvas {
	// 	canvas[i] = make([]color.Color, gifimg.Config.Height)
	// 	for j := range canvas[i] {
	// 		canvas[i][j] = color.Black
	// 	}
	// }

	fmt.Fprintln(os.Stderr, "preperation finished, sending pings..")
	for {
		canvas := image.NewRGBA(image.Rect(0, 0, gifimg.Config.Width, gifimg.Config.Height))
		draw.Draw(canvas, canvas.Bounds(), gifimg.Image[0], image.ZP, draw.Src)

		for i := range gifimg.Image {
			start := time.Now()

			img := gifimg.Image[i]
			delay := gifimg.Delay[i]

			draw.Draw(canvas, canvas.Bounds(), img, image.ZP, draw.Over)
			/*
				bounds := img.Bounds()
				for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
					for x := bounds.Min.X; x < bounds.Max.X; x++ {
						r1, g1, b1, _ := canvas.At(x, y).RGBA()

						colorr := img.At(x, y)
						r2, g2, b2, a2 := colorr.RGBA()
						a := float64(a2) / 65535
						r := uint16(float64(r2)*float64(a) + float64(r1)*(float64(1)-a))
						g := uint16(float64(g2)*float64(a) + float64(g1)*(float64(1)-a))
						b := uint16(float64(b2)*float64(a) + float64(b1)*(float64(1)-a))

						canvas.Set(x, y, color.RGBA64{
							R: r,
							G: g,
							B: b,
							A: 65535,
						})
					}
				}*/

			// f, _ := os.Create(fmt.Sprintf("frame%d.png", i))
			// png.Encode(f, canvas)
			// f.Close()

			// print canvas
			fmt.Fprintln(os.Stderr, "scaling image..")
			scaledImg := resize.Resize(uint(newWidth), uint(newHeight), canvas, resize.Lanczos3)
			bounds := scaledImg.Bounds()
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				for x := bounds.Min.X; x < bounds.Max.X; x++ {
					r, g, b, _ := scaledImg.At(x, y).RGBA()

					r = r / 256
					g = g / 256
					b = b / 256

					//fmt.Printf("r: %d, g: %d, b: %d\n", r, g, b)

					addr := fmt.Sprintf("2400:8902:e001:233:%02x%02x:%02x:%02x:%02x", x_offset+x, y_offset+y, r, g, b)
					dst, _ := net.ResolveIPAddr("ip6", addr)

					for {
						_, err := conn.WriteToIP(packetBytes, dst)
						if err == nil {
							break
						}
					}
				}
			}

			// wait
			elapsed := time.Since(start)
			fmt.Fprintln(os.Stderr, "sent whole image in %s", elapsed)

			time.Sleep(time.Second / 100 * time.Duration(delay))
			time.Sleep(time.Second / 4)
		}
	}
	conn.Close()
}
