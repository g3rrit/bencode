package bencode

import (
	"testing"
)

type Case struct {
	elem Element
	enc  string
	str  string
}

var cases []Case = []Case {
	{ Integer(10), "i10e", "10"},
	{ Integer(0), "i0e", "0"},
	{ Integer(-10), "i-10e", "-10"},

	{ ByteString([]byte("Hello")), "5:Hello", "Hello"},
	{ ByteString([]byte{  }), "0:", ""},

	{ List([]Element{ Integer(10) }), "li10ee", "[\n10\n]" },

	{ Dictionary(map[string]Element{}), "de", "{\n}" },
}

func TestString(t *testing.T) {
	for _, c := range cases {
		got := c.elem.String()
		if got != c.str {
			t.Errorf("(%q).String() == %q, want %q", c.elem, got, c.str)
		}
	}
}

func TestEncode(t *testing.T) {
	for _, c := range cases {
		got := c.elem.Encode()
		if string(got) != c.enc {
			t.Errorf("(%q).Encode() == %q, want %q", c.elem, got, c.enc)
		}
	}
}

func TestDecode(t *testing.T) {
	for _, c := range cases {
		got, err := D([]byte(c.enc))
		if err != nil {
			t.Errorf("Error(%q -> %q): %q", c.enc, c.elem, err.Error())
		} else if got.String() != c.str {
			t.Errorf("Decode([]byte(%q)) == %q, want %q", c.enc, got, c.str)
		}
	}
}
