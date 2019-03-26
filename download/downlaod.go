package download

import (
	"compress/gzip"
	"github.com/anacrolix/torrent/iplist"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const torrentBlockListURL = "http://john.bitsurge.net/public/biglist.p2p.gz"

// Download and add the blocklist.
func GetBlocklist() iplist.Ranger {
	var err error
	blocklistPath := os.TempDir() + "/go-peerflix-blocklist.gz"

	if _, err = os.Stat(blocklistPath); os.IsNotExist(err) {
		err = DownloadBlockList(blocklistPath)
	}

	if err != nil {
		log.Printf("Error downloading blocklist: %s", err)
		return nil
	}

	// Load blocklist.
	// #nosec
	// We trust our temporary directory as we just wrote the file there ourselves.
	blocklistReader, err := os.Open(blocklistPath)
	if err != nil {
		log.Printf("Error opening blocklist: %s", err)
		return nil
	}

	// Extract file.
	gzipReader, err := gzip.NewReader(blocklistReader)
	if err != nil {
		log.Printf("Error extracting blocklist: %s", err)
		return nil
	}

	// Read as iplist.
	blocklist, err := iplist.NewFromReader(gzipReader)
	if err != nil {
		log.Printf("Error reading blocklist: %s", err)
		return nil
	}

	log.Printf("Loading blocklist.\nFound %d ranges\n", blocklist.NumRanges())
	return blocklist
}

func DownloadBlockList(blocklistPath string) (err error) {
	log.Printf("Downloading blocklist")
	fileName, err := DownloadFile(torrentBlockListURL)
	if err != nil {
		log.Printf("Error downloading blocklist: %s\n", err)
		return
	}

	return os.Rename(fileName, blocklistPath)
}

func DownloadFile(URL string) (fileName string, err error) {
	var file *os.File
	if file, err = ioutil.TempFile(os.TempDir(), "go-peerflix"); err != nil {
		return
	}

	defer func() {
		if ferr := file.Close(); ferr != nil {
			log.Printf("Error closing torrent file: %s", ferr)
		}
	}()

	// #nosec
	// We are downloading the url the user passed to us, we trust it is a torrent file.
	response, err := http.Get(URL)
	if err != nil {
		return
	}

	defer func() {
		if ferr := response.Body.Close(); ferr != nil {
			log.Printf("Error closing torrent file: %s", ferr)
		}
	}()

	_, err = io.Copy(file, response.Body)

	return file.Name(), err
}
