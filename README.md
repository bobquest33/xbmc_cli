xbmc_cli
========

XBMC cli using the JSON interface


Usage: xbmc_cli [command]

Supported commands are:
  clean                    -- cleans the library of non-existing things.
  notify <title> <message> -- Send notification.
  ping                     -- sends a jsonrpc ping to the specificed host
  reboot                   -- Reboots the system
  scantv                   -- Scans for new TV episodes (broken currently)
  scanmovies               -- Scans for new movies. (broken currently)
  scan <path>              -- Scans a specific path. (doesn't take path)
  sendtext <string>        -- Sends quoted text as input
  suspend                  -- Put system into suspend mode
  wake                     -- Sends Wake-on-LAN magic packet
  Use quotes around strings with spaces and such.
Config File:
  xbmc_config file format (put this in $HOME/.xbmc_config):
  [default]
  host = your.ip.addr
  port = 9090
  tv_path = smb://path/to/tv/
  movie_path = /path/to/movies
  music_path = /path/to/music
  # for WOL function
  mac_addr = 00:01:23:45:67:89  
  bcast_addr = your.network.bcast.addr (ex. 192.168.1.255)
  bcast_port = your.xbmc.wol.port # mine seems to prefer 9, ymmv
