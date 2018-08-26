package torrent

import (
	"bytes"
	"crypto"
	"io"
	"log"
)

func CalculateSha1(b []byte) []byte {
	buf := bytes.NewBuffer(b)
	h := crypto.SHA1.New()
	if _, err := io.Copy(h, buf); err != nil {
		log.Fatal(err)
	}
	return h.Sum(nil)
}
