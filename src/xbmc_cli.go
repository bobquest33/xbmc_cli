// Copyright 2013 Joseph Bironas.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/dlintw/goconf"
	"fmt"
	"net/rpc"
	"os"
	"strconv"
	"xbmcjson"
	"xbmcwol"
)


type clientRequest struct {
	Version string                 `json:"jsonrpc"`
	Method string                  `json:"method"`
	Params map[string] interface{} `json:"params"`
	Id     uint64                  `json:"id"`
}

func 

func Connect(host, port string) (*rpc.Client) {
	client, err := xbmcjson.Dial("tcp", host + ":" + port)
	if err != nil {
		fmt.Printf("ERROR(dial):%v", err)
	} else {
		fmt.Printf("connected to %v:%v\n", host, port)
	}
	return client
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
	fmt.Printf("  # these are required for WOL support\n")
	fmt.Printf("  mac_addr = 00:01:23:45:67:89\n")
	fmt.Printf("  bcast_addr = your.bcast.addr\n")
	fmt.Printf("  bcasr_port = your.bcast.port\n")
	os.Exit(0)
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
	bcast_addr, _ := f.GetString("", "bcast_addr")
	bcast_port, _ := f.GetString("", "bcast_port")

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
		c := Connect(host, port)
		defer c.Close()
		req.Method = "JSONRPC.Ping"
		// Returns string
		var res string
		response, _ := Request(c, req, res)
		fmt.Printf("%v\n", response)
	case "reboot":
		c := Connect(host, port)
		defer c.Close()
		req.Method = "System.Reboot"
		var res string
		response, _ := Request(c, req, res)
		fmt.Printf("%v\n", response)
	case "scan":
		c := Connect(host, port)
		defer c.Close()
		req.Method = "VideoLibrary.Scan"
		var res string
		response, _ := Request(c, req, res)
		fmt.Printf("%v\n", response)
	case "scanmusic":
		c := Connect(host, port)
		defer c.Close()
		req.Method = "AudioLibrary.Scan"
		req.Params["directory"] = music_path
		var res string
		response, _ := Request(c, req, res)
		fmt.Printf("%v\n", response)
	case "scantv":
		c := Connect(host, port)
		defer c.Close()
		req.Method = "VideoLibrary.Scan"
		//req.Params["directory"] = string
		req.Params["directory"] = tv_path
		// Returns string
		var res string
		response, _ := Request(c, req, res)
		fmt.Printf("%v\n", response)
	case "scanmovies":
		c := Connect(host, port)
		defer c.Close()
		req.Method = "VideoLibrary.Scan"
		req.Params["directory"] = movie_path
		// Still returns string
		var res string
		response, _ := Request(c, req, res)
		fmt.Printf("%v\n", response)
	case "clean":
		c := Connect(host, port)
		defer c.Close()
		req.Method = "VideoLibrary.Clean"
		// Returns string
		var res string
		response, _ := Request(c, req, res)
		fmt.Printf("%v\n", response)
	case "cleanmusic":
		c := Connect(host, port)
		defer c.Close()
		req.Method = "AudioLibrary.Clean"
		var res string
		response, _ := Request(c, req, res)
		fmt.Printf("%v\n", response)
	case "sendtext":
		c := Connect(host, port)
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
		c := Connect(host, port)
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
		c := Connect(host, port)
		defer c.Close()
		req.Method = "Application.SetVolume"
		req.Params["volume"], _ = strconv.Atoi(args[0])
		//res := &clientResponse{}
		var res string
		response, _ := Request(c, req, res)
		fmt.Printf("%v\n", response)
	case "pause":
		c := Connect(host, port)
		defer c.Close()
		req.Method = "Player.PlayPause"
		req.Params[""] = ""
	case "wake":
		err := xbmcwol.SendMagicPacket(mac_addr, bcast_addr, bcast_port)
		if err != nil {
			fmt.Printf("ERROR:%v\n", err)
		}
	case "suspend":
		c := Connect(host, port)
		defer c.Close()
	    req.Method = "System.Suspend"
		var res string
		response, _ := Request(c, req, res)
		fmt.Printf("%v\n", response)
	default:
		Usage()
	}
}
