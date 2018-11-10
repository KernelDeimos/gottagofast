package toolparse

/*
TODO: Eventually this file should be generated from a generic algorithm,
	  operating on []string instead of []D.Token, using an icga script against
	  the GottaGoFast package's implementation.
*/

import (
	"errors"
	"strings"

	"github.com/KernelDeimos/gottagofast/datastack"
)

const (
	ParseTokenOpen   = '('
	ParseTokenClose  = ')'
	ParseTokenQuote  = '\''
	ParseTokenEscape = '\\'
	ParseTokenSep    = ' '

	// Should be all symbols above, but excluding the secape character
	ParseSyntaxMask = ` ()'`
)

// ParseListSimple parses a list from a string with a few hard-coded syntax
// rules:
// - items are separated by spaces
// - strings containing spaces can be surrounded by single-quotes
// - The `\` symbol can escape single-quotes and other `\`s within strings
// - parethesis can be used to next lists inside other lists
func ParseListSimple(input string) ([]interface{}, error) {
	input = strings.TrimSpace(input)
	input += " "

	// PARSER PARAMETERS
	sep := ' '

	// PARSER STATE VARIABLES
	var buffer string
	var escape bool
	// -- ensure valid bracket count
	var bracketLevel int

	// PARSER STACK
	root := []interface{}{}
	stack := datastack.New()
	current := &root

	// Parser states
	const (
		StateFindNext = iota
		StateEatQuote
	)

	state := StateFindNext

	for _, thisRune := range input {
		if state == StateFindNext {
			isSyntax := strings.Contains(ParseSyntaxMask, string(thisRune))

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

			if thisRune == sep {
				continue
			}
			if thisRune == ParseTokenOpen /* ( */ {
				bracketLevel++
				newList := []interface{}{}
				stack.Push(current)
				current = &newList
				continue
			}
			if thisRune == ParseTokenClose /* ) */ {
				bracketLevel--
				if bracketLevel < 0 {
					goto ErrorCase1
				}
				doneList := current
				lastList := stack.Pop().(*[]interface{})
				current = lastList

				// Append current list as a literal list (non-pointer)
				*current = append(*current, *doneList)
				continue
			}
			if thisRune == ParseTokenQuote /* ' */ {
				if buffer == "" {
					state = StateEatQuote
				}
				continue
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
			if thisRune == ParseTokenQuote {
				*current = append(*current, buffer)
				buffer = ""
				state = StateFindNext
				continue
			}

			// Check for escape character
			if thisRune == ParseTokenEscape {
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
