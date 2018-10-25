package utilstr

import (
	"strings"
)

// IndexAll reports the index of every occurrance of a substring. When two
// substrings overlap, the index of the first instance will be taken.
func IndexAll(
	input, substr string,
) []int {
	indices := []int{}
	var startAt int
	// While an instance still exists after the start location
	for {
		// Get index of substring if one exists
		index := strings.Index(input[startAt:], substr)
		if index == -1 {
			break
		}

		// Add index as startAt+index (index on input string)
		indices = append(indices, startAt+index)
		// Increement startAt to be the position immediately after this instance
		// of the substring
		startAt += index + len(substr)
	}

	return indices
}
