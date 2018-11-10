package listsimpleopts

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestParseListSimple(t *testing.T) {
	o := NewOptions(' ')
	o.AddBracket(Bracket{'(', ')', ""})
	o.AddBracket(Bracket{'[', ']', "vector"})
	// Test with valid input
	{
		out, err := ParseListSimpleOptions(*o, "a b (c d (e)) [f g(h)i] j (k)")
		if err != nil {
			t.Error(err)
		}
		outBytes, err := json.Marshal(out)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(string(outBytes))
	}
}
