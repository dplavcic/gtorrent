package bencode

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
)

// TODO(dplavcic) find better way to calculate announce hash
// announce hash helper variable
//do not uomarshal bztes at this step..
// /unmarshal later, save šđbzte to
var originalBuferLength int
var infoDictStartPosition int
var infoDictEndPosition int
var announceURLHash = make([]byte, 160)

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
	d.readNext(v)
	return nil
}

func (d *Decoder) readNext(v reflect.Value) (interface{}, error) {
	b, _ := d.buf.ReadByte()
	var val interface{}
	var err error

	if b == 'd' {
		d.readDict(v)
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
		if v.Kind() == reflect.Struct {
			// key, value
			// assume key value is a string
			key, err := d.readKey()
			field := v.FieldByName(key)
			// fmt.Printf("field name: %v\n", field)
			// field.SetInt(10)

			// fmt.Printf("key v: %v\n", v)
			_, err = d.readNext(field)
			// fmt.Printf("val v: %v\n", v)

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
		} else {
			// key, value
			// assume key value is a string
			key, err := d.readKey()
			field := v.Elem().FieldByName(key)
			// fmt.Printf("field name: %v\n", field)
			// field.SetInt(10)

			// fmt.Printf("key v: %v\n", v)
			_, err = d.readNext(field)
			// fmt.Printf("val v: %v\n", v)

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

	}
	return nil

}

func (d *Decoder) readKey() (string, error) {
	strLen, err := d.buf.ReadString(byte(':'))
	if err != nil {
		return "", err
	}
	strLenValue, err := strconv.ParseInt(strLen[0:len(strLen)-1], 10, 64)
	if err != nil {
		return "", err
	}
	buffer := make([]byte, strLenValue)
	n, err := io.ReadFull(d.buf, buffer)
	if err != nil {
		return "", err
	}
	if int64(n) != strLenValue {
		return "", fmt.Errorf("could not read all bytes. expected: %d, got: %d", strLenValue, n)
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
			// if i == v.Elem().Len() {
			// if v.Elem().CanSet() {
			// v.Elem().Set(reflect.Append(v.Elem(), reflect.ValueOf(make([]interface{}, 0, 0))))
			// } else {
			v.Set(reflect.Append(v, reflect.ValueOf(make([]interface{}, 0, 0))))
			// }
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
	strLen, err := d.buf.ReadString(byte(':'))
	if err != nil {
		return "", err
	}
	strLenValue, err := strconv.ParseInt(strLen[0:len(strLen)-1], 10, 64)
	if err != nil {
		return "", err
	}
	buffer := make([]byte, strLenValue)
	n, err := io.ReadFull(d.buf, buffer)
	if err != nil {
		return "", err
	}
	if int64(n) != strLenValue {
		return "", fmt.Errorf("could not read all bytes. expected: %d, got: %d", strLenValue, n)
	}
	fmt.Printf("buffer: %#v\n", buffer)
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

func printAddress(i interface{}) {
	fmt.Printf("addr: %p\n", i)
}

func InfoDictStartPosition() int {
	return infoDictStartPosition
}

func InfoDictEndPosition() int {
	return infoDictEndPosition
}
