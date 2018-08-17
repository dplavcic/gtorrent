package torrent

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

const (
	EventStarted   = "started"
	EventCompleted = "completed"
	EventStopped   = "stopped"
)

type TrackerRequest struct {
	announceURL string
	infoHash    []byte
	peerID      string
	port        int
	uploaded    int
	downloaded  int
	left        int64
	events      string
}

type TrackerResponse struct {
	failureReason string
	interval      int
	peers         []Peer
}

type Peer struct {
	id   string
	ip   string
	port int
}

func PingTracker(t *TorrentFile) []byte {
	tr := &TrackerRequest{
		announceURL: t.Announce,
		infoHash:    CalculateAnnounceHash("debian.torrent"),
		peerID:      "-TR2920-wbhc3m277yr6",
		port:        51413,
		uploaded:    0,
		downloaded:  0,
		left:        t.Info.Length,
		events:      EventStarted,
	}

	u, err := url.Parse(tr.announceURL)
	if err != nil {
		log.Fatal(err)
	}
	u.Scheme = "http"
	q := u.Query()
	q.Set("info_hash", string(tr.infoHash))
	q.Set("peer_id", tr.peerID)
	q.Set("port", strconv.Itoa(tr.port))
	q.Set("uploaded", "0")
	q.Set("downloaded", strconv.Itoa(tr.downloaded))
	q.Set("left", strconv.FormatInt(tr.left, 10))
	q.Set("event", tr.events)
	q.Set("compact", strconv.Itoa(0))

	u.RawQuery = q.Encode()
	r, err := http.Get(u.String())

	if err != nil {
		log.Panic(err)
	}

	defer r.Body.Close()
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	return content
}

func ParseTrackerResponse(b []byte) {
	// br := bytes.NewBuffer(b)
	// r := bencode.Unmarshall(br)
	fmt.Printf("%#v\n\n", nil)
}
