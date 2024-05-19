package code49

import (
	"slices"
	"strings"
)

// Decode decodes a bunch of numbers `bsLines` which transferred from bar/space data, returns raw text.
// It treads the Code49 as Mode-0, ignores other non-data characters.
// If `skipChecksum` presents to true, Decode will ignore the last line of `bsLines`.
func Decode(bsLines []string, skipChecksum bool) string {
	var (
		buffer strings.Builder
		//index  int
		chars = []rune(charTable)
	)

	// For each line expect last line, as ignoring checksum
	bound := len(bsLines)
	if skipChecksum {
		bound = len(bsLines) - 1
	}
	for row := 0; row < bound; row++ {
		// Extract pattern, remove head '11' and tail '4'
		rowValue := bsLines[row]
		payload := rowValue[2 : len(rowValue)-1]

		// For each symbol
		for col := 0; col < 4; col++ {
			val := 0
			parity := rowParity[row][col]
			if parity == 'E' || len(bsLines) == row+1 {
				val = slices.Index(evenParityPatterns, payload[col<<3:(col+1)<<3])
			} else {
				val = slices.Index(oddParityPatterns, payload[col<<3:(col+1)<<3])
			}

			// Extract 2 characters from val
			c1, c2 := val/49, val%49

			// Find and append buffer
			suffix, prefix := chars[c1], chars[c2]
			if slices.Contains(nonDataMethodChart, prefix) {
				buffer.WriteRune(suffix)
				continue
			}

			tableMark := string(prefix) + string(suffix)
			asciiVal := slices.Index(c49Table7, tableMark)
			if asciiVal == -1 {
				if slices.Contains(nonDataChart, suffix) {
					buffer.WriteRune(prefix)
				} else {
					buffer.WriteString(string(suffix) + string(prefix))
				}
			} else {
				buffer.WriteRune(rune(asciiVal))
			}
		}
	}

	return buffer.String()
}
