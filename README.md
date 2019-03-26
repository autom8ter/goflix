# GoFlix

`go get github.com/autom8ter/goflix`

## Usage

```text
   _____     _______     
  / ___/__  / __/ (_)_ __
 / (_ / _ \/ _// / /\ \ /
 \___/\___/_/ /_/_//_\_\
 
 Usage:
   goflix [command]
 
 Available Commands:
   help        Help about any command
   serve       start goflix server
 
 Flags:
   -h, --help   help for goflix
 
 Use "goflix [command] --help" for more information about a command.

```

`goflix serve --help`

```text
start goflix server

Usage:
  goflix serve [flags]

Flags:
  -h, --help               help for serve
      --max-conn int       max connections (default 200)
  -p, --player string      Open the stream with a video player (VLC, MPV, MPlayer) (default "vlc")
      --port int           port to serve torrent (default 8080)
      --seed               to seed or not to seed
      --tcp                to tcp or not to tcp (default true)
  -t, --torrent string     torrent to download
      --torrent-port int   port to torrent from (default 50007)

```