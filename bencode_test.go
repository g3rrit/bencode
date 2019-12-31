package bencode

import (
	"testing"
	"math/big"
)

func TestGet(t *testing.T) {
	cases := []struct {
		in Dictionary
		want Element
	}{
		{ Dictionary { Value: [][2]Element {
			[2]Element {ByteString{ Value: []byte("abe") }, ByteString { Value: []byte("ef") } },
			 } },
			 ByteString { Value: []byte("ef") } },

	}
	for _, c := range cases {
		got, err := c.in.Get("abe")
		if err != nil {
			t.Errorf("element not present in Dictionary (%q) -> (%q)", c.in, c.want)
			continue
		}
		if got.String() != c.want.String() {
			t.Errorf("(%q).String() == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestString(t *testing.T) {
	cases := []struct {
		in Element
		want string
	}{
		{ Integer { Value: big.NewInt(10) }, "10"},
		{ Integer { Value: big.NewInt(0) }, "0"},
		{ Integer { Value: big.NewInt(-10) }, "-10"},

		{ ByteString { Value: []byte("Hello") }, "Hello"},
		{ ByteString { Value: []byte{  } }, ""},

		{ List { Value: []Element{ Integer{ Value: big.NewInt(10) } }}, "{\n10\n}" },

		{ Dictionary { Value: [][2]Element { } }, "{\n}" },
	}
	for _, c := range cases {
		got := c.in.String()
		if got != c.want {
			t.Errorf("(%q).String() == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestEncode(t *testing.T) {
	cases := []struct {
		in Element
		want []byte
	}{
		{ Integer { Value: big.NewInt(10) }, []byte("i10e")},
		{ Integer { Value: big.NewInt(0) }, []byte("i0e")},
		{ Integer { Value: big.NewInt(-10) }, []byte("i-10e")},

		{ ByteString { Value: []byte("test")}, []byte("b4:teste")},
		{ ByteString { Value: []byte{  } }, []byte("b0:e")},

		{ List { Value: []Element{ Integer{ Value: big.NewInt(10) } }}, []byte("li10ee") },

		{ Dictionary { Value: [][2]Element{ } }, []byte("de") },
	}
	for _, c := range cases {
		got := c.in.Encode()
		if string(got) != string(c.want) {
			t.Errorf("(%q).Encode() == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestDecode(t *testing.T) {
	cases := []struct {
		in string
		want Element
	}{
		{ "i10e", Integer{ Value: big.NewInt(10) } },
		{ "i-10e", Integer{ Value: big.NewInt(-10) } },
		{ "i0e", Integer{ Value: big.NewInt(0) } },

		{ "0:", ByteString{ Value: []byte{  } } },
		{ "1:e", ByteString{ Value: []byte("e") } },
		{ "3:abe", ByteString{ Value: []byte("abe") } },

		{ "l3:abee", List { Value: []Element { ByteString{ Value: []byte("abe") } } } },
		{ "l3:abei43ee", List { Value: []Element { ByteString{ Value: []byte("abe") },
			Integer { Value: big.NewInt(43) } } } },
		{ "l3:abeli10eei43ee", List { Value: []Element { ByteString{ Value: []byte("abe") },
			List { Value: []Element { Integer { Value: big.NewInt(10) } } },
			Integer { Value: big.NewInt(43) } } } },

		{ "d3:abe2:efe", Dictionary { Value: [][2]Element {
			[2]Element {ByteString{ Value: []byte("abe") }, ByteString { Value: []byte("ef") } },
			 } } },
	}
	for _, c := range cases {
		got, e, err := Decode([]byte(c.in))
		if err != nil {
			t.Errorf("Error(%q -> %q): %q", c.in, c.want, err.Error())
		} else if got.String() != c.want.String() {
			t.Errorf("Decode([]byte(%q)) == %q, want %q - e: %q", c.in, got, c.want, e)
		}
	}
}
