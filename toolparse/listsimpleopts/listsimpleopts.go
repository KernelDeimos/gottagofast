package listsimpleopts

/*
TODO: Eventually this file should be generated from a generic algorithm,
	  operating on []string instead of []D.Token, using an icga script against
	  the GottaGoFast package's implementation.
*/

import (
	"errors"
	"fmt"
	"strings"

	"github.com/KernelDeimos/gottagofast/datastack"
)

type Options struct {
	Sep      rune
	Brackets map[rune]Bracket
	Quotes   map[rune]Quote
}

func NewOptions(sep rune) *Options {
	return &Options{
		Sep:      sep,
		Brackets: map[rune]Bracket{},
		Quotes:   map[rune]Quote{},
	}
}

// Bracket represents a set of characters which contain lists
type Bracket struct {
	Open    rune
	Close   rune
	Keyword string
}

// Quote represents a set of characters which contain scalar values
type Quote struct {
	Open   rune
	Close  rune
	Escape rune
}

func (opts *Options) AddQuote(quote Quote) {
	opts.Quotes[quote.Open] = quote
}

func (opts *Options) AddBracket(bracket Bracket) {
	opts.Brackets[bracket.Open] = bracket
}

func (opts Options) GetSyntaxMask() string {
	mask := map[rune]struct{}{}
	for _, bracket := range opts.Brackets {
		mask[bracket.Open] = struct{}{}
		mask[bracket.Close] = struct{}{}
	}
	for _, quote := range opts.Quotes {
		mask[quote.Open] = struct{}{}
		mask[quote.Close] = struct{}{}
	}
	maskStr := ""
	for run := range mask {
		maskStr += string(run)
	}

	return maskStr + string(opts.Sep)
}

func ParseListSimpleOptions(opts Options, input string) ([]interface{}, error) {
	input = strings.TrimSpace(input)
	input += " "

	syntaxMask := opts.GetSyntaxMask()

	// PARSER STATE VARIABLES
	var buffer string
	var escape bool
	// -- ensure valid bracket count
	var bracketLevel int
	var currentCloser *rune
	var stringCloser rune
	var stringEscaper rune

	// PARSER STACK
	root := []interface{}{}
	stack := datastack.New()
	current := &root

	// BRACKETS STACK (parallel stack)
	brackets := datastack.New()

	// Parser states
	const (
		StateFindNext = iota
		StateEatQuote
	)

	state := StateFindNext

	goInside := func(closer rune) {
		bracketLevel++
		newList := []interface{}{}
		stack.Push(current)
		current = &newList

		// set closer
		fmt.Println("--setcloser ", string(closer))
		brackets.Push(closer)
		currentCloser = &closer
	}

	goOutside := func() bool /*hasError*/ {
		bracketLevel--
		if bracketLevel < 0 {
			return true
		}
		doneList := current
		lastList := stack.Pop().(*[]interface{})
		current = lastList

		// Append current list as a literal list (non-pointer)
		*current = append(*current, *doneList)

		// set closer
		if bracketLevel > 0 {
			brackets.Pop()
			closer := brackets.Peek().(rune)
			fmt.Println("--retcloser ", string(closer))
			currentCloser = &closer
		} else {
			currentCloser = nil
		}

		return false
	}

	for _, thisRune := range input {
		fmt.Println(">", string(thisRune))
		if state == StateFindNext {
			isSyntax := strings.Contains(syntaxMask, string(thisRune))

			// Save buffer value if syntax character was found, otherwise
			// append the next character and continue iterating.
			if isSyntax {
				if buffer != "" {
					*current = append(*current, buffer)
					buffer = ""
				}
			} else {
				buffer += string(thisRune)
				continue
			}

			if thisRune == opts.Sep {
				continue
			}

			// Search for open brackets
			for test, bracket := range opts.Brackets {
				if test == thisRune {
					goInside(bracket.Close)
					if bracket.Keyword != "" {
						*current = append(*current, bracket.Keyword)
					}
					continue
				}
			}

			// Search for close brackets
			if currentCloser != nil {
				fmt.Println("??", string(*currentCloser), ",", string(thisRune))
				if *currentCloser == thisRune {
					fmt.Println("::", string(*currentCloser))
					if goOutside() {
						goto ErrorCase1
					}
					continue
				}
			}

			// Search for quotes
			for test, quote := range opts.Quotes {
				if test == thisRune && buffer == "" {
					state = StateEatQuote
					stringCloser = quote.Close
					stringEscaper = quote.Escape
					continue
				}
			}
		}
		if state == StateEatQuote {
			// If in escape mode, add this rune unconditionally,
			// then exit escape mode
			if escape {
				buffer += string(thisRune)
				escape = false
				continue
			}

			// Check for end quote
			if thisRune == stringCloser {
				*current = append(*current, buffer)
				buffer = ""
				state = StateFindNext
				continue
			}

			// Check for escape character
			if thisRune == stringEscaper {
				escape = true
				continue
			}

			buffer += string(thisRune)
		}
	}

	if bracketLevel > 0 {
		goto ErrorCase2
	}

	goto SuccessCase

ErrorCase1:
	return root, errors.New("unexpected close")

ErrorCase2:
	return root, errors.New("unexpected open")

SuccessCase:
	return root, nil

}
