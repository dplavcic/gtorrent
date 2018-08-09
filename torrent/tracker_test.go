package torrent

import (
	"net/url"
	"strings"
	"testing"
)

func TestInfoHashEscape(t *testing.T) {
	ReadTorrentFile("../debian.torrent")
	h := CalculateAnnounceHash("../debian.torrent")

	got := url.PathEscape(string(h))
	want := strings.ToUpper("%3b%1d%85%f8x%0e%f8%c4%d8S%8f%80%9azc%fcR%991%8e")

	if want != got {
		t.Errorf("want: %s, got: %s\n", want, got)
	}

}

