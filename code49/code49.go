package code49

import (
	"errors"
	"regexp"
	"slices"
	"strings"
)

// SetMode set data type in mode ECI, GS1 or HIBC.
// The default mode is ECI.
func SetMode(m DataType) {
	mode = m
}

// Encode encodes a given text into patterns in Code49
func Encode(text string) (patterns []string, encodationPatterns [][]int, err error) {
	var (
		codewords              = make([]int, 170)
		codewordCount          = 0
		padCount               = 0
		xCount, yCount, zCount int
		posnVal, localVal      int
		chars                  = []rune(charTable)
		cGrid                  [][]int // chart grid
		wGrid                  [][]int // symbol chart grid

		rowHeight []int
	)
	for i := 0; i < 8; i++ {
		cGrid = append(cGrid, make([]int, 8))
		wGrid = append(wGrid, make([]int, 4))
	}

	if !regexp.MustCompile("[\u0000-\u007F]+").MatchString(text) {
		return nil, nil, errors.New("invalid characters in input")
	}

	inputData := toBytes(text)
	intermediateBuilder := strings.Builder{}
	if mode == GS1 {
		intermediateBuilder.WriteRune('*') // FNC1
	}
	for i := 0; i < len(inputData); i++ {
		c := inputData[i]
		if c == FNC1 {
			intermediateBuilder.WriteRune('*') // FNC1
		} else {
			intermediateBuilder.WriteString(c49Table7[c])
		}
	}

	// ------ analyse code words ------

	intermediate := intermediateBuilder.String()
	h := len([]rune(intermediate))
	i := 0
	for {
		if intermediate[i] >= '0' && intermediate[i] <= '9' {

			// numeric data
			latch, j := 0, 0
			for {
				if (i + j) >= h {
					latch = 1
				} else {
					if intermediate[i+j] >= '0' && intermediate[i+j] <= '9' {
						j++
					} else {
						latch = 1
					}
				}

				if latch != 0 {
					break
				}
			}

			if j >= 5 {

				// use numeric encodation method
				var c, blockCount, blockRemain, blockValue int

				codewords[codewordCount] = 48 // <NS> (numeric shift)
				codewordCount++

				blockCount = j / 5
				blockRemain = j % 5

				for c = 0; c < blockCount; c++ {
					if (c == blockCount-1) && (blockRemain == 2) {
						// rule d
						blockValue = 100000
						blockValue += int(intermediate[i]-'0') * 1000
						blockValue += int(intermediate[i+1]-'0') * 100
						blockValue += int(intermediate[i+2]-'0') * 10
						blockValue += int(intermediate[i+3] - '0')

						codewords[codewordCount] = blockValue / (48 * 48)
						blockValue -= (48 * 48) * codewords[codewordCount]
						codewordCount++
						codewords[codewordCount] = blockValue / 48
						blockValue -= 48 * codewords[codewordCount]
						codewordCount++
						codewords[codewordCount] = blockValue
						codewordCount++
						i += 4
						blockValue = int(intermediate[i]-'0') * 100
						blockValue += int(intermediate[i+1]-'0') * 10
						blockValue += int(intermediate[i+2] - '0')

						codewords[codewordCount] = blockValue / 48
						blockValue -= 48 * codewords[codewordCount]
						codewordCount++
						codewords[codewordCount] = blockValue
						codewordCount++
					} else {
						blockValue = int(intermediate[i]-'0') * 10000
						blockValue += int(intermediate[i+1]-'0') * 1000
						blockValue += int(intermediate[i+2]-'0') * 100
						blockValue += int(intermediate[i+3]-'0') * 10
						blockValue += int(intermediate[i+4] - '0')

						codewords[codewordCount] = blockValue / (48 * 48)
						blockValue -= (48 * 48) * codewords[codewordCount]
						codewordCount++
						codewords[codewordCount] = blockValue / 48
						blockValue -= 48 * codewords[codewordCount]
						codewordCount++
						codewords[codewordCount] = blockValue
						codewordCount++
						i += 5
					}
				}

				switch blockRemain {
				case 1:
					// rule a
					codewords[codewordCount] = slices.Index(chars, rune(intermediate[i]))
					codewordCount++
					i++
				case 3:
					// rule b
					blockValue = int(intermediate[i]-'0') * 100
					blockValue += int(intermediate[i+1]-'0') * 10
					blockValue += int(intermediate[i+2] - '0')

					codewords[codewordCount] = blockValue / 48
					blockValue -= 48 * codewords[codewordCount]
					codewordCount++
					codewords[codewordCount] = blockValue
					codewordCount++
					i += 3
				case 4:
					// rule c
					blockValue = 100000
					blockValue += int(intermediate[i]-'0') * 1000
					blockValue += int(intermediate[i+1]-'0') * 100
					blockValue += int(intermediate[i+2]-'0') * 10
					blockValue += int(intermediate[i+3] - '0')

					codewords[codewordCount] = blockValue / (48 * 48)
					blockValue -= (48 * 48) * codewords[codewordCount]
					codewordCount++
					codewords[codewordCount] = blockValue / 48
					blockValue -= 48 * codewords[codewordCount]
					codewordCount++
					codewords[codewordCount] = blockValue
					codewordCount++
					i += 4
				}
				if i < h {
					// there's more to add
					codewords[codewordCount] = 48 // numeric shift
					codewordCount++
				}

			} else {
				codewords[codewordCount] = slices.Index(chars, rune(intermediate[i]))
				codewordCount++
				i++
			}

		} else {
			codewords[codewordCount] = slices.Index(chars, rune(intermediate[i]))
			codewordCount++
			i++
		}

		// do {...} while ()
		if i >= h {
			break
		}
	}

	// ------ start mode ------

	var M int
	switch codewords[0] { // set starting mode value
	case 48:
		M = 2
	case 43:
		M = 4
	case 44:
		M = 5
	default:
		M = 0
	}

	if M != 0 {
		for i = 0; i < codewordCount; i++ {
			codewords[i] = codewords[i+1]
		}
		codewordCount--
	}

	if codewordCount > 49 {
		return nil, nil, errors.New("input too long")
	}

	// ------ padding ------

	// place codewords in code character grid
	rows := 0
	for {
		for i = 0; i < 7; i++ {
			if ((rows * 7) + i) < codewordCount {
				cGrid[rows][i] = codewords[(rows*7)+i]
			} else {
				cGrid[rows][i] = pad // Pad
				padCount++
			}
		}
		rows++

		if rows*7 >= codewordCount {
			break
		}
	}

	if ((rows <= 6) && (padCount < 5)) || (rows > 6) || (rows == 1) {
		// add a row
		for i = 0; i < 7; i++ {
			cGrid[rows][i] = pad // Pad
		}
		rows++
	}

	// add row count and mode character
	cGrid[rows-1][6] = (7 * (rows - 2)) + M

	// ------ checksum ------

	// add checksum for each row
	for i = 0; i < rows-1; i++ {
		rowSum := 0
		for j := 0; j < 7; j++ {
			rowSum += cGrid[i][j]
		}
		cGrid[i][7] = rowSum % 49
	}

	// calculate symbol check characters
	posnVal = 0
	xCount = cGrid[rows-1][6] * xWeight00
	yCount = cGrid[rows-1][6] * yWeight00
	zCount = cGrid[rows-1][6] * zWeight00
	for i = 0; i < rows-1; i++ {
		for j := 0; j < 4; j++ {
			localVal = (cGrid[i][2*j] * 49) + cGrid[i][(2*j)+1]
			xCount += xWeights[posnVal] * localVal
			yCount += yWeights[posnVal] * localVal
			zCount += zWeights[posnVal] * localVal
			posnVal++
		}
	}

	if rows > 6 {
		// add z symbol check
		cGrid[rows-1][0] = (zCount % 2401) / 49
		cGrid[rows-1][1] = (zCount % 2401) % 49
	}

	localVal = (cGrid[rows-1][0] * 49) + cGrid[rows-1][1]
	xCount += xWeights[posnVal] * localVal
	yCount += yWeights[posnVal] * localVal
	posnVal++

	// add y symbol check
	cGrid[rows-1][2] = (yCount % 2401) / 49
	cGrid[rows-1][3] = (yCount % 2401) % 49

	localVal = (cGrid[rows-1][2] * 49) + cGrid[rows-1][3]
	xCount += xWeights[posnVal] * localVal

	// add x symbol check
	cGrid[rows-1][4] = (xCount % 2401) / 49
	cGrid[rows-1][5] = (xCount % 2401) % 49

	// add last row as checksum
	sum := 0
	for i = 0; i < 7; i++ {
		sum += cGrid[rows-1][i]
	}
	cGrid[rows-1][7] = sum % 49

	// transfer data to symbol character grid
	for i = 0; i < rows; i++ {
		for j := 0; j < 4; j++ {
			wGrid[i][j] = (cGrid[i][2*j] * 49) + cGrid[i][(2*j)+1]
		}
	}

	patterns = make([]string, rows)
	encodationPatterns = make([][]int, rows)
	rowHeight = make([]int, rows)

	for i = 0; i < rows; i++ {
		rowPattern := strings.Builder{}
		rowPattern.WriteString("11") // pattern start with prefix '11'
		for j := 0; j < 4; j++ {
			symbolChart := wGrid[i][j]
			encodationPatterns[i] = append(encodationPatterns[i], symbolChart)
			if i != (rows - 1) {
				if rowParity[i][j] == 'E' {
					// even parity
					rowPattern.WriteString(evenEncodationPatterns[symbolChart])
				} else {
					// odd parity
					rowPattern.WriteString(oddEncodationPatterns[symbolChart])
				}
			} else {
				// last row uses all even parity
				rowPattern.WriteString(evenEncodationPatterns[symbolChart])
			}
		}
		rowPattern.WriteString("4") // patter stop with suffix '4'
		patterns[i] = rowPattern.String()
		rowHeight[i] = 10
	}

	return patterns, encodationPatterns, nil
}

// Decode decodes barcode `patterns` which contains multiple rows with scanned bar/space amounts,
// start with prefix `11` and end with suffix `4`, returns basic text.
// The basic text means that text decoded by this method contains checksum mixin characters.
// It treads the Code49 as Mode-0, ignores other non-data characters.
// If `skipChecksum` presents to true, Decode will ignore the last line of `patterns`.
func Decode(patterns []string, skipChecksum bool) string {
	var (
		buffer strings.Builder
		chars  = []rune(charTable)
	)

	// For each line expect last line, as ignoring checksum
	bound := len(patterns)
	if skipChecksum {
		bound = len(patterns) - 1
	}
	for row := 0; row < bound; row++ {
		// Extract pattern, remove head '11' and tail '4'
		bsLine := patterns[row]
		payload := bsLine[2 : len(bsLine)-1]

		// For each symbol
		for col := 0; col < 4; col++ {
			value := 0
			parity := rowParity[row][col]
			block := payload[col<<3 : (col+1)<<3]
			if parity == 'E' || len(patterns) == row+1 {
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

// DecodeRaw decodes barcode `patterns` which contains multiple rows with scanned bar/space amounts,
// start with prefix `11` and end with suffix `4`, returns raw text.
// It treads the Code49 as Mode-0, ignores other non-data characters.
func DecodeRaw(patterns []string) string {
	var (
		chars     = []rune(charTable)
		hibcChars = []rune(hibcCharTable)
		cGrid     = make([][]int, 0)
	)

	// For each line decode into a chart grid
	for row := 0; row < len(patterns); row++ {
		bsLine := patterns[row]
		payload := bsLine[2 : len(bsLine)-1]

		var cGridRow []int
		for col := 0; col < 4; col++ {
			value := 0
			parity := rowParity[row][col]
			block := payload[col<<3 : (col+1)<<3]
			if parity == 'E' || len(patterns) == row+1 {
				value = slices.Index(evenEncodationPatterns, block)
			} else {
				value = slices.Index(oddEncodationPatterns, block)
			}

			c1, c2 := value/49, value%49
			cGridRow = append(cGridRow, c1, c2)
		}
		cGrid = append(cGrid, cGridRow)
	}

	// Analyse the starting mode from checksum
	rows := len(patterns)
	M := cGrid[rows-1][6] - (7 * (rows - 2))
	startingModeValue := 0
	if M == 2 {
		startingModeValue = 48
	} else if M == 4 {
		startingModeValue = 43
	} else if M == 5 {
		startingModeValue = 44
	}

	// Recovery chart code words
	codewords := make([]int, 0)
	if M != 0 {
		codewords = append(codewords, startingModeValue)
	}
	for row := 0; row < len(cGrid)-1; row++ {
		for col := 0; col < len(cGrid[row])-1; col++ {
			codewords = append(codewords, cGrid[row][col])
		}
	}

	// Use Table7 to transfer code words into plain text
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
