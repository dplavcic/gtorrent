package bencode

import (
	"fmt"
	"io/ioutil"
	"testing"
)

type Torrent struct {
	Announce     string
	AnnounceList []interface{}
	Comment      string
	CreationDate int
	Info         Info
	HttpSeeds    []interface{}
}

type Info struct {
	Length      int
	Name        string
	PieceLength int
	Pieces      string
}

func TestUnmarshal(t *testing.T) {

	var tr Torrent
	tr.Announce = "Hi"

	// f, e := ioutil.ReadFile("../ubuntu_18.torrent.custom_field_names")
	f, e := ioutil.ReadFile("../debian.torrent.custom_field_names")
	if e != nil {
		t.Error(e)
	}

	fmt.Printf("%v\n", string(f))
	Unmarshal(f, &tr)

	t.Errorf("%v\n\n", tr)

}
