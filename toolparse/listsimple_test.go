package toolparse

import (
	"encoding/json"
	"testing"
)

func TestParseListSimple(t *testing.T) {
	// Test with valid input
	{
		out, err := ParseListSimple("a b (c d (e)) f g(h)i j (k)")
		if err != nil {
			t.Error(err)
		}
		outBytes, err := json.Marshal(out)
		if err != nil {
			t.Error(err)
		}
		if string(outBytes) !=
			`["a","b",["c","d",["e"]],"f","g",["h"],"i","j",["k"]]` {
			t.Fail()
		}
	}

	// Test with invalid input
	{
		_, err := ParseListSimple("a b (c d (e) f g(h)i j (k)")
		if err == nil {
			t.Fail()
		}
	}
	{
		_, err := ParseListSimple("a b (c d (e))) f g(h)i j (k)")
		if err == nil {
			t.Fail()
		}
	}

}
