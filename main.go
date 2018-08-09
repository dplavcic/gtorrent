package main

import (
	_ "crypto/sha1"
	"fmt"

	"github.com/dplavcic/gtorrent/torrent"
)

func main() {

	t := torrent.ReadTorrentFile("debian.torrent")
	tr := torrent.MapToStruct(t)

	r := torrent.PingTracker(tr)
	fmt.Println("response: " + string(r))
	//	ioutil.WriteFile("debian_compact_mode", r, 066)
	torrent.ParseTrackerResponse(r)

}
