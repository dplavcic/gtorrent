package bencode

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
)

type Decode struct {
	buf *bytes.Buffer
}

func Unmarshal(data []byte, v interface{}) error {
	buf := bytes.NewBuffer(data)
	d := Decode{buf: buf}
	return d.unmarshal(v)
}

func (d *Decode) unmarshal(v interface{}) (err error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("rv.Kind() != reflect.Ptr || rv.IsNil()")
	}

	d.parseValue(rv)
	return nil
}

func (d *Decode) parseValue(v reflect.Value) {
	d.readNext(v)
}

// v is *ptr
func (d *Decode) readNext(v reflect.Value) {
	b, err := d.buf.ReadByte()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		log.Panic()
	}
	switch b {
	case 'd':
		d.parseDict(v)
	case 'l':
		d.parseList(v)
	case 'i':
		d.parseInt(v)
	case 'e':
		return
	default: //string
		d.buf.UnreadByte()
		d.parseString(v)
	}
}

func (d *Decode) parseDict(v reflect.Value) {
	for {
		key := d.readKey()
		var fv reflect.Value
		if v.Kind() == reflect.Struct {
			fv = d.fieldName(v, key)
		} else {
			fv = d.fieldName(v.Elem(), key)
		}

		// this sets parsed value to v
		d.readNext(fv)
		b, err := d.buf.ReadByte()
		if err != nil {
			return
		}

		if b == 'e' { //end of dict
			return
		}

		//not end of dict
		d.buf.UnreadByte()
	}
}

// reads dict key, assume string
func (d *Decode) readKey() string {
	line, err := d.buf.ReadString(byte(':')) //read until :, : included
	if err != nil {
		log.Fatal(err)
	}

	line = line[0 : len(line)-1] //n - number of characters to read
	strLen, err := strconv.ParseInt(line, 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, strLen)
	n, err := io.ReadFull(d.buf, buffer)
	if err != nil {
		log.Fatal(err)
	}

	if int64(n) != strLen {
		log.Fatal("n != strLen")
	}

	return string(buffer)
}

func (d *Decode) parseList(v reflect.Value) {
	i := 0
	for ; ; i++ {
		switch v.Kind() {
		case reflect.Slice:
			v.Set(reflect.AppendSlice(v, reflect.MakeSlice(v.Type(), 1, 1)))
			d.readNext(v.Index(i))
		case reflect.Interface:
			// create new slice
			ns := make([]interface{}, 1, 1)
			nsv := reflect.ValueOf(&ns)
			// add elements to new slice
			d.readNext(nsv.Elem().Index(i))
			// set old to new slice
			v.Set(nsv.Elem())
		default:
			log.Panic(fmt.Sprintf("kind: %v", v.Kind()))
		}

		b, err := d.buf.ReadByte()
		if err != nil {
			log.Fatal(err)
		}

		if b == 'e' { //end of list
			return
		}
		d.buf.UnreadByte()
	}
}

func (d *Decode) parseInt(v reflect.Value) {

	line, err := d.buf.ReadString(byte('e'))
	if err != nil {
		log.Fatal(err)
	}
	n, err := strconv.ParseInt(line[0:len(line)-1], 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	v.SetInt(n)
}

func (d *Decode) parseString(v reflect.Value) {
	value := d.readKey()
	v.Set(reflect.ValueOf(value))
}

func (d *Decode) fieldName(v reflect.Value, key string) reflect.Value {

	val := reflect.Indirect(v)
	if val.Kind() != reflect.Struct {
		log.Panic("unmarshall can only take structs")
	}

	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := fieldName(f)
		if tag == key {
			fbi := v.FieldByIndex([]int{i})
			return fbi
		}
	}
	return reflect.ValueOf("")
}

func fieldName(f reflect.StructField) string {
	if t := f.Tag.Get("bencode"); t != "" {
		return t //return tag
	}
	return f.Name //fall back to field name
}
