package jsonb

import (
	"golang.org/x/xerrors"
)

var whiteSpaceMap = map[byte]bool{
	' ':  true,
	'\n': true,
	'\t': true,
	'\r': true,
}

var valueMap = make([]jsonType, 256)

func init() {
	valueMap['"'] = jsonString
	valueMap['-'] = jsonNumber
	valueMap['0'] = jsonNumber
	valueMap['1'] = jsonNumber
	valueMap['2'] = jsonNumber
	valueMap['3'] = jsonNumber
	valueMap['4'] = jsonNumber
	valueMap['5'] = jsonNumber
	valueMap['6'] = jsonNumber
	valueMap['7'] = jsonNumber
	valueMap['8'] = jsonNumber
	valueMap['9'] = jsonNumber
	valueMap['t'] = jsonBoolean
	valueMap['f'] = jsonBoolean
	valueMap['n'] = jsonNull
	valueMap['['] = jsonArray
	valueMap['{'] = jsonObject
	valueMap[' '] = jsonWhitespace
	valueMap['\r'] = jsonWhitespace
	valueMap['\t'] = jsonWhitespace
	valueMap['\n'] = jsonWhitespace
}

var emptyJSON = []byte(`null`)

// Reader :
type Reader struct {
	typ   jsonType
	b     []byte
	pos   int
	len   int
	start int
	end   int
}

// NewReader :
func NewReader(b []byte) *Reader {
	return &Reader{b: b, len: len(b)}
}

// Bytes :
func (r *Reader) Bytes() []byte {
	return r.b
}

// ReadNext :
func (r *Reader) nextToken() byte {
	var c byte
	for i := r.pos; i < r.len; i++ {
		c = r.b[i]
		if _, isOk := whiteSpaceMap[c]; isOk {
			r.b = append(r.b[:i], r.b[i+1:]...)
			r.len = r.len - 1
			i--
			continue
		}
		r.pos = i + 1
		return c
	}
	return 0
}

func (r *Reader) prevToken() byte {
	if r.pos > 0 {
		return r.b[r.pos-1]
	}
	return 0
}

func (r *Reader) peekType() jsonType {
	c := r.nextToken()
	defer r.unreadByte()
	typ := valueMap[c]
	return typ
}

// GetBytes :
func (r *Reader) GetBytes() (b []byte) {
	r.start = r.pos
	c := r.nextToken()
	switch c {
	case '"':
		// r.skipString()
	case 'n':
		// r.skipThreeBytes('u', 'l', 'l') // null
	case 't':
		// iter.skipThreeBytes('r', 'u', 'e') // true
	case 'f':
		// iter.skipFourBytes('a', 'l', 's', 'e') // false
	case '0':
		// iter.unreadByte()
		// iter.ReadFloat32()
	case '-', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// iter.skipNumber()
	case '[':
		// r.skipArray()
	case '{':
		// r.skipObject()
	default:
		// iter.ReportError("Skip", fmt.Sprintf("do not know how to skip: %v", c))
		return
	}
	return
}

func (r *Reader) skipArray() {
	level := 1
	c := r.nextToken()
	if c != '[' {
		return
	}

	for i := r.pos; i < r.len; i++ {
		switch r.b[i] {
		case '"': // If inside string, skip it
			// iter.head = i + 1
			// iter.skipString()
			// i = iter.head - 1 // it will be i++ soon
		case '[': // If open symbol, increase level
			level++
		case ']': // If close symbol, increase level
			level--

			// If we have returned to the original level, we're done
			if level <= 0 {
				r.pos = i + 1
				return
			}
		}
	}
}

// ReadBytes :
func (r *Reader) ReadBytes() ([]byte, error) {
	i := r.pos
	r.skip()
	return r.b[i:r.pos], nil
}

// ReadValue :
func (r *Reader) ReadValue() (interface{}, error) {
	typ := r.peekType()
	switch typ {
	case jsonString:
		return r.ReadString()
	case jsonNumber:
		return r.ReadNumber()
	case jsonBoolean:
		return r.ReadBoolean()
	case jsonNull:
		return r.ReadNull(), nil
	case jsonArray:
		var v []interface{}
		if err := r.ReadArray(func(it *Reader) error {
			x, err := it.ReadValue()
			if err != nil {
				return err
			}
			v = append(v, x)
			return nil
		}); err != nil {
			return v, err
		}
		return v, nil
	case jsonObject:
		var v map[string]interface{}
		if err := r.ReadObject(func(it *Reader, k string) error {
			if v == nil {
				v = make(map[string]interface{})
			}
			x, err := it.ReadValue()
			if err != nil {
				return err
			}
			v[k] = x
			return nil
		}); err != nil {
			return nil, err
		}
		return v, nil
	default:
		return nil, xerrors.New("invalid json format")
	}
}

func (r *Reader) unreadByte() *Reader {
	if r.pos > 0 {
		r.pos--
	}
	return r
}

// ReadBoolean :
func (r *Reader) ReadBoolean() (bool, error) {
	c := r.nextToken()
	if c == 't' {
		r.skipBytes([]byte{'r', 'u', 'e'})
		return true, nil
	}
	if c == 'f' {
		r.skipBytes([]byte{'a', 'l', 's', 'e'})
		return false, nil
	}
	return false, xerrors.New("invalid boolean value")
}

// ReadNull :
func (r *Reader) ReadNull() (b bool) {
	c := r.nextToken()
	if c == 'n' {
		r.skipBytes([]byte{'u', 'l', 'l'})
		return true
	}
	return
}
