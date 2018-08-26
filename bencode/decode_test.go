package bencode

import (
	"bytes"
	"crypto"
	_ "crypto/sha1"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"testing"
)

// todo(dplavcic) ignore unknown fields
type Torrent struct {
	Announce     string        `bencode:"announce"`
	AnnounceList []interface{} `bencode:"announce-list"`
	Comment      string        `bencode:"comment"`
	CreationDate int           `bencode:"creation date"`
	Info         []byte        `bencode:"info"`
	HTTPSeeds    []interface{} `bencode:"http seeds"`
	Checksum     string        `bencode:"checksum"`
	CreatedBy    string        `bencode:"created by"`
	Encoding     string        `bencode:"encoding"`
}

type Info struct {
	Name        string `bencode:"name"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
	Length      int    `bencode:"length"` //length of the file in bytes
	Files       []File `bencode:"files"`
	Private     int    `bencode:"private"`
	Entropy     string `bencode:"entropy"`
	Source      string `bencode:"source"`
	XCrossSeed  string `bencode:"x cross seed"`
}

type File struct {
	Length int      `bencode:"length"`
	Path   []string `bencode:"path"`
	Name   string   `bencode:"name"`
}

func TestUnmarshalByteField(t *testing.T) {
	f, e := ioutil.ReadFile("../ubuntu_server_info_dict.torrent")
	if e != nil {
		t.Error(e)
	}
	var infoDict Info
	Unmarshal(f, &infoDict)

	t.Errorf("info dict: %#v\n", infoDict)
}

func TestUnmarshal(t *testing.T) {
	var tr Torrent
	f, e := ioutil.ReadFile("../ubuntu_server.torrent")
	if e != nil {
		t.Error(e)
	}

	Unmarshal(f, &tr)
	t.Errorf("%v\n", tr)
}

func TestInfoDictHash(t *testing.T) {
	var tr Torrent
	f, e := ioutil.ReadFile("../ubuntu_server.torrent")
	if e != nil {
		t.Error(e)
	}

	Unmarshal(f, &tr)
	h := calculateSha1(tr.Info)

	got := url.PathEscape(string(h))
	want := strings.ToUpper("%965%ca%bd%20%1a%92%fd%a5%d6%c7%0d5%cb%fbp~%29%05%89")

	if strings.ToLower(want) != strings.ToLower(got) {
		t.Errorf("want: %s, got: %s\n", want, got)
	}
}

func calculateSha1(b []byte) []byte {
	buf := bytes.NewBuffer(b)
	h := crypto.SHA1.New()
	if _, err := io.Copy(h, buf); err != nil {
		log.Fatal(err)
	}
	return h.Sum(nil)
}
