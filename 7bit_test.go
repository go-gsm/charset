package charset

import (
	"log"
	"testing"
)

func TestEncode7Bit(t *testing.T) {
	data := []struct {
		str string
		exp []byte
	}{
		{"hello[world]! ы?", mustBytes("E8329BFDDEF0EE6F399BBCF18540BF1F")},
		{"AAAAAAAAAAAAAAB\r", mustBytes("C16030180C0683C16030180C0A1B0D")},
		{"AAAAAAAAAAAAAAB", mustBytes("C16030180C0683C16030180C0A1B")},
		{"height of eifel", mustBytes("E872FA8CA683DE6650396D2EB31B")},
	}
	for _, d := range data {
		if string(d.exp) != string(Pack7Bit(Encode7Bit(d.str))) {
			t.Errorf("Expected %v, got %v\n", d.exp, Encode7Bit(d.str))
		}
	}
}

func TestDecode7Bit(t *testing.T) {
	data := []struct {
		exp   string
		pack7 []byte
	}{
		// ы -> ?
		{"hello[world]! ??", mustBytes("E8329BFDDEF0EE6F399BBCF18540BF1F")},
		{"AAAAAAAAAAAAAAB\r", mustBytes("C16030180C0683C16030180C0A1B0D")},
		{"AAAAAAAAAAAAAAB", mustBytes("C16030180C0683C16030180C0A1B")},
		{"height of eifel", mustBytes("E872FA8CA683DE6650396D2EB31B")},
	}
	for _, d := range data {
		log.Println(displayPack(d.pack7))
		out, err := Decode7Bit(Unpack7Bit(d.pack7))
		if err != nil {
			t.Errorf("Expected nil, got %v\n", err)
		}
		if string(d.exp) != string(out) {
			t.Errorf("Expected %v, got %v\n", d.exp, out)
		}
	}
}

func TestPack7Bit(t *testing.T) {
	raw7 := []byte{Esc, 0x3c, Esc, 0x3e}
	exp := []byte{0x1b, 0xde, 0xc6, 0x7}
	if string(Pack7Bit(raw7)) != string(exp) {
		t.Errorf("Expected %v, got %v\n", exp, Pack7Bit(raw7))
	}
}

func TestUnpack7Bit(t *testing.T) {
	pack7 := []byte{0x1b, 0xde, 0xc6, 0x7}
	exp := []byte{Esc, 0x3c, Esc, 0x3e}
	if string(Unpack7Bit(pack7)) != string(exp) {
		t.Errorf("Expected %v, got %v\n", exp, Unpack7Bit(pack7))
	}
}

func TestIsGsmAlpha(t *testing.T) {
	data := []struct {
		expected bool
		str      string
	}{
		{true, `@£$¥èéùìòÇØøÅåΔ_ΦΓΛΩΠΨΣΘΞÆæßÉ!"#¤%&'()*+,-./0123456789:;<=>?¡ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÑÜ§¿abcdefghijklmnopqrstuvwxyzäöñüà|^€{}[]~\`},
		{false, "😁"},
		{false, "abcdedf💥"},
	}

	for _, d := range data {
		result := IsGsmAlpha(d.str)
		if result != d.expected {
			t.Errorf("Expected %v, got %v for %v\n", d.expected, result, d.str)
		}
	}
}

// mustBytes is an alias for parseOddHexStr, except that it will panic
// if there is any parse error.
func mustBytes(hex string) []byte {
	b, err := ParseOddHexStr(hex)
	if err != nil {
		panic(err)
	}
	return b
}
