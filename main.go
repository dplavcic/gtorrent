package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strconv"
)

type singleFileMode struct {
	name   string
	length int64
	md5sum string
}
type multipleFileMode struct {
	name  string
	files []singleFileMode
}

type infoDictionary struct {
	pieceLength int64 // number of bytes in each piece
	// string consisting of the concatenation of all 20-byye SHA1 hash values,
	// one per piece (byte string)
	pieces  string
	private int8 // 1 - true
}

type TorrentFile struct {
	info         string // TODO(dplavcic) check dir struct
	announce     string // the announce URL of the tracker
	announceList []string
	creationDate string //timestamp
	comment      string
	createdBy    string
	encoding     string
}

func main() {
	data, err := ioutil.ReadFile("listen.pls.torrent")
	if err != nil {
		log.Fatal(err)
	}
	buf := bytes.NewBuffer(data)
	dict, err := readTorrentFile(buf)
	fmt.Printf("%v\n\n", dict)
}

func readTorrentFile(buf *bytes.Buffer) (map[string]interface{}, error) {
	b, err := buf.ReadByte()
	if err != nil {
		log.Panic("Could not read byte")
	}
	if b == 'd' {
		return createDict(buf)
	}
	return nil, nil
}

func createDict(buf *bytes.Buffer) (map[string]interface{}, error) {
	dict := make(map[string]interface{})

	for {
		key, _ := readString(buf)
		value, _ := readNext(buf)
		dict[key] = value

		b, _ := buf.ReadByte()
		if b == 'e' {
			break
		} else {
			buf.UnreadByte()
		}
	}

	return dict, nil
}

func readNext(buf *bytes.Buffer) (interface{}, error) {
	var item interface{}
	b, err := buf.ReadByte()
	if err != nil {
		log.Panic(err)
	}

	if b == 'd' {
		item, err = createDict(buf)
	} else if b == 'i' {
		item, err = readInt(buf)
	} else if b == 'l' {
		item, err = readList(buf)
	} else {
		buf.UnreadByte()
		item, err = readString(buf)
	}

	fmt.Printf("returning: %v\n", item)
	return item, err
}

func readList(buf *bytes.Buffer) ([]interface{}, error) {
	var list []interface{}

	for {
		next, _ := buf.ReadByte()
		fmt.Println(string(next))
		fmt.Println(list)
		if next == 'e' {
			fmt.Println("break")
			break
		}

		buf.UnreadByte()
		value, _ := readNext(buf)
		list = append(list, value)
	}
	return list, nil
}

func readInt(buf *bytes.Buffer) (int64, error) {
	stringLine, err := buf.ReadString(byte('e'))
	if err != nil {
		return 0, err
	}
	value, err := strconv.ParseInt(stringLine[0:len(stringLine)-1], 10, 64)
	return value, err
}

func readString(buf *bytes.Buffer) (string, error) {
	strLen, _ := buf.ReadString(byte(':'))
	strLenValue, _ := strconv.ParseInt(strLen[0:len(strLen)-1], 10, 64)
	buffer := make([]byte, strLenValue)
	_, err := io.ReadFull(buf, buffer)
	return string(buffer), err
}
