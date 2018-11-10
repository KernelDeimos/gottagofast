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

func TestParseListSimpleReportsCorrectTypes(t *testing.T) {
	// Test with valid input
	{
		out, err := ParseListSimple("a b (c d (e)) f g(h)i j (k)")
		if err != nil {
			t.Error(err)
		}

		check := func(ok bool, index int) {
			if !ok {
				t.Errorf("item at index %d has the wrong type", index)
			}
		}

		var ok bool
		_, ok = out[0].(string)
		check(ok, 0)
		_, ok = out[1].(string)
		check(ok, 1)
		_, ok = out[2].([]interface{})
		check(ok, 2)
		_, ok = out[8].([]interface{})
		check(ok, 3)
	}

}
