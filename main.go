package main

import (
	_ "crypto/sha1"
	"io/ioutil"
	"log"
	"os"

	"github.com/dplavcic/gtorrent/bencode"
	"github.com/dplavcic/gtorrent/torrent"
)

func main() {

	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	var tr torrent.Torrent
	bencode.Unmarshal(data, &tr)

	var info torrent.Info
	bencode.Unmarshal(tr.InfoByte, &info)
	tr.Info = info

	r := torrent.PingTracker(tr)
	torrent.ParseTrackerResponse(r)

}
