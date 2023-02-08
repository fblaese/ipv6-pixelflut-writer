package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {

	fmt.Fprintln(os.Stderr, "socket setup..")
	packetByteBuffer := new(bytes.Buffer)
	binary.Write(packetByteBuffer, binary.BigEndian, uint8(128)) // type
	binary.Write(packetByteBuffer, binary.BigEndian, uint8(0))   // code
	binary.Write(packetByteBuffer, binary.BigEndian, uint8(0))   // checksum
	binary.Write(packetByteBuffer, binary.BigEndian, uint8(0))   // identifier
	binary.Write(packetByteBuffer, binary.BigEndian, uint8(0))   // sequencenumber
	packetBytes := packetByteBuffer.Bytes()
	_ = packetBytes

	local, _ := net.ResolveIPAddr("ip6", "2a02:810d:ab40:2500:b6af:a0fc:c84b:4ea2")
	remote, _ := net.ResolveIPAddr("ip6", "2400:8902:e001:233:32::")
	_ = remote
	conn, err := net.ListenIP("ip6:ipv6-icmp", local)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(os.Stderr, "preperation finished, sending pings..")
	for {
		for i := 0; i < 256; {
			start := time.Now()
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
			for y := 28; y < 64; y++ {
				for x := 128; x < 160; x++ {
					r := 0
					g := uint32(0)
					b := uint32(0)

					addr := fmt.Sprintf("2400:8902:e001:233:%02x%02x:%02x:%02x:%02x", x, y, r, g, b)
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

			time.Sleep(time.Second / 4)
		}
	}
	conn.Close()
}
