package client

import (
	"bufio"
	"fmt"
	"github.com/autom8ter/goflix/download"
	"github.com/autom8ter/goflix/reader"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/dustin/go-humanize"
)

const clearScreen = "\033[H\033[2J"

var isHTTP = regexp.MustCompile(`^https?:\/\/`)

// ClientError formats errors coming from the client.
type ClientError struct {
	Type   string
	Origin error
}

func (clientError ClientError) Error() string {
	return fmt.Sprintf("Error %s: %s\n", clientError.Type, clientError.Origin)
}

// Client manages the torrent downloading.
type Client struct {
	Client   *torrent.Client
	Torrent  *torrent.Torrent
	Progress int64
	Uploaded int64
	Config   *ClientConfig
}

// ClientConfig specifies the behaviour of a client.
type ClientConfig struct {
	TorrentPath    string
	Port           int
	TorrentPort    int
	Seed           bool
	TCP            bool
	MaxConnections int
}

func DefaultClientConfig() ClientConfig {
	return ClientConfig{
		Port:           8080,
		TorrentPort:    50007,
		Seed:           false,
		TCP:            true,
		MaxConnections: 200,
	}
}

// NewClient creates a new torrent client based on a magnet or a torrent file.
// If the torrent file is on http, we try downloading it.
func New(cfg *ClientConfig) (*Client, error) {
	var (
		client = &Client{}
		t      *torrent.Torrent
		c      *torrent.Client
		err    error
	)
	client.Config = cfg
	blocklist := download.GetBlocklist()
	torrentConfig := torrent.NewDefaultClientConfig()
	torrentConfig.DataDir = os.TempDir()
	torrentConfig.NoUpload = !cfg.Seed
	torrentConfig.DisableTCP = !cfg.TCP
	torrentConfig.ListenPort = cfg.TorrentPort
	torrentConfig.IPBlocklist = blocklist

	// Create client.
	c, err = torrent.NewClient(torrentConfig)

	if err != nil {
		return client, ClientError{Type: "creating torrent client", Origin: err}
	}

	client.Client = c

	// Add torrent.

	// Add as magnet url.
	if strings.HasPrefix(cfg.TorrentPath, "magnet:") {
		if t, err = c.AddMagnet(cfg.TorrentPath); err != nil {
			return client, ClientError{Type: "adding torrent", Origin: err}
		}
	} else {
		// Otherwise add as a torrent file.

		// If it's online, we try downloading the file.
		if isHTTP.MatchString(cfg.TorrentPath) {
			if cfg.TorrentPath, err = download.DownloadFile(cfg.TorrentPath); err != nil {
				return client, ClientError{Type: "downloading torrent file", Origin: err}
			}
		}

		if t, err = c.AddTorrentFromFile(cfg.TorrentPath); err != nil {
			return client, ClientError{Type: "adding torrent to the client", Origin: err}
		}
	}

	client.Torrent = t
	client.Torrent.SetMaxEstablishedConns(cfg.MaxConnections)

	go func() {
		<-t.GotInfo()
		t.DownloadAll()

		// Prioritize first 5% of the file.
		largestFile := client.GetLargestFile()
		firstPieceIndex := largestFile.Offset() * int64(t.NumPieces()) / t.Length()
		endPieceIndex := (largestFile.Offset() + largestFile.Length()) * int64(t.NumPieces()) / t.Length()
		for idx := firstPieceIndex; idx <= endPieceIndex*5/100; idx++ {
			t.Piece(int(idx)).SetPriority(torrent.PiecePriorityNow)
		}
	}()

	return client, err
}

// Close cleans up the connections.
func (c *Client) Close() {
	c.Torrent.Drop()
	c.Client.Close()
}

// Render outputs the command line interface for the client.
func (c *Client) Render() {
	t := c.Torrent

	if t.Info() == nil {
		return
	}

	currentProgress := t.BytesCompleted()
	downloadSpeed := humanize.Bytes(uint64(currentProgress-c.Progress)) + "/s"
	c.Progress = currentProgress

	complete := humanize.Bytes(uint64(currentProgress))
	size := humanize.Bytes(uint64(t.Info().TotalLength()))

	bytesWrittenData := t.Stats().BytesWrittenData
	uploadProgress := (&bytesWrittenData).Int64() - c.Uploaded
	uploadSpeed := humanize.Bytes(uint64(uploadProgress)) + "/s"
	c.Uploaded = uploadProgress
	percentage := c.Percentage()
	totalLength := t.Info().TotalLength()

	output := bufio.NewWriter(os.Stdout)

	fmt.Fprint(output, clearScreen)
	fmt.Fprint(output, t.Info().Name+"\n")
	fmt.Fprint(output, strings.Repeat("=", len(t.Info().Name))+"\n")
	if c.ReadyForPlayback() {
		fmt.Fprintf(output, "Stream: \thttp://localhost:%d\n", c.Config.Port)
	}
	if currentProgress > 0 {
		fmt.Fprintf(output, "Progress: \t%s / %s  %.2f%%\n", complete, size, percentage)
	}
	if currentProgress < totalLength {
		fmt.Fprintf(output, "Download speed: %s\n", downloadSpeed)
	}
	if c.Config.Seed {
		fmt.Fprintf(output, "Upload speed: \t%s", uploadSpeed)
	}

	output.Flush()
}

//
func (c *Client) GetLargestFile() *torrent.File {
	var target *torrent.File
	var maxSize int64

	for _, file := range c.Torrent.Files() {
		if maxSize < file.Length() {
			maxSize = file.Length()
			target = file
		}
	}

	return target
}

// ReadyForPlayback checks if the torrent is ready for playback or not.
// We wait until 5% of the torrent to start playing.
func (c *Client) ReadyForPlayback() bool {
	return c.Percentage() > 5
}

// GetFile is an http handler to serve the biggest file managed by the client.
func (c *Client) GetFile(w http.ResponseWriter, r *http.Request) {
	target := c.GetLargestFile()
	entry, err := reader.NewFileReader(target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer func() {
		if err := entry.Close(); err != nil {
			log.Printf("Error closing file reader: %s\n", err)
		}
	}()

	w.Header().Set("Content-Disposition", "attachment; filename=\""+c.Torrent.Info().Name+"\"")
	http.ServeContent(w, r, target.DisplayPath(), time.Now(), entry)
}

func (c *Client) Percentage() float64 {
	info := c.Torrent.Info()

	if info == nil {
		return 0
	}

	return float64(c.Torrent.BytesCompleted()) / float64(info.TotalLength()) * 100
}

/*
func (c Client) RenderPieces() (output string) {
	pieces := c.Torrent.PieceStateRuns()
	for i := range pieces {
		piece := pieces[i]
		if piece.Priority == torrent.PiecePriorityReadahead {
			output += "!"
		}
		if piece.Partial {
			output += "P"
		} else if piece.Checking {
			output += "c"
		} else if piece.Complete {
			output += "d"
		} else {
			output += "_"
		}
	}
	return
}
*/
