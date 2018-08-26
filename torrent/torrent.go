package torrent

type Torrent struct {
	Announce     string        `bencode:"announce"`
	AnnounceList []interface{} `bencode:"announce-list"`
	Comment      string        `bencode:"comment"`
	CreationDate int           `bencode:"creation date"`
	InfoByte     []byte        `bencode:"info"`
	Info         Info
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
