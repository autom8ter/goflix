# goflix
--
    import "github.com/autom8ter/goflix"

## Cli Usage

`go get github.com/autom8ter/goflix/cmd/goflix`

```text
  _____     _______     
 / ___/__  / __/ (_)_ __
/ (_ / _ \/ _// / /\ \ /
\___/\___/_/ /_/_//_\_\

Usage:
  goflix [command]

Available Commands:
  help        Help about any command
  play        Open stream in video player
  serve       start goflix server

Flags:
  -h, --help   help for goflix

Use "goflix [command] --help" for more information about a command.

```

## Usage

#### type ClientConfig

```go
type ClientConfig struct {
	TorrentPath    string
	Port           int
	TorrentPort    int
	Seed           bool
	TCP            bool
	MaxConnections int
}
```

ClientConfig specifies the behaviour of a GoFlix.

#### func  DefaultClientConfig

```go
func DefaultClientConfig() ClientConfig
```

#### type ClientError

```go
type ClientError struct {
	Type   string
	Origin error
}
```

ClientError formats errors coming from the GoFlix.

#### func (ClientError) Error

```go
func (clientError ClientError) Error() string
```

#### type GoFlix

```go
type GoFlix struct {
	Client   *torrent.Client
	Torrent  *torrent.Torrent
	Progress int64
	Uploaded int64
	Config   *ClientConfig
}
```

GoFlix manages the torrent downloading.

#### func  New

```go
func New(cfg *ClientConfig) (*GoFlix, error)
```
NewClient creates a new torrent GoFlix based on a magnet or a torrent file. If
the torrent file is on http, we try downloading it.

#### func (*GoFlix) Close

```go
func (c *GoFlix) Close()
```
Close cleans up the connections.

#### func (*GoFlix) FileReader

```go
func (c *GoFlix) FileReader() (io.ReadSeeker, error)
```

#### func (*GoFlix) GetLargestFile

```go
func (c *GoFlix) GetLargestFile() *torrent.File
```

#### func (*GoFlix) HandlerFunc

```go
func (c *GoFlix) HandlerFunc() http.HandlerFunc
```
GetFile is an http handler to serve the biggest file managed by the GoFlix.

#### func (*GoFlix) Percentage

```go
func (c *GoFlix) Percentage() float64
```

#### func (*GoFlix) ReadyForPlayback

```go
func (c *GoFlix) ReadyForPlayback() bool
```
ReadyForPlayback checks if the torrent is ready for playback or not. We wait
until 5% of the torrent to start playing.

#### func (*GoFlix) Render

```go
func (c *GoFlix) Render()
```
Render outputs the command line interface for the GoFlix.
