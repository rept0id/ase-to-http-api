package main

import (
	"errors"
	"log"
	"math"
	"net"
	"strconv"
	"time"
)

func ase(ip string, port int) (aseRes TypeAseRes, err error) {
	var udpConn *net.UDPConn

	var buffer []byte
	var bufferLen int
	var bufferPos int

	/*** * * ***/

	buffer = make([]byte, 4096)

	/*** * * ***/

	portStr := strconv.Itoa(port)

	udpAddr, err := net.ResolveUDPAddr("udp", ip+":"+portStr)
	if err != nil {
		return aseRes, err
	}

	udpConn, err = net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return aseRes, err
	}
	defer udpConn.Close()

	udpConn.SetDeadline(time.Now().Add(2 * time.Second))

	_, err = udpConn.Write([]byte("s"))
	if err != nil {
		return aseRes, err
	}

	bufferLen, _, err = udpConn.ReadFromUDP(buffer)
	if err != nil {
		return aseRes, err
	}

	/*** * * ***/

	// header
	aseRes.Header = string(buffer[:4])
	bufferPos += 4 // header consumed
	if bufferPos >= bufferLen {
		err = errors.New("bufferPos >= bufferLen")

		return aseRes, err
	}

	// info
	for bufferOff := bufferPos; bufferOff < bufferLen; bufferOff++ {
		bufOffNext := bufferOff + 1

		if bufOffNext >= bufferLen {
			err = errors.New("bufferPos >= bufferLen")

			return aseRes, err
			// break
		}

		if buffer[bufOffNext] == 0x3f {
			aseRes.Info = string(buffer[bufferPos:bufferOff])

			bufferPos = bufferOff

			break
		}
	}

	// players
	// check : https://github.com/multitheftauto/mtasa-blue/blob/62d4a53bd32b4acae27837c7592931758307b762/Server/mods/deathmatch/logic/ASE.cpp#L236
	for bufferOff := bufferPos; bufferOff < bufferLen; bufferOff++ {
		bufferOffNext := bufferOff + 1

		if bufferOffNext >= bufferLen || buffer[bufferOffNext] == 0x3f {
			var e error // small error

			var b []byte // small buffer
			var bPos int // small buffer pos
			var bLen int // small buffer length

			var player TypeAseResPlayer

			/*** * * ***/

			b = buffer[bufferPos:bufferOff]
			bPos = 0
			bLen = len(b)

			// delimiter ("?")
			//
			bPos += 2
			if bPos >= bLen {
				log.Print(errors.New("bufPos >= bufLen"))

				continue
			}

			// name
			//
			nameLen := int(b[bPos]) - 1
			bPos++ // name length consumed
			if bPos >= bLen {
				log.Print(errors.New("bufPos >= bufLen"))

				continue
			}
			//
			//
			if bPos+nameLen >= bLen {
				log.Print(errors.New("bufPos >= bufLen"))

				continue
			}
			player.Name = string(b[bPos : bPos+nameLen])
			//
			//
			bPos += nameLen // name consumed

			// team (skipped)
			bPos++ // skip `(unsigned char)1;`

			// skin (skipped)
			bPos++ // skip `(unsigned char)1;`

			// score
			//
			scoreLen := int(b[bPos]) - 1
			bPos++ // score length consumed
			if bPos >= bLen {
				log.Print(errors.New("bufPos >= bufLen"))

				continue
			}
			//
			//
			if bPos+scoreLen >= bLen {
				log.Print(errors.New("bufPos+scoreLen >= bufLen"))

				continue
			}
			player.Score, e = strconv.Atoi(string(b[bPos : bPos+scoreLen]))
			if e != nil {
				player.Score = math.MinInt
				log.Print(e)
			}
			//
			//
			bPos += scoreLen // score consumed

			// ping
			//
			pingLen := int(b[bPos]) - 1
			bPos++ // ping length consumed
			if bPos >= bLen {
				log.Print(errors.New("bufPos >= bufLen"))

				continue
			}
			//
			//
			if bPos+pingLen >= bLen {
				log.Print(errors.New("bufPos >= bufLen"))

				continue
			}
			player.Ping, e = strconv.Atoi(string(b[bPos : bPos+pingLen]))
			if e != nil {
				player.Ping = math.MinInt
				log.Print(e)
			}
			//
			//
			bPos += scoreLen // score consumed

			// time (skipped)
			bPos++ // skip `(unsigned char)1;`

			/*** * * ***/

			aseRes.Players = append(aseRes.Players, player)

			/*** * * ***/

			bufferPos = bufferOff
		}
	}

	/*** * * ***/

	return aseRes, err
}
