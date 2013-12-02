// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// All modifications from the original Copyright 2013 Joseph Bironas

package main

import (
	"github.com/dlintw/goconf"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"xbmcjson"
)


type clientRequest struct {
	Version string                 `json:"jsonrpc"`
	Method string                  `json:"method"`
	Params map[string] interface{} `json:"params"`
	Id     uint64                  `json:"id"`
}

func Connect(host, port string) (*rpc.Client, error) {
	client, err := xbmcjson.Dial("tcp", host + ":" + port)
	if err != nil {
		fmt.Printf("ERROR(dial):%v", err)
	} else {
		fmt.Printf("connected to %v:%v\n", host, port)
	}
	return client, err
}

func Request(c *rpc.Client, req *clientRequest, res string) (string, error){
	err := c.Call(req.Method, req.Params, &res)
	if err != nil {
		fmt.Printf("ERROR(call): %v\n", err)
	}
	return res, err
}

func Usage() {
	fmt.Printf("usage: %v [command]\n\n", os.Args[0])
	fmt.Printf("Supported commands are:\n")
	fmt.Printf("  clean                    -- cleans the library of non-existing things.\n")
	fmt.Printf("  notify <title> <message> -- Send notification.\n")
	fmt.Printf("  ping                     -- sends a jsonrpc ping to the specificed host\n")
	fmt.Printf("  reboot                   -- Reboots the system\n")
	fmt.Printf("  scantv                   -- Scans for new TV episodes\n")
	fmt.Printf("  scanmovies               -- Scans for new movies.\n")
	fmt.Printf("  scan <path>              -- Scans a specific path.\n")
	fmt.Printf("  sendtext <string>        -- Sends quoted text as input\n")
	fmt.Printf("  suspend                  -- Put system into suspend mode\n")
	fmt.Printf("  wake                     -- Sends Wake-on-LAN magic packet\n")
	fmt.Printf("  Use quotes around strings with spaces and such.\n\n")
	fmt.Printf("Config File:\n")
	fmt.Printf("  xbmc_config file format (put this in $HOME/.xbmc_config):\n\n")
	fmt.Printf("  [default]\n")
	fmt.Printf("  host = your.ip.addr\n")
	fmt.Printf("  port = 9090\n")
	fmt.Printf("  tv_path = smb://path/to/tv/\n")
	fmt.Printf("  movie_path = /path/to/movies\n")
	fmt.Printf("  music_path = /path/to/music\n")
	os.Exit(0)
}

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

func main() {
	homedir := os.ExpandEnv("$HOME")

	// Pull host/port from config file
	f, err := goconf.ReadConfigFile(homedir + "/.xbmc_config")
	if err != nil {
		fmt.Printf("ERROR:%v\n",err)
		os.Exit(1)
	}
	host, _ := f.GetString("default", "host")
        port, _ := f.GetString("default", "port")
	tv_path, _ := f.GetString("", "tv_path")
	movie_path, _ := f.GetString("", "movie_path")
	music_path, _ := f.GetString("", "music_path")
	mac_addr, _ := f.GetString("", "mac_addr")

	if len(os.Args) < 2 {
		Usage()
	}

	cmd := os.Args[1]
	args := os.Args[2:]

	req := &clientRequest{}
	req.Version = "2.0"
	req.Params = make(map[string] interface{})

	fmt.Printf("Sending %v... ", cmd)
	switch cmd {
	case "ping":
		c, err := Connect(host, port)
		if err != nil {
			fmt.Printf("ERROR:%v\n",err)
		}
		defer c.Close()
		req.Method = "JSONRPC.Ping"
		// Returns string
		var res string
		response, _ := Request(c, req, res)
		fmt.Printf("%v\n", response)
	case "reboot":
		c, err := Connect(host, port)
		if err != nil {
			fmt.Printf("ERROR:%v\n",err)
		}
		defer c.Close()
		req.Method = "System.Reboot"
		var res string
		response, _ := Request(c, req, res)
		fmt.Printf("%v\n", response)
	case "scan":
		c, err := Connect(host, port)
		if err != nil {
			fmt.Printf("ERROR:%v\n",err)
		}
		defer c.Close()
		req.Method = "VideoLibrary.Scan"
		var res string
		response, _ := Request(c, req, res)
		fmt.Printf("%v\n", response)
	case "scanmusic":
		c, err := Connect(host, port)
		if err != nil {
			fmt.Printf("ERROR:%v\n",err)
		}
		defer c.Close()
		req.Method = "AudioLibrary.Scan"
		req.Params["directory"] = music_path
		var res string
		response, _ := Request(c, req, res)
		fmt.Printf("%v\n", response)
	case "scantv":
		c, err := Connect(host, port)
		if err != nil {
			fmt.Printf("ERROR:%v\n",err)
		}
		defer c.Close()
		req.Method = "VideoLibrary.Scan"
		//req.Params["directory"] = string
		req.Params["directory"] = tv_path
		// Returns string
		var res string
		response, _ := Request(c, req, res)
		fmt.Printf("%v\n", response)
	case "scanmovies":
		c, err := Connect(host, port)
		if err != nil {
			fmt.Printf("ERROR:%v\n",err)
		}
		defer c.Close()
		req.Method = "VideoLibrary.Scan"
		req.Params["directory"] = movie_path
		// Still returns string
		var res string
		response, _ := Request(c, req, res)
		fmt.Printf("%v\n", response)
	case "clean":
		c, err := Connect(host, port)
		if err != nil {
			fmt.Printf("ERROR:%v\n",err)
		}
		defer c.Close()
		req.Method = "VideoLibrary.Clean"
		// Returns string
		var res string
		response, _ := Request(c, req, res)
		fmt.Printf("%v\n", response)
	case "cleanmusic":
		c, err := Connect(host, port)
		if err != nil {
			fmt.Printf("ERROR:%v\n",err)
		}
		defer c.Close()
		req.Method = "AudioLibrary.Clean"
		var res string
		response, _ := Request(c, req, res)
		fmt.Printf("%v\n", response)
	case "sendtext":
		c, err := Connect(host, port)
		if err != nil {
			fmt.Printf("ERROR:%v\n",err)
		}
		defer c.Close()
		// takes next command line argument as input
		// use quotes for strings on the command line
		req.Method = "Input.SendText"
		req.Params["text"] = args[0]
		var res string
 		// To be honest, I don't know what this returns
		response, _ := Request(c, req, res)
		fmt.Printf("%v\n", response)
	case "notify":
		c, err := Connect(host, port)
		if err != nil {
			fmt.Printf("ERROR:%v\n",err)
		}
		defer c.Close()
		// takes next two command line arguments as title
		// and message
		req.Method = "GUI.ShowNotification"
 		req.Params["title"] = args[0]
		req.Params["message"] = args[1]
		var res string
		// Returns something other than a string
		response, _ := Request(c, req, res)
		fmt.Printf("%v\n", response)
	case "setvolume":
		c, err := Connect(host, port)
		if err != nil {
			fmt.Printf("ERROR:%v\n",err)
		}
		defer c.Close()
		req.Method = "Application.SetVolume"
		req.Params["volume"], _ = strconv.Atoi(args[0])
		//res := &clientResponse{}
		var res string
		response, _ := Request(c, req, res)
		fmt.Printf("%v\n", response)
	case "pause":
		c, err := Connect(host, port)
		if err != nil {
			fmt.Printf("ERROR:%v\n",err)
		}
		defer c.Close()
		req.Method = "Player.PlayPause"
		req.Params[""] = ""
	case "wake":
		bcastAddr := "192.168.2.255"
		bcastPort := "9"
		err := SendMagicPacket(mac_addr, bcastAddr, bcastPort)
		if err != nil {
			fmt.Printf("ERROR:%v\n", err)
		}
	case "suspend":
		c, err := Connect(host, port)
		if err != nil {
			fmt.Printf("ERROR:%v\n",err)
		}
		defer c.Close()
	    req.Method = "System.Suspend"
		var res string
		response, _ := Request(c, req, res)
		fmt.Printf("%v\n", response)
	default:
		Usage()
	}
}
