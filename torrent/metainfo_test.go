package torrent

import (
	_ "crypto/sha1"
	"encoding/hex"
	"testing"
)

func TestCalculateAnnounceSHA1(t *testing.T) {
	ReadTorrentFile("../debian.torrent")
	h := CalculateAnnounceHash("../debian.torrent")

	if hex.EncodeToString(h) != "3b1d85f8780ef8c4d8538f809a7a63fc5299318e" {
		t.Errorf("Got: %s, expected: %s\n", hex.Dump(h), "3b1d85f8780ef8c4d8538f809a7a63fc5299318e")
	}

	

}
