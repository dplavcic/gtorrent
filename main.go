package main

import (
	_ "crypto/sha1"
	"fmt"

	"github.com/dplavcic/gtorrent/torrent"
)

func main() {

	torrent.ReadTorrentFile("debian.torrent")
	h := torrent.CalculateAnnounceHash("debian.torrent")

	fmt.Printf("%x\n", h)
}
