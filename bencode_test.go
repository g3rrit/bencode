package bencode

import (
	"testing"
	"math/big"
)

func TestString(t *testing.T) {
	cases := []struct {
		in Element
		want string
	}{
		{ Integer { val: big.NewInt(10) }, "10"},
		{ Integer { val: big.NewInt(0) }, "0"},
		{ Integer { val: big.NewInt(-10) }, "-10"},

		{ ByteString { val: []byte{ 0x0f, 0xab, 0xfd } }, "\\X0fabfd"},
		{ ByteString { val: []byte{  } }, "\\X"},

		{ List { val: []Element{ Integer{ val: big.NewInt(10) } }}, "{\n10\n}" },

		{ Dictionary { val: [][2]Element { } }, "{\n}" },
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
		{ Integer { val: big.NewInt(10) }, []byte("i10e")},
		{ Integer { val: big.NewInt(0) }, []byte("i0e")},
		{ Integer { val: big.NewInt(-10) }, []byte("i-10e")},

		{ ByteString { val: []byte("test")}, []byte("b4:teste")},
		{ ByteString { val: []byte{  } }, []byte("b0:e")},

		{ List { val: []Element{ Integer{ val: big.NewInt(10) } }}, []byte("li10ee") },

		{ Dictionary { val: [][2]Element{ } }, []byte("de") },
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
		{ "i10e", Integer{ val: big.NewInt(10) } },
		{ "i-10e", Integer{ val: big.NewInt(-10) } },
		{ "i0e", Integer{ val: big.NewInt(0) } },

		{ "0:", ByteString{ val: []byte{  } } },
		{ "1:e", ByteString{ val: []byte("e") } },
		{ "3:abe", ByteString{ val: []byte("abe") } },

		{ "l3:abee", List { val: []Element { ByteString{ val: []byte("abe") } } } },
		{ "l3:abei43ee", List { val: []Element { ByteString{ val: []byte("abe") },
			Integer { val: big.NewInt(43) } } } },
		{ "l3:abeli10eei43ee", List { val: []Element { ByteString{ val: []byte("abe") },
			List { val: []Element { Integer { val: big.NewInt(10) } } },
			Integer { val: big.NewInt(43) } } } },

		{ "d3:abe2:efe", Dictionary { val: [][2]Element {
			[2]Element {ByteString{ val: []byte("abe") }, ByteString { val: []byte("ef") } },
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
