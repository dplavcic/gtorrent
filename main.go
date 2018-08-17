package main

import (
	_ "crypto/sha1"
	"fmt"
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

	fmt.Printf("%#v\n\n", string(data))
	var tr torrent.TorrentFile
	bencode.Unmarshal(data, &tr)

}
