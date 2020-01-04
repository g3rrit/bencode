package bencode

import (
	"strconv"
	"errors"
)

//--------------------------------------------------------
// ELEMENT INTERFACE
//--------------------------------------------------------
type Element interface {
	isElement()
	String() string
	Encode() []byte
}

//--------------------------------------------------------
// ELEMENT TYPES
//--------------------------------------------------------
type Integer int

type ByteString []byte

type List []Element

type Dictionary map[string]Element

//--------------------------------------------------------
// INTERFACE FUNCTIONS
//--------------------------------------------------------

func (_ Integer) isElement() {}
func (_ ByteString) isElement() {}
func (_ List) isElement() {}
func (_ Dictionary) isElement() {}

func (elem Integer) String() string {
	return strconv.Itoa(int(elem))
}

func (elem ByteString) String() string {
	res := ""
	for _, val := range elem {
		if val >= 32 && val <= 126 {
			res += string(val)
		} else {
			res += string(".")
		}
	}
	return res
}

func (elem List) String() string {
	res := "[\n"
	for i := range elem {
		res += elem[i].String() + "\n"
	}
	res += "]"
	return res
}

func (elem Dictionary) String() string {
	res := "{\n"
	for k, v := range elem {
		res += k + " => " + v.String() + "\n"
	}
	res += "}"
	return res
}

func (elem Integer) Encode() []byte {
	res := append([]byte("i"), []byte(elem.String())...)
	return append(res, []byte("e")...)
}

func (elem ByteString) Encode() []byte {
	res := []byte(strconv.Itoa(len(elem)))
	res = append(res, []byte(":")...)
	return append(res, elem...)
}

func (elem List) Encode() []byte {
	res := []byte("l")
	for i := range elem {
		res = append(res, elem[i].Encode()...)
	}
	return append(res, []byte("e")...)
}

func (elem Dictionary) Encode() []byte {
	res := []byte("d")
	for k, v := range elem {
		res = append(res, (ByteString([]byte(k))).Encode()...)
		res = append(res, v.Encode()...)
	}
	return append(res, []byte("e")...)
}

func D(data []byte) (Element, error) {
	res, _, err := Decode(data)
	if err != nil {
		return res, err
	}
	return res, nil
}

func Decode(data []byte) (Element, int, error) {
	if len(data) <= 1 {
		return nil, 0, errors.New("missing data")
	}
	switch data[0] {
	case 'i':
		return decodeInteger(data)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return decodeByteString(data)
	case 'l':
		return decodeList(data)
	case 'd':
		return decodeDictionary(data)
	default:
		return nil, 0, errors.New("invalid format")
	}
}

func getEnd(data []byte) (int, error) {
	for i := 0; i < len(data); i++ {
		if data[i] == 'e' {
			return i, nil
		}
	}
	return 0, errors.New("missing end")
}

func decodeInteger(data []byte) (Integer, int, error) {
	var res Integer = 0
	e, err := getEnd(data)
	if err != nil {
		return res, 0, err
	}
	i, err := strconv.Atoi(string(data[1:e]))
	if err != nil {
		return res, 0, err
	}
	res = Integer(i)
	return res, e, nil
}

func decodeByteString(data []byte) (ByteString, int, error) {
	dataSize := len(data)
	if dataSize <= 1 {
		return nil,0, errors.New("bytestring does not contain length")
	}
	var i int
	for i = 1; data[i] != ':'; i++ {
		if i >= dataSize - 1 {
			return nil, 0, errors.New("bytestring does not contain length")
		}
	}
	l, err := strconv.Atoi(string(data[:i]))

	if err != nil {
		return nil, 0, errors.New("unable to parse length of bytestring")
	}
	if dataSize <= i + l {
		return nil, 0, errors.New("missing data")
	}

	var res ByteString = ByteString(make([]byte, i))
	if l == 0 {
		res = []byte{}
	} else {
		res = data[i + 1:i + l + 1]
	}
	return res, i + l, nil
}

func decodeList(data []byte) (List, int, error) {
	var res List = List(make([]Element, 0))
	var i int
	for i = 1; data[i] != 'e'; {
		elem, e, err := Decode(data[i:])
		if err != nil {
			return res, 0, err
		}
		res = append(res, elem)
		i += e + 1
		if i >= len(data) {
			return res, 0, errors.New("missing end of list")
		}
	}
	return res, i, nil
}

func decodeDictionary(data []byte) (Dictionary, int, error) {
	var res Dictionary = Dictionary(make(map[string] Element))
	var i int
	for i = 1; data[i] != 'e'; {
		k, e, err := Decode(data[i:])
		if err != nil {
			return res, 0, err
		}
		i += e + 1
		if i >= len(data) {
			return res, 0, errors.New("missing key in Dictionary")
		}
		v, e, err := Decode(data[i:])
		if err != nil {
			return res, 0, err
		}

		switch nt := k.(type) {
		case ByteString:
			res[string(nt)] = v
		default:
			return res, 0, errors.New("keys in Dictionary need to be ByteStrings")
		}

		i += e + 1
		if i >= len(data) {
			return res, 0, errors.New("missing end of Dictionary")
		}
	}
	return res, i, nil
}
