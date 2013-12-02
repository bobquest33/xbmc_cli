/*
  Originally pulled this from github.com/ghthor/gowol 
  Extended interface to support bcastPort for xbmc.
*/

package xbmcwol

import (
	"encoding/hex"
	"errors"
	"log"
	"net"
	"strings"
)

func SendMagicPacket(macAddr string, bcastAddr string, bcastPort string) error {

	if len(macAddr) != (6*2 + 5) {
		return errors.New("Invalid MAC Address String: " + macAddr)
	}
	
	packet, err := constructMagicPacket(macAddr)
	if err != nil {
		return err
	}

	a, err := net.ResolveUDPAddr("udp", bcastAddr+":"+bcastPort)
	if err != nil {
		return err
	}

	c, err := net.DialUDP("udp", nil, a)
	if err != nil {
		return err
	}

	written, err := c.Write(packet)
	c.Close()

	// Packet must be 102 bytes in length
	if written != 102 {
		return err
	}

	return nil
}

func constructMagicPacket(macAddr string) ([]byte, error) {
	macBytes, err := hex.DecodeString(strings.Join(strings.Split(macAddr, ":"), ""))
	if err != nil {
		log.Fatalln("Error Hex Decoding:", err)
		return nil, err
	}

	b := []uint8{255, 255, 255, 255, 255, 255}
	for i := 0; i < 16; i++ {
		b = append(b, macBytes...)
	}
	return b, err
}
