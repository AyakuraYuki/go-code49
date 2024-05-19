package code49

import (
	"slices"
	"strings"
)

// Decode decodes a bunch of numbers `bsLines` which transferred from bar/space data, returns basic text.
// The basic text means that text decoded by this method contains checksum mixin characters.
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
		bsLine := bsLines[row]
		payload := bsLine[2 : len(bsLine)-1]

		// For each symbol
		for col := 0; col < 4; col++ {
			value := 0
			parity := rowParity[row][col]
			block := payload[col<<3 : (col+1)<<3]
			if parity == 'E' || len(bsLines) == row+1 {
				value = slices.Index(evenEncodationPatterns, block)
			} else {
				value = slices.Index(oddEncodationPatterns, block)
			}

			// Extract 2 characters from value
			c1, c2 := value/49, value%49

			// Find and append buffer
			char, control := chars[c1], chars[c2]
			if slices.Contains(nonDataMethodChart, control) {
				buffer.WriteRune(char)
				continue
			}

			asciiChart := string(control) + string(char)
			asciiVal := slices.Index(c49Table7, asciiChart)
			if asciiVal == -1 {
				if slices.Contains(nonDataChart, char) {
					buffer.WriteRune(control)
				} else {
					buffer.WriteString(string(char) + string(control))
				}
			} else {
				buffer.WriteRune(rune(asciiVal))
			}
		}
	}

	return buffer.String()
}

// DecodeRaw decodes a bunch of numbers `bsLines` which transferred from bar/space data, returns raw text.
// It treads the Code49 as Mode-0, ignores other non-data characters.
// If `skipChecksum` presents to true, Decode will ignore the last line of `bsLines`.
func DecodeRaw(bsLines []string) string {
	var (
		chars     = []rune(charTable)
		hibcChars = []rune(hibcCharTable)
		cGrid     = make([][]int, 0)
	)

	for row := 0; row < len(bsLines); row++ {
		bsLine := bsLines[row]
		payload := bsLine[2 : len(bsLine)-1]

		var cGridRow []int
		for col := 0; col < 4; col++ {
			value := 0
			parity := rowParity[row][col]
			block := payload[col<<3 : (col+1)<<3]
			if parity == 'E' || len(bsLines) == row+1 {
				value = slices.Index(evenEncodationPatterns, block)
			} else {
				value = slices.Index(oddEncodationPatterns, block)
			}

			c1, c2 := value/49, value%49
			cGridRow = append(cGridRow, c1, c2)
		}
		cGrid = append(cGrid, cGridRow)
	}

	rows := len(bsLines)
	M := cGrid[rows-1][6] - (7 * (rows - 2))
	startingModeValue := 0
	if M == 2 {
		startingModeValue = 48
	} else if M == 4 {
		startingModeValue = 43
	} else if M == 5 {
		startingModeValue = 44
	}

	codewords := make([]int, 0)
	if M != 0 {
		codewords = append(codewords, startingModeValue)
	}
	for row := 0; row < len(cGrid)-1; row++ {
		for col := 0; col < len(cGrid[row])-1; col++ {
			codewords = append(codewords, cGrid[row][col])
		}
	}

	var intermediate []rune
	var stack strings.Builder
	for _, v := range codewords {
		if v == pad {
			continue
		}
		if stack.Len() == 2 {
			chart := stack.String()
			asciiVal := slices.Index(c49Table7, chart)
			intermediate = append(intermediate, rune(asciiVal))
			stack.Reset()
		}
		if v > len(hibcChars) {
			stack.WriteRune(chars[v])
			continue
		}
		if stack.Len() == 1 {
			stack.WriteRune(chars[v])
			continue
		}
		intermediate = append(intermediate, hibcChars[v])
		stack.Reset()
	}
	if stack.Len() == 2 {
		chart := stack.String()
		asciiVal := slices.Index(c49Table7, chart)
		intermediate = append(intermediate, rune(asciiVal))
	}

	return string(intermediate)
}
