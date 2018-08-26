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
	r := d.unmarshal(v)
	return r
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
	var i []uint8
	if v.Kind() == reflect.Slice && v.Type() == reflect.TypeOf(i) {
		d.buf.UnreadByte()
		d.parseByteArray(v)
	}

	for {
		b, err := d.buf.ReadByte()
		if err != nil {
			return
		}

		if b == 'e' { //end of dict
			return
		}

		//not end of dict
		d.buf.UnreadByte()
		key := d.readKey()

		var fv reflect.Value
		if v.Kind() == reflect.Struct {
			fv, err = d.fieldName(v, key)
		} else {
			fv, err = d.fieldName(v.Elem(), key)
		}

		//ignore unknown fields
		if err != nil {
			d.skipValue(fv)
			continue
		}

		// this sets parsed value to v
		d.readNext(fv)
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

func (d *Decode) fieldName(v reflect.Value, key string) (reflect.Value, error) {

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
			return fbi, nil
		}
	}

	return reflect.ValueOf(nil), fmt.Errorf("could not find field tag by key: %v", key)
}

func fieldName(f reflect.StructField) string {
	if t := f.Tag.Get("bencode"); t != "" {
		return t //return tag
	}
	return f.Name //fall back to field name
}

// skip
func (d *Decode) skipValue(v reflect.Value) {
	b, err := d.buf.ReadByte()
	if err != nil {
		fmt.Println("could not skip value: ", string(b))
	}
	switch b {
	case 'd':
		d.skipDict(v)
	case 'l':
		d.skipList(v)
	case 'i':
		d.skipInt(v)
	default:
		d.buf.UnreadByte()
		d.skipString(v)
	}
}

func (d *Decode) skipDict(v reflect.Value) {
	for {
		//d.readKey()

		b, err := d.buf.ReadByte()
		if err != nil {
			log.Panic()
		}

		if b == 'e' {
			return //end of dict
		}
		d.buf.UnreadByte()
		d.skipValue(v)
	}
}

func (d *Decode) skipInt(v reflect.Value) {
	_, err := d.buf.ReadString(byte('e'))
	if err != nil {
		log.Fatal(err)
	}
}

func (d *Decode) skipString(v reflect.Value) {
	d.readKey()
}

func (d *Decode) skipList(v reflect.Value) {
	for {
		b, err := d.buf.ReadByte()
		if err != nil {
			fmt.Println("could not skip list: ", err)
		}

		if b == 'e' {
			return //end of list
		}
		d.buf.UnreadByte()
		d.skipValue(v)
	}
}

func (d *Decode) parseByteArray(v reflect.Value) {
	d.readNextByte(v)
}

func (d *Decode) readNextByte(v reflect.Value) {
	b, err := d.buf.ReadByte()
	if err != nil {
		fmt.Errorf("%v", err)
	}
	switch b {
	case 'd':
		d.buf.UnreadByte()
		d.parseDictByte(v)
	case 'i':
		d.buf.UnreadByte()
		d.parseIntByte(v)
	default:
		d.buf.UnreadByte()
		d.parseStringByte(v)
	}
}

func (d *Decode) parseDictByte(v reflect.Value) {
	dv, e := d.buf.ReadByte() // d
	for i := 0; ; i++ {       //add key
		if e != nil {
			fmt.Errorf("%v", e)
		}
		sl, e := d.buf.ReadBytes(byte(':')) // 4:
		if e != nil {
			fmt.Errorf("%v", e)
		}
		slv := sl[:len(sl)-1]                            // 4
		sliv, e := strconv.ParseInt(string(slv), 10, 64) // 4 as int64
		if e != nil {
			fmt.Errorf("%v", e)
		}
		buffer := make([]byte, sliv, sliv) // []byte{info}
		d.buf.Read(buffer)

		var r []byte
		if i == 0 {
			r = make([]byte, 0, 1+len(sl)+len(buffer))
			r = append(r, dv)
			r = append(r, sl...)
			r = append(r, buffer...)
		} else {
			r = make([]byte, 0, len(sl)+len(buffer))
			r = append(r, sl...)
			r = append(r, buffer...)
		}

		rv := reflect.AppendSlice(v, reflect.ValueOf(r))
		v.Set(rv)
		d.readNextByte(v)
		b, e := d.buf.ReadByte()
		if e != nil {
			fmt.Errorf("could not read byte: %v", e)
		}
		if b == 'e' {
			rv = reflect.Append(v, reflect.ValueOf(b))
			v.Set(rv)
			return
		}
		d.buf.UnreadByte()
	}
}

func (d *Decode) parseIntByte(v reflect.Value) {
	i, e := d.buf.ReadByte() // i
	if e != nil {
		fmt.Errorf("could not read byte: %v", e)
	}
	iv, e := d.buf.ReadBytes(byte('e')) // 3e
	if e != nil {
		fmt.Errorf("could not read bytes : %v", e)
	}

	r := make([]byte, 0, 1+len(iv))
	r = append(r, i)
	r = append(r, iv...)

	rv := reflect.AppendSlice(v, reflect.ValueOf(r))
	v.Set(rv)
}

func (d *Decode) parseStringByte(v reflect.Value) {
	sl, e := d.buf.ReadBytes(byte(':')) // 4:
	if e != nil {
		fmt.Errorf("%v", e)
	}
	slv := sl[:len(sl)-1]                            // 4
	sliv, e := strconv.ParseInt(string(slv), 10, 64) // 4 as int64
	if e != nil {
		fmt.Errorf("%v", e)
	}
	buffer := make([]byte, sliv, sliv) // []byte{info}
	d.buf.Read(buffer)

	r := make([]byte, 0, 1+len(sl)+len(buffer))
	r = append(r, sl...)
	r = append(r, buffer...)

	rv := reflect.AppendSlice(v, reflect.ValueOf(r))
	v.Set(rv)
}
