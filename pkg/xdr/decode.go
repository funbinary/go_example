package xdr

import (
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
	"time"
)

var (
	errMaxSlice = "data exceeds max slice limit"
	errIODecode = "%s while decoding %d bytes"
)

func Read(r io.Reader, val interface{}) error {
	_, err := Unmarshal(r, val)
	return err
}

func ReadUint32(r io.Reader) (uint32, error) {
	var n uint32
	if err := Read(r, &n); err != nil {
		return n, err
	}

	return n, nil
}

func ReadOpaque(r io.Reader) ([]byte, error) {
	length, err := ReadUint32(r)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, length)
	if _, err = r.Read(buf); err != nil {
		return nil, err
	}

	return buf, nil
}

func Unmarshal(r io.Reader, v interface{}) (int, error) {
	d := Decoder{r: r}
	return d.Decode(v)
}

// A Decoder wraps an io.Reader that is expected to provide an XDR-encoded byte
// stream and provides several exposed methods to manually decode various XDR
// primitives without relying on reflection.  The NewDecoder function can be
// used to get a new Decoder directly.
//
// Typically, Unmarshal should be used instead of manual decoding.  A Decoder
// is exposed so it is possible to perform manual decoding should it be
// necessary in complex scenarios where automatic reflection-based decoding
// won't work.
type Decoder struct {
	r io.Reader

	// maxReadSize is the default maximum bytes an element can contain.  0
	// is unlimited and provides backwards compatability.  Setting it to a
	// non-zero value caps reads.
	maxReadSize uint
}

// Decode operates identically to the Unmarshal function with the exception of
// using the reader associated with the Decoder as the source of XDR-encoded
// data instead of a user-supplied reader.  See the Unmarhsal documentation for
// specifics.
func (d *Decoder) Decode(v interface{}) (int, error) {
	if v == nil {
		msg := "can't unmarshal to nil interface"
		return 0, unmarshalError("Unmarshal", ErrNilInterface, msg, nil,
			nil)
	}

	vv := reflect.ValueOf(v)
	if vv.Kind() != reflect.Ptr {
		msg := fmt.Sprintf("can't unmarshal to non-pointer '%v' - use "+
			"& operator", vv.Type().String())
		err := unmarshalError("Unmarshal", ErrBadArguments, msg, nil, nil)
		return 0, err
	}
	if vv.IsNil() && !vv.CanSet() {
		msg := fmt.Sprintf("can't unmarshal to unsettable '%v' - use "+
			"& operator", vv.Type().String())
		err := unmarshalError("Unmarshal", ErrNotSettable, msg, nil, nil)
		return 0, err
	}

	return d.decode(vv)
}

// decode is the main workhorse for unmarshalling via reflection.  It uses
// the passed reflection value to choose the XDR primitives to decode from
// the encapsulated reader.  It is a recursive function,
// so cyclic data structures are not supported and will result in an infinite
// loop.  It returns the  the number of bytes actually read.
func (d *Decoder) decode(v reflect.Value) (int, error) {
	if !v.IsValid() {
		msg := fmt.Sprintf("type '%s' is not valid", v.Kind().String())
		err := unmarshalError("decode", ErrUnsupportedType, msg, nil, nil)
		return 0, err
	}

	// Indirect through pointers allocating them as needed.
	ve, err := d.indirect(v)
	if err != nil {
		return 0, err
	}

	// Handle time.Time values by decoding them as an RFC3339 formatted
	// string with nanosecond precision.  Check the type string rather
	// than doing a full blown conversion to interface and type assertion
	// since checking a string is much quicker.
	if ve.Type().String() == "time.Time" {
		// Read the value as a string and parse it.
		timeString, n, err := d.DecodeString()
		if err != nil {
			return n, err
		}
		ttv, err := time.Parse(time.RFC3339, timeString)
		if err != nil {
			err := unmarshalError("decode", ErrParseTime,
				err.Error(), timeString, err)
			return n, err
		}
		ve.Set(reflect.ValueOf(ttv))
		return n, nil
	}

	// Handle native Go types.
	switch ve.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int:
		i, n, err := d.DecodeInt()
		if err != nil {
			return n, err
		}
		if ve.OverflowInt(int64(i)) {
			msg := fmt.Sprintf("signed integer too large to fit '%s'",
				ve.Kind().String())
			err = unmarshalError("decode", ErrOverflow, msg, i, nil)
			return n, err
		}
		ve.SetInt(int64(i))
		return n, nil

	case reflect.Int64:
		i, n, err := d.DecodeHyper()
		if err != nil {
			return n, err
		}
		ve.SetInt(i)
		return n, nil

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint:
		ui, n, err := d.DecodeUint()
		if err != nil {
			return n, err
		}
		if ve.OverflowUint(uint64(ui)) {
			msg := fmt.Sprintf("unsigned integer too large to fit '%s'",
				ve.Kind().String())
			err = unmarshalError("decode", ErrOverflow, msg, ui, nil)
			return n, err
		}
		ve.SetUint(uint64(ui))
		return n, nil

	case reflect.Uint64:
		ui, n, err := d.DecodeUhyper()
		if err != nil {
			return n, err
		}
		ve.SetUint(ui)
		return n, nil

	case reflect.Bool:
		b, n, err := d.DecodeBool()
		if err != nil {
			return n, err
		}
		ve.SetBool(b)
		return n, nil

	case reflect.Float32:
		f, n, err := d.DecodeFloat()
		if err != nil {
			return n, err
		}
		ve.SetFloat(float64(f))
		return n, nil

	case reflect.Float64:
		f, n, err := d.DecodeDouble()
		if err != nil {
			return n, err
		}
		ve.SetFloat(f)
		return n, nil

	case reflect.String:
		s, n, err := d.DecodeString()
		if err != nil {
			return n, err
		}
		ve.SetString(s)
		return n, nil

	case reflect.Array:
		n, err := d.decodeFixedArray(ve, false)
		if err != nil {
			return n, err
		}
		return n, nil

	case reflect.Slice:
		n, err := d.decodeArray(ve, false)
		if err != nil {
			return n, err
		}
		return n, nil

	case reflect.Struct:
		n, err := d.decodeStruct(ve)
		if err != nil {
			return n, err
		}
		return n, nil

	case reflect.Map:
		n, err := d.decodeMap(ve)
		if err != nil {
			return n, err
		}
		return n, nil

	case reflect.Interface:
		n, err := d.decodeInterface(ve)
		if err != nil {
			return n, err
		}
		return n, nil
	}

	// The only unhandled types left are unsupported.  At the time of this
	// writing the only remaining unsupported types that exist are
	// reflect.Uintptr and reflect.UnsafePointer.
	msg := fmt.Sprintf("unsupported Go type '%s'", ve.Kind().String())
	err = unmarshalError("decode", ErrUnsupportedType, msg, nil, nil)
	return 0, err
}

// indirect dereferences pointers allocating them as needed until it reaches
// a non-pointer.  This allows transparent decoding through arbitrary levels
// of indirection.
func (d *Decoder) indirect(v reflect.Value) (reflect.Value, error) {
	rv := v
	for rv.Kind() == reflect.Ptr {
		// Allocate pointer if needed.
		isNil := rv.IsNil()
		if isNil && !rv.CanSet() {
			msg := fmt.Sprintf("unable to allocate pointer for '%v'",
				rv.Type().String())
			err := unmarshalError("indirect", ErrNotSettable, msg,
				nil, nil)
			return rv, err
		}
		if isNil {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}
	return rv, nil
}

// DecodeString treats the next bytes as a variable length XDR encoded string
// and returns the result as a string along with the number of bytes actually
// read.  Character encoding is assumed to be UTF-8 and therefore ASCII
// compatible.  If the underlying character encoding is not compatibile with
// this assumption, the data can instead be read as variable-length opaque data
// (DecodeOpaque) and manually converted as needed.
//
// An UnmarshalError is returned if there are insufficient bytes remaining or
// the string data is larger than the max length of a Go slice.
//
// Reference:
// 	RFC Section 4.11 - String
// 	Unsigned integer length followed by bytes zero-padded to a multiple of
// 	four
func (d *Decoder) DecodeString() (string, int, error) {
	dataLen, n, err := d.DecodeUint()
	if err != nil {
		return "", n, err
	}
	if uint(dataLen) > uint(math.MaxInt32) ||
		(d.maxReadSize != 0 && uint(dataLen) > d.maxReadSize) {
		err = unmarshalError("DecodeString", ErrOverflow, errMaxSlice,
			dataLen, nil)
		return "", n, err
	}

	opaque, n2, err := d.DecodeFixedOpaque(int32(dataLen))
	n += n2
	if err != nil {
		return "", n, err
	}
	return string(opaque), n, nil
}

// DecodeHyper treats the next 8 bytes as an XDR encoded hyper value and
// returns the result as an int64  along with the number of bytes actually read.
//
// An UnmarshalError is returned if there are insufficient bytes remaining.
//
// Reference:
// 	RFC Section 4.5 - Hyper Integer
// 	64-bit big-endian signed integer in range [-9223372036854775808, 9223372036854775807]
func (d *Decoder) DecodeHyper() (int64, int, error) {
	var buf [8]byte
	n, err := io.ReadFull(d.r, buf[:])
	if err != nil {
		msg := fmt.Sprintf(errIODecode, err.Error(), 8)
		err := unmarshalError("DecodeHyper", ErrIO, msg, buf[:n], err)
		return 0, n, err
	}

	rv := int64(buf[7]) | int64(buf[6])<<8 |
		int64(buf[5])<<16 | int64(buf[4])<<24 |
		int64(buf[3])<<32 | int64(buf[2])<<40 |
		int64(buf[1])<<48 | int64(buf[0])<<56
	return rv, n, err
}

// DecodeUhyper treats the next 8  bytes as an XDR encoded unsigned hyper value
// and returns the result as a uint64  along with the number of bytes actually
// read.
//
// An UnmarshalError is returned if there are insufficient bytes remaining.
//
// Reference:
// 	RFC Section 4.5 - Unsigned Hyper Integer
// 	64-bit big-endian unsigned integer in range [0, 18446744073709551615]
func (d *Decoder) DecodeUhyper() (uint64, int, error) {
	var buf [8]byte
	n, err := io.ReadFull(d.r, buf[:])
	if err != nil {
		msg := fmt.Sprintf(errIODecode, err.Error(), 8)
		err := unmarshalError("DecodeUhyper", ErrIO, msg, buf[:n], err)
		return 0, n, err
	}

	rv := uint64(buf[7]) | uint64(buf[6])<<8 |
		uint64(buf[5])<<16 | uint64(buf[4])<<24 |
		uint64(buf[3])<<32 | uint64(buf[2])<<40 |
		uint64(buf[1])<<48 | uint64(buf[0])<<56
	return rv, n, nil
}

// DecodeInt treats the next 4 bytes as an XDR encoded integer and returns the
// result as an int32 along with the number of bytes actually read.
//
// An UnmarshalError is returned if there are insufficient bytes remaining.
//
// Reference:
// 	RFC Section 4.1 - Integer
// 	32-bit big-endian signed integer in range [-2147483648, 2147483647]
func (d *Decoder) DecodeInt() (int32, int, error) {
	var buf [4]byte
	n, err := io.ReadFull(d.r, buf[:])
	if err != nil {
		msg := fmt.Sprintf(errIODecode, err.Error(), 4)
		err := unmarshalError("DecodeInt", ErrIO, msg, buf[:n], err)
		return 0, n, err
	}

	rv := int32(buf[3]) | int32(buf[2])<<8 |
		int32(buf[1])<<16 | int32(buf[0])<<24
	return rv, n, nil
}

// DecodeUint treats the next 4 bytes as an XDR encoded unsigned integer and
// returns the result as a uint32 along with the number of bytes actually read.
//
// An UnmarshalError is returned if there are insufficient bytes remaining.
//
// Reference:
// 	RFC Section 4.2 - Unsigned Integer
// 	32-bit big-endian unsigned integer in range [0, 4294967295]
func (d *Decoder) DecodeUint() (uint32, int, error) {
	var buf [4]byte
	n, err := io.ReadFull(d.r, buf[:])
	if err != nil {
		msg := fmt.Sprintf(errIODecode, err.Error(), 4)
		err := unmarshalError("DecodeUint", ErrIO, msg, buf[:n], err)
		return 0, n, err
	}

	rv := uint32(buf[3]) | uint32(buf[2])<<8 |
		uint32(buf[1])<<16 | uint32(buf[0])<<24
	return rv, n, nil
}

// DecodeBool treats the next 4 bytes as an XDR encoded boolean value and
// returns the result as a bool along with the number of bytes actually read.
//
// An UnmarshalError is returned if there are insufficient bytes remaining or
// the parsed value is not a 0 or 1.
//
// Reference:
// 	RFC Section 4.4 - Boolean
// 	Represented as an XDR encoded enumeration where 0 is false and 1 is true
func (d *Decoder) DecodeBool() (bool, int, error) {
	val, n, err := d.DecodeInt()
	if err != nil {
		return false, n, err
	}
	switch val {
	case 0:
		return false, n, nil
	case 1:
		return true, n, nil
	}

	err = unmarshalError("DecodeBool", ErrBadEnumValue, "bool not 0 or 1",
		val, nil)
	return false, n, err
}

// DecodeFloat treats the next 4 bytes as an XDR encoded floating point and
// returns the result as a float32 along with the number of bytes actually read.
//
// An UnmarshalError is returned if there are insufficient bytes remaining.
//
// Reference:
// 	RFC Section 4.6 - Floating Point
// 	32-bit single-precision IEEE 754 floating point
func (d *Decoder) DecodeFloat() (float32, int, error) {
	var buf [4]byte
	n, err := io.ReadFull(d.r, buf[:])
	if err != nil {
		msg := fmt.Sprintf(errIODecode, err.Error(), 4)
		err := unmarshalError("DecodeFloat", ErrIO, msg, buf[:n], err)
		return 0, n, err
	}

	val := uint32(buf[3]) | uint32(buf[2])<<8 |
		uint32(buf[1])<<16 | uint32(buf[0])<<24
	return math.Float32frombits(val), n, nil
}

// DecodeDouble treats the next 8 bytes as an XDR encoded double-precision
// floating point and returns the result as a float64 along with the number of
// bytes actually read.
//
// An UnmarshalError is returned if there are insufficient bytes remaining.
//
// Reference:
// 	RFC Section 4.7 -  Double-Precision Floating Point
// 	64-bit double-precision IEEE 754 floating point
func (d *Decoder) DecodeDouble() (float64, int, error) {
	var buf [8]byte
	n, err := io.ReadFull(d.r, buf[:])
	if err != nil {
		msg := fmt.Sprintf(errIODecode, err.Error(), 8)
		err := unmarshalError("DecodeDouble", ErrIO, msg, buf[:n], err)
		return 0, n, err
	}

	val := uint64(buf[7]) | uint64(buf[6])<<8 |
		uint64(buf[5])<<16 | uint64(buf[4])<<24 |
		uint64(buf[3])<<32 | uint64(buf[2])<<40 |
		uint64(buf[1])<<48 | uint64(buf[0])<<56
	return math.Float64frombits(val), n, nil
}

// DecodeEnum treats the next 4 bytes as an XDR encoded enumeration value and
// returns the result as an int32 after verifying that the value is in the
// provided map of valid values.   It also returns the number of bytes actually
// read.
//
// An UnmarshalError is returned if there are insufficient bytes remaining or
// the parsed enumeration value is not one of the provided valid values.
//
// Reference:
// 	RFC Section 4.3 - Enumeration
// 	Represented as an XDR encoded signed integer
func (d *Decoder) DecodeEnum(validEnums map[int32]bool) (int32, int, error) {
	val, n, err := d.DecodeInt()
	if err != nil {
		return 0, n, err
	}

	if !validEnums[val] {
		err := unmarshalError("DecodeEnum", ErrBadEnumValue,
			"invalid enum", val, nil)
		return 0, n, err
	}
	return val, n, nil
}

// decodeFixedArray treats the next bytes as a series of XDR encoded elements
// of the same type as the array represented by the reflection value and decodes
// each element into the passed array.  The ignoreOpaque flag controls whether
// or not uint8 (byte) elements should be decoded individually or as a fixed
// sequence of opaque data.  It returns the  the number of bytes actually read.
//
// An UnmarshalError is returned if any issues are encountered while decoding
// the array elements.
//
// Reference:
// 	RFC Section 4.12 - Fixed-Length Array
// 	Individually XDR encoded array elements
func (d *Decoder) decodeFixedArray(v reflect.Value, ignoreOpaque bool) (int, error) {
	// Treat [#]byte (byte is alias for uint8) as opaque data unless
	// ignored.
	if !ignoreOpaque && v.Type().Elem().Kind() == reflect.Uint8 {
		data, n, err := d.DecodeFixedOpaque(int32(v.Len()))
		if err != nil {
			return n, err
		}
		reflect.Copy(v, reflect.ValueOf(data))
		return n, nil
	}

	// Decode each array element.
	var n int
	for i := 0; i < v.Len(); i++ {
		n2, err := d.decode(v.Index(i))
		n += n2
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

// decodeArray treats the next bytes as a variable length series of XDR encoded
// elements of the same type as the array represented by the reflection value.
// The number of elements is obtained by first decoding the unsigned integer
// element count.  Then each element is decoded into the passed array. The
// ignoreOpaque flag controls whether or not uint8 (byte) elements should be
// decoded individually or as a variable sequence of opaque data.  It returns
// the number of bytes actually read.
//
// An UnmarshalError is returned if any issues are encountered while decoding
// the array elements.
//
// Reference:
// 	RFC Section 4.13 - Variable-Length Array
// 	Unsigned integer length followed by individually XDR encoded array
// 	elements
func (d *Decoder) decodeArray(v reflect.Value, ignoreOpaque bool) (int, error) {
	dataLen, n, err := d.DecodeUint()
	if err != nil {
		return n, err
	}
	if uint(dataLen) > uint(math.MaxInt32) ||
		(d.maxReadSize != 0 && uint(dataLen) > d.maxReadSize) {
		err := unmarshalError("decodeArray", ErrOverflow, errMaxSlice,
			dataLen, nil)
		return n, err
	}

	// Allocate storage for the slice elements (the underlying array) if
	// existing slice does not have enough capacity.
	sliceLen := int(dataLen)
	if v.Cap() < sliceLen {
		v.Set(reflect.MakeSlice(v.Type(), sliceLen, sliceLen))
	}
	if v.Len() < sliceLen {
		v.SetLen(sliceLen)
	}

	// Treat []byte (byte is alias for uint8) as opaque data unless ignored.
	if !ignoreOpaque && v.Type().Elem().Kind() == reflect.Uint8 {
		data, n2, err := d.DecodeFixedOpaque(int32(sliceLen))
		n += n2
		if err != nil {
			return n, err
		}
		v.SetBytes(data)
		return n, nil
	}

	// Decode each slice element.
	for i := 0; i < sliceLen; i++ {
		n2, err := d.decode(v.Index(i))
		n += n2
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func (d *Decoder) decodeArrayByLength(v reflect.Value, dataLen int) (int, error) {
	n := 0
	// Allocate storage for the slice elements (the underlying array) if
	// existing slice does not have enough capacity.
	sliceLen := int(dataLen)
	if v.Cap() < sliceLen {
		v.Set(reflect.MakeSlice(v.Type(), sliceLen, sliceLen))
	}
	if v.Len() < sliceLen {
		v.SetLen(sliceLen)
	}

	// Treat []byte (byte is alias for uint8) as opaque data unless ignored.
	if v.Type().Elem().Kind() == reflect.Uint8 {
		data, n2, err := d.DecodeFixedOpaque(int32(sliceLen))
		n += n2
		if err != nil {
			return n, err
		}
		v.SetBytes(data)
		return n, nil
	}

	// Decode each slice element.
	for i := 0; i < sliceLen; i++ {
		n2, err := d.decode(v.Index(i))
		n += n2
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

// decodeStruct treats the next bytes as a series of XDR encoded elements
// of the same type as the exported fields of the struct represented by the
// passed reflection value.  Pointers are automatically indirected and
// allocated as necessary.  It returns the  the number of bytes actually read.
//
// An UnmarshalError is returned if any issues are encountered while decoding
// the elements.
//
// Reference:
// 	RFC Section 4.14 - Structure
// 	XDR encoded elements in the order of their declaration in the struct
func (d *Decoder) decodeStruct(v reflect.Value) (int, error) {
	var n int
	var union string
	vt := v.Type()

	for i := 0; i < v.NumField(); i++ {
		// Skip unexported fields.
		vtf := vt.Field(i)
		if vtf.PkgPath != "" {
			continue
		}

		vf := v.Field(i)
		tag := parseTag(vtf.Tag)

		// RFC Section 4.19 - Optional data
		if tag.Get("optional") == "true" {
			if vf.Type().Kind() != reflect.Ptr {
				msg := fmt.Sprintf("optional must be a pointer, not '%v'",
					vf.Type().String())
				err := unmarshalError("decodeStruct", ErrBadOptional,
					msg, nil, nil)
				return n, err
			}

			hasopt, n2, err := d.DecodeBool()
			n += n2
			if err != nil {
				return n, err
			}
			if !hasopt {
				continue
			}
		}

		// Indirect through pointers allocating them as needed and
		// ensure the field is settable.
		vf, err := d.indirect(vf)
		if err != nil {
			return n, err
		}

		if !vf.CanSet() {
			msg := fmt.Sprintf("can't decode to unsettable '%v'",
				vf.Type().String())
			err := unmarshalError("decodeStruct", ErrNotSettable,
				msg, nil, nil)
			return n, err
		}
		limit := tag.Get("limit")
		if limit != "" {
			var l int
			l, err = strconv.Atoi(limit)
			if err == nil {
				switch vf.Kind() {
				case reflect.Slice:
					n2, err := d.decodeArrayByLength(vf, l)
					n += n2
					if err != nil {
						return n, err
					}
					continue

				}
			}
		}

		// Handle non-opaque data to []uint8 and [#]uint8 based on
		// struct tag.
		if tag.Get("opaque") == "false" {
			switch vf.Kind() {
			case reflect.Slice:
				n2, err := d.decodeArray(vf, true)
				n += n2
				if err != nil {
					return n, err
				}
				continue

			case reflect.Array:
				n2, err := d.decodeFixedArray(vf, true)
				n += n2
				if err != nil {
					return n, err
				}
				continue
			}
		}

		if union != "" {
			ucase := tag.Get("unioncase")
			if ucase != "" && ucase != union {
				continue
			}
		}

		// Decode each struct field.
		n2, err := d.decode(vf)
		n += n2
		if err != nil {
			return n, err
		}

		if tag.Get("union") == "true" {
			if vf.Type().ConvertibleTo(reflect.TypeOf(0)) {
				union = strconv.Itoa(int(vf.Convert(reflect.TypeOf(0)).Int()))
			} else if vf.Kind() == reflect.Bool {
				if vf.Bool() {
					union = "1"
				} else {
					union = "0"
				}
			} else {
				msg := fmt.Sprintf("type '%s' is not valid", vf.Kind().String())
				return n, unmarshalError("decodeStruct", ErrBadDiscriminant, msg, nil, nil)
			}
		}

	}

	return n, nil
}

// RFC Section 4.16 - Void
// RFC Section 4.17 - Constant
// RFC Section 4.18 - Typedef
// RFC Section 4.19 - Optional data
// RFC Sections 4.15 though 4.19 only apply to the data specification language
// which is not implemented by this package.

// decodeMap treats the next bytes as an XDR encoded variable array of 2-element
// structures whose fields are of the same type as the map keys and elements
// represented by the passed reflection value.  Pointers are automatically
// indirected and allocated as necessary.  It returns the  the number of bytes
// actually read.
//
// An UnmarshalError is returned if any issues are encountered while decoding
// the elements.
func (d *Decoder) decodeMap(v reflect.Value) (int, error) {
	dataLen, n, err := d.DecodeUint()
	if err != nil {
		return n, err
	}

	// Allocate storage for the underlying map if needed.
	vt := v.Type()
	if v.IsNil() {
		v.Set(reflect.MakeMap(vt))
	}

	// Decode each key and value according to their type.
	keyType := vt.Key()
	elemType := vt.Elem()
	for i := uint32(0); i < dataLen; i++ {
		key := reflect.New(keyType).Elem()
		n2, err := d.decode(key)
		n += n2
		if err != nil {
			return n, err
		}

		val := reflect.New(elemType).Elem()
		n2, err = d.decode(val)
		n += n2
		if err != nil {
			return n, err
		}
		v.SetMapIndex(key, val)
	}
	return n, nil
}

// RFC Section 4.8 -  Quadruple-Precision Floating Point
// 128-bit quadruple-precision floating point
// Not Implemented

// DecodeFixedOpaque treats the next 'size' bytes as XDR encoded opaque data and
// returns the result as a byte slice along with the number of bytes actually
// read.
//
// An UnmarshalError is returned if there are insufficient bytes remaining to
// satisfy the passed size, including the necessary padding to make it a
// multiple of 4.
//
// Reference:
// 	RFC Section 4.9 - Fixed-Length Opaque Data
// 	Fixed-length uninterpreted data zero-padded to a multiple of four
func (d *Decoder) DecodeFixedOpaque(size int32) ([]byte, int, error) {
	// Nothing to do if size is 0.
	if size == 0 {
		return nil, 0, nil
	}

	pad := (4 - (size % 4)) % 4
	paddedSize := size + pad
	if uint(paddedSize) > uint(math.MaxInt32) {
		err := unmarshalError("DecodeFixedOpaque", ErrOverflow,
			errMaxSlice, paddedSize, nil)
		return nil, 0, err
	}

	buf := make([]byte, paddedSize)
	n, err := io.ReadFull(d.r, buf)
	if err != nil {
		msg := fmt.Sprintf(errIODecode, err.Error(), paddedSize)
		err := unmarshalError("DecodeFixedOpaque", ErrIO, msg, buf[:n],
			err)
		return nil, n, err
	}
	return buf[0:size], n, nil
}

// decodeInterface examines the interface represented by the passed reflection
// value to detect whether it is an interface that can be decoded into and
// if it is, extracts the underlying value to pass back into the decode function
// for decoding according to its type.  It returns the  the number of bytes
// actually read.
//
// An UnmarshalError is returned if any issues are encountered while decoding
// the interface.
func (d *Decoder) decodeInterface(v reflect.Value) (int, error) {
	if v.IsNil() || !v.CanInterface() {
		msg := fmt.Sprintf("can't decode to nil interface")
		err := unmarshalError("decodeInterface", ErrNilInterface, msg,
			nil, nil)
		return 0, err
	}

	// Extract underlying value from the interface and indirect through
	// pointers allocating them as needed.
	ve := reflect.ValueOf(v.Interface())
	ve, err := d.indirect(ve)
	if err != nil {
		return 0, err
	}
	if !ve.CanSet() {
		msg := fmt.Sprintf("can't decode to unsettable '%v'",
			ve.Type().String())
		err := unmarshalError("decodeInterface", ErrNotSettable, msg,
			nil, nil)
		return 0, err
	}
	return d.decode(ve)
}
