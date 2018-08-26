package torrent

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/dplavcic/gtorrent/bencode"
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
	compact     int
}

type TrackerResponse struct {
	FailureReason string `bencode:"failureReason"`
	Interval      int    `bencode:"interval"`
	PeerData      string `bencode:"peers"`
	Peers         []Peer
}

// comapct form - each peer is represented using only 6 bytes
// first 4 byte contain the 32-bit ipv4 address
// last 2 byte contain the port number
// network byte order is used
type Peer struct {
	// does not appear in compact mode
	ID   string `bencode:"peer id"` // 20 bytes + 3 bytes bencoding overhead
	IP   net.IP `bencode:"ip"`      // max 225 bytes + 4bytes bencoding overhead
	Port uint16 `bencode:"port"`    // max 7 bytes +
}

func PingTracker(t Torrent) []byte {
	tr := TrackerRequest{
		announceURL: t.Announce,
		infoHash:    CalculateSha1(t.InfoByte),
		peerID:      "-TR2920-wbhc3m277yr6",
		port:        51413,
		uploaded:    0,
		downloaded:  0,
		left:        int64(0),
		events:      EventStarted,
		compact:     1,
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

func ParseTrackerResponse(b []byte) TrackerResponse {
	var tr TrackerResponse
	bencode.Unmarshal(b, &tr)
	tr.Peers = ParsePeerList(tr)
	return tr
}

func ParsePeerList(tr TrackerResponse) []Peer {
	var pl []Peer
	size := len(tr.PeerData)
	for i := 0; i < size; i += 6 {
		ip := tr.PeerData[i : i+6]
		ipp := net.IPv4(ip[0], ip[1], ip[2], ip[3])
		port := []byte{ip[4], ip[5]}
		p := binary.BigEndian.Uint16(port)

		peer := Peer{
			IP:   ipp,
			Port: p,
		}
		pl = append(pl, peer)
	}
	return pl
}
