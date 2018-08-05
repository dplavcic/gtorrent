package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/dplavcic/gtorrent/bencode"
)

type singleFileMode struct {
	name   string
	length int64
	md5sum string
}
type multipleFileMode struct {
	name  string
	files []singleFileMode
}

type infoDictionary struct {
	pieceLength int64 // number of bytes in each piece
	// string consisting of the concatenation of all 20-byye SHA1 hash values,
	// one per piece (byte string)
	pieces  string
	private int8 // 1 - true
}

type TorrentFile struct {
	info         string // TODO(dplavcic) check dir struct
	announce     string // the announce URL of the tracker
	announceList []string
	creationDate string //timestamp
	comment      string
	createdBy    string
	encoding     string
}

func main() {
	dict := readTorrentFile("listen.pls.torrent")
	fmt.Printf("%v\n\n", dict)
}

func readTorrentFile(fileName string) interface{} {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	buf := bytes.NewBuffer(data)
	return bencode.Unmarshall(buf)
}
