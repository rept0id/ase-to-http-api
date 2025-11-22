package main

import (
	"log"
	"net"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	ginEngine := gin.Default()

	ginEngine.GET("/ase", func(c *gin.Context) {
		var err error

		var udpConn *net.UDPConn

		var buffer []byte
		var bufferLen int
		var bufferPos int

		var res struct {
			TypeRes

			Data     TypeAseRes     `json:"Data"`
			Metadata map[string]any `json:"Metadata"`
		}

		/*** * * ***/

		buffer = make([]byte, 4096)

		res.Metadata = make(map[string]any)

		/*** * * ***/

		udpAddr, err := net.ResolveUDPAddr("udp", "94.23.158.180:22126")
		if err != nil {
			panic(err)
		}

		udpConn, err = net.DialUDP("udp", nil, udpAddr)
		if err != nil {
			panic(err)
		}
		defer udpConn.Close()

		udpConn.SetDeadline(time.Now().Add(2 * time.Second))

		_, err = udpConn.Write([]byte("s"))
		if err != nil {
			panic(err)
		}

		bufferLen, _, err = udpConn.ReadFromUDP(buffer)
		if err != nil {
			return
		}

		// header
		res.Data.Header = string(buffer[:4])
		bufferPos += 4
		// info
		for bufferOff := bufferPos; bufferOff < bufferLen; bufferOff++ {
			bufOffNext := bufferOff + 1

			if bufOffNext > bufferLen {
				break
			}

			if buffer[bufOffNext] == 0x3f {
				res.Data.Info = string(buffer[bufferPos:bufferOff])

				bufferPos = bufferOff

				break
			}
		}
		// players
		// check : https://github.com/multitheftauto/mtasa-blue/blob/62d4a53bd32b4acae27837c7592931758307b762/Server/mods/deathmatch/logic/ASE.cpp#L236
		for bufferOff := bufferPos; bufferOff < bufferLen; bufferOff++ {
			bufferOffNext := bufferOff + 1

			if bufferOffNext > bufferLen || buffer[bufferOffNext] == 0x3f {
				var player TypeAseResPlayer

				buf := buffer[bufferPos:bufferOff]
				bufPos := 0
				bufLen := len(buf)

				// delimiter ("?")
				//
				bufPos += 2
				if bufPos > bufLen {
					continue
				}

				// name
				//
				nameLen := int(buf[bufPos]) - 1
				bufPos++ // name length consumed
				if bufPos > bufLen {
					continue
				}
				//
				//
				if bufPos+nameLen > bufLen {
					continue
				}
				player.Name = string(buf[bufPos : bufPos+nameLen])
				//
				//
				bufPos += nameLen // name consumed

				// team (skipped)
				bufPos++ // skip `(unsigned char)1;`

				// skin (skipped)
				bufPos++ // skip `(unsigned char)1;`

				// score
				//
				scoreLen := int(buf[bufPos]) - 1
				bufPos++ // score length consumed
				if bufPos > bufLen {
					continue
				}
				//
				//
				if bufPos+scoreLen > bufLen {
					continue
				}
				player.Score, err = strconv.Atoi(string(buf[bufPos : bufPos+scoreLen]))
				if err != nil {
					log.Println(err)
				}
				//
				//
				bufPos += scoreLen // score consumed

				// ping
				//
				pingLen := int(buf[bufPos]) - 1
				bufPos++ // ping length consumed
				if bufPos > bufLen {
					continue
				}
				//
				//
				if bufPos+pingLen > bufLen {
					continue
				}
				player.Ping, err = strconv.Atoi(string(buf[bufPos : bufPos+pingLen]))
				if err != nil {
					log.Print(err)
				}
				//
				//
				bufPos += scoreLen // score consumed

				// time (skipped)
				bufPos++ // skip `(unsigned char)1;`

				/*** * * ***/

				res.Data.Players = append(res.Data.Players, player)

				bufferPos = bufferOff

				continue
			}
		}

		/*** * * ***/

		c.JSON(200, res)
	})

	ginEngine.Run(":80")
}
