package torrent

import (
	"bytes"
	"crypto"
	"io"
	"io/ioutil"
	"log"

	"github.com/dplavcic/gtorrent/bencode"
)

func CalculateAnnounceHash(fileName string) []byte {

	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Panic(err)
	}

	return calculateSha1(file[bencode.InfoDictStartPosition():bencode.InfoDictEndPosition()])
}

func calculateSha1(b []byte) []byte {
	buf := bytes.NewBuffer(b)
	h := crypto.SHA1.New()
	if _, err := io.Copy(h, buf); err != nil {
		log.Fatal(err)
	}
	return h.Sum(nil)
}
