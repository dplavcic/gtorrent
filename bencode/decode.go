package bencode

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
)

type Decoder struct {
	buf *bytes.Buffer
	v   interface{}
}

func Unmarshal(data []byte, v interface{}) (err error) {
	buf := bytes.NewBuffer(data)
	e := Decoder{buf: buf, v: v}
	err = e.Decode(v)
	return err
}

func (d *Decoder) Decode(v interface{}) error {
	pv := reflect.ValueOf(v)
	if pv.Kind() != reflect.Ptr || pv.IsNil() {
		return fmt.Errorf("invalid arg reveived: %s", reflect.TypeOf(v))
	}
	return d.parseValue(pv)
}

func (d *Decoder) parseValue(v reflect.Value) error {
	_, err := d.readNext(v)
	if err != nil {
		return err
	}
	return nil
}

func (d *Decoder) readNext(v reflect.Value) (interface{}, error) {
	b, _ := d.buf.ReadByte()
	var val interface{}
	var err error

	if b == 'd' {
		err = d.readDict(v)
	} else if b == 'i' {
		val, err = d.readInt(v)
	} else if b == 'l' {
		val, err = d.readList(v)
	} else {
		d.buf.UnreadByte()
		val, err = d.readString(v)
	}
	return val, err
}

func (d *Decoder) readDict(v reflect.Value) error {

	for {
		// assume key value is a string
		key, err := d.readKey()
		if err != nil {
			log.Fatal(err)
		}

		var field reflect.Value
		if v.Kind() == reflect.Struct {
			field = v.FieldByName(key)
		} else {
			field = v.Elem().FieldByName(key)

		}

		_, err = d.readNext(field)
		if err != nil {
			fmt.Println(err)
		}

		b, err := d.buf.ReadByte()
		d.buf.UnreadByte()
		if err != nil {
			return err
		}

		if b == 'e' { // dict end
			break
		}
	}
	return nil
}

func (d *Decoder) readKey() (string, error) {
	strCharCount, err := d.buf.ReadString(byte(':'))
	if err != nil {
		return "", err
	}
	strCharCount = strCharCount[0 : len(strCharCount)-1]
	strLen, err := strconv.ParseInt(strCharCount, 10, 64)
	if err != nil {
		return "", err
	}

	buffer := make([]byte, strLen)
	n, err := io.ReadFull(d.buf, buffer)
	if err != nil {
		return "", err
	}
	if int64(n) != strLen {
		return "", fmt.Errorf("could not read all bytes. expected: %d, got: %d", strLen, n)
	}

	return string(buffer), nil
}

func (d *Decoder) readList(v reflect.Value) (interface{}, error) {

	v.Set(reflect.Append(v, reflect.ValueOf(make([]interface{}, 0, 0))))

	i := 0
	for ; ; i++ {
		b, err := d.buf.ReadByte()
		if err != nil {
			return nil, err
		}
		if b == 'e' { //end of list
			break
		}

		if i == v.Len() {
			v.Set(reflect.Append(v, reflect.ValueOf(make([]interface{}, 0, 0))))
		}

		d.buf.UnreadByte() //not end, unread one byte

		_, err = d.readNext(v.Index(i))
		if err != nil {
			return nil, err
		}
	}
	return v, nil
}

func (d *Decoder) readString(v reflect.Value) (interface{}, error) {
	strCharCount, err := d.buf.ReadString(byte(':'))
	if err != nil {
		return "", err
	}
	strCharCount = strCharCount[0 : len(strCharCount)-1]
	strLen, err := strconv.ParseInt(strCharCount, 10, 64)
	if err != nil {
		return "", err
	}

	buffer := make([]byte, strLen)
	n, err := io.ReadFull(d.buf, buffer)
	if err != nil {
		return "", err
	}
	if int64(n) != strLen {
		return "", fmt.Errorf("could not read all bytes. expected: %d, got: %d", strLen, n)
	}

	switch v.Kind() {
	case reflect.String:
		v.Set(reflect.ValueOf(string(buffer)))
	default:
		v.Set(reflect.ValueOf(buffer))

	}

	return string(buffer), nil
}

func (d *Decoder) readInt(v reflect.Value) (int64, error) {
	stringLine, err := d.buf.ReadString(byte('e'))
	if err != nil {
		return 0, err
	}
	value, err := strconv.ParseInt(stringLine[0:len(stringLine)-1], 10, 64)
	if err != nil {
		return 0, err
	}

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(value)
	case reflect.Interface:
		v.Set(reflect.ValueOf(v))
	default:
		log.Panic("read int")
	}
	return value, nil
}
