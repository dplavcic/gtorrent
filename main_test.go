package main

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"testing"
)

func TestReadInt(t *testing.T) {
	bNumber := []byte("i4194304e")
	buf := bytes.NewBuffer(bNumber)
	result, err := readInt(buf)
	if err != nil {
		t.Error(err)
	}

	if result != 4194304 {
		t.Errorf("Expected %v, got: %v\n", 1, result)
	}
}

func TestReadString(t *testing.T) {
	b := []byte("39:http://torrent.ubuntu.com:6969/announce")
	buf := bytes.NewBuffer(b)
	result, err := readString(buf)
	if err != nil {
		t.Error(err)
	}

	if result != "http://torrent.ubuntu.com:6969/announce" {
		t.Errorf("Expected %v, got: %v\n", "http://torrent.ubuntu.com:6969/announce", result)
	}
}

func TestReadString2(t *testing.T) {
	b := []byte("39:http://torrent.ubuntu.com:6969/announce")
	buf := bytes.NewBuffer(b)
	strLen, _ := buf.ReadString(byte(':'))
	if strLen[:len(strLen)-1] != "39" {
		t.Error("string len: " + strLen)
	}

	strLenValue, _ := strconv.ParseInt(strLen[0:len(strLen)-1], 10, 64)
	t.Error(strLenValue)
}

func TestCreateDict(t *testing.T) {
	b, err := ioutil.ReadFile("ubuntu.torrent")
	if err != nil {
		t.Error(err)
	}
	buf := bytes.NewBuffer(b)

	res, _ := readTorrentFile(buf)
	t.Errorf("%#v\n", res)
}

func TestCreateList(t *testing.T) {
	b := []byte("l4:spam4:eggse")
	buf := bytes.NewBuffer(b)

	buf.ReadByte()
	r, err := readList(buf)
	t.Error(r)
	if err != nil {
		t.Error(err)
	}
}
