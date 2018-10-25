package utilstr

// AtomicReplace replaces multiple substrings in a string atomically. Here,
// atomically means it will replace all the substrings without replacing
// instances which occur after replacing other substrings.
func AtomicReplace(
	input string, replacements map[string]string,
) (output string, count map[string]int) {
	count = map[string]int{}
	output = ""

	// Map input string index to substring being replaced
	replacementIndexes := map[int]string{}

	// Iterate over all the strings being replaced and populate the map of
	// indexes that need replacements
	for key := range replacements {
		indices := IndexAll(input, key)
		for _, index := range indices {
			replacementIndexes[index] = key
		}
		// Ensure this key is added to the map of counts
		count[key] = 0
	}

	// Take a sip of coffee
	// (this didn't require writing any code)

	// Idea: using a sorted list of replacement keys here may allow for a more
	// efficient replacement algorithm, depending on how Golang is optimized.

	// Perform replacements
	linput := len(input)
	for i := 0; i < linput; i++ {
		if _, exists := replacementIndexes[i]; exists {
			// Get the original and replacement substrings for this index
			replaceFrom := replacementIndexes[i]
			replaceTo := replacements[replaceFrom]
			// Add replacement to output and scan past original substring
			output += replaceTo
			i += len(replaceFrom) - 1 // (-1 because of loop increment)
			// Update replacement count
			count[replaceFrom]++
		} else {
			// Add the character at this index to the output string
			output += input[i : i+1]
		}
	}

	return // NB: named return values
}
