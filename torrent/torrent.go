package torrent

import (
	"bytes"
	"io/ioutil"
	"log"
	"time"
)

type InfoDictionary struct {
	Name        string
	PieceLength int64
	Pieces      string
	Length      int64
	Path        string
	Private     bool
}

type TorrentFile struct {
	Announce     string `bencode:"announce"`
	AnnounceList []string
	Comment      string         `bencode:"comment"`
	CreatedBy    string         `bencode:"createdBy"`
	CreationDate time.Time      `bencode:"creation date"`
	Encoding     string         `bencode:"encoding"`
	Info         InfoDictionary // TODO(dplavcic) check dir struct
	InfoBytes    []byte         `bencode:"info"`
}

func ReadTorrentFile(fileName string) interface{} {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	bytes.NewBuffer(data)

	// return bencode.Unmarshall(buf)
	return nil
}

func MapToStruct(data interface{}) *TorrentFile {
	tf := &TorrentFile{
		Announce:     getString(data, "announce"),
		Comment:      getString(data, "comment"),
		CreatedBy:    getString(data, "created by"),
		CreationDate: getDate(data, "creation date"),
		Encoding:     getString(data, "encoding"),
		AnnounceList: nil,
		Info: InfoDictionary{
			Name:        getString(data.(map[string]interface{})["info"], "name"),
			PieceLength: getInt64(data.(map[string]interface{})["info"], "piece length"),
			Pieces:      getString(data.(map[string]interface{})["info"], "pieces"),
			Length:      getInt64(data.(map[string]interface{})["info"], "length"),
			Private:     getBool(data.(map[string]interface{})["info"], "private"),
		},
	}

	return tf
}

func getString(data interface{}, value string) string {
	r, ok := data.(map[string]interface{})[value].(string)
	if !ok {
		return ""
	}
	return r
}

func getInt64(data interface{}, value string) int64 {
	r, ok := data.(map[string]interface{})[value].(int64)
	if !ok {
		return -1
	}
	return r
}

func getDate(data interface{}, value string) time.Time {
	r, ok := data.(map[string]interface{})[value].(int64)
	if !ok {
		return time.Now()
	}
	return time.Unix(r, 0)
}

func getBool(data interface{}, value string) bool {
	r, _ := data.(map[string]interface{})[value].(int64)
	if r == 1 {
		return true
	}
	return false
}
