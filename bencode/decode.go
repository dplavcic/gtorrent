package bencode

import (
	"bytes"
	"io"
	"log"
	"strconv"
)

// TODO(dplavcic) find better way to calculate announce hash
// announce hash helper variable
var originalBuferLength int
var infoDictStartPosition int
var infoDictEndPosition int
var announceURLHash = make([]byte, 160)

func Unmarshall(buf *bytes.Buffer) interface{} {
	originalBuferLength = buf.Len()
	return ReadNext(buf)
}

func ReadNext(buf *bytes.Buffer) interface{} {
	var item interface{}
	b, err := buf.ReadByte()
	if err != nil {
		log.Fatal(err)
	}

	if b == 'd' {
		buf.UnreadByte()
		item = ReadDict(buf)
	} else if b == 'i' {
		buf.UnreadByte()
		item = ReadInt(buf)
	} else if b == 'l' {
		buf.UnreadByte()
		item = ReadList(buf)
	} else {
		buf.UnreadByte()
		item = ReadString(buf)
	}

	return item
}

func ReadDict(buf *bytes.Buffer) map[string]interface{} {
	dict := make(map[string]interface{})
	buf.ReadByte() ////ignore prefix 'd - list'

	for {
		key := ReadString(buf)
		value := ReadNext(buf)
		dict[key] = value

		b, err := buf.ReadByte()
		if err != nil {
			log.Fatal(err)
		}
		if b == 'e' {
			if infoDictEndPosition == 0 {
				infoDictEndPosition = (originalBuferLength - buf.Len())
			}
			break
		} else {
			buf.UnreadByte() //not end of the dict, unread byte
		}
	}

	return dict
}

func ReadList(buf *bytes.Buffer) []interface{} {
	var list []interface{}
	buf.ReadByte() //ignore prefix 'l - list'
	for {
		next, err := buf.ReadByte()
		if err != nil {
			log.Fatal(err)
		}
		if next == 'e' {
			break
		}

		buf.UnreadByte() //not end of the list, unread byte
		value := ReadNext(buf)
		list = append(list, value)
	}
	return list
}

func ReadInt(buf *bytes.Buffer) int64 {
	buf.ReadByte() // ignore prefix 'i - integer'
	stringLine, err := buf.ReadString(byte('e'))
	if err != nil {
		log.Fatal(err)
	}
	value, err := strconv.ParseInt(stringLine[0:len(stringLine)-1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return value
}

func ReadString(buf *bytes.Buffer) string {
	strLen, err := buf.ReadString(byte(':'))
	if err != nil {
		log.Fatal(err)
	}
	strLenValue, err := strconv.ParseInt(strLen[0:len(strLen)-1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	buffer := make([]byte, strLenValue)
	n, err := io.ReadFull(buf, buffer)
	if err != nil {
		log.Fatal(err)
	}
	if int64(n) != strLenValue {
		log.Fatalf("Could not read all bytes. Expected: %d, got: %d\n", strLenValue, n)
	}

	// we need this to calculate announce hash
	if string(buffer) == "info" {
		infoDictStartPosition = originalBuferLength - buf.Len()
	}

	return string(buffer)
}

func InfoDictStartPosition() int {
	return infoDictStartPosition
}

func InfoDictEndPosition() int {
	return infoDictEndPosition
}
