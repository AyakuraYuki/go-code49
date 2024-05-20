package code49

const (
	FNC1 = -1
	FNC2 = -2
	FNC3 = -3
	FNC4 = -4

	FNC1String = "\\<FNC1>"
	FNC2String = "\\<FNC2>"
	FNC3String = "\\<FNC3>"
	FNC4String = "\\<FNC4>"
)

func toBytes(s string) []int {
	fnc1 := []byte(FNC1String)
	fnc2 := []byte(FNC2String)
	fnc3 := []byte(FNC3String)
	fnc4 := []byte(FNC4String)

	bytes := []byte(s)
	data := make([]int, len(bytes))

	i, j := 0, 0
	for ; i < len(bytes); i++ {
		if containsAt(bytes, fnc1, i) {
			data[j] = FNC1
			i += len(fnc1) - 1
		} else if containsAt(bytes, fnc2, i) {
			data[j] = FNC2
			i += len(fnc2) - 1
		} else if containsAt(bytes, fnc3, i) {
			data[j] = FNC3
			i += len(fnc3) - 1
		} else if containsAt(bytes, fnc4, i) {
			data[j] = FNC4
			i += len(fnc4) - 1
		} else {
			data[j] = int(bytes[i]) & 0xff
		}
		j++
	}

	if j < i {
		data = data[:j]
	}

	return data
}

// returns true if the specified array contains the specified sub-array at the specified index
func containsAt(array, searchFor []byte, index int) bool {
	for i := 0; i < len(searchFor); i++ {
		if index+i >= len(array) || array[index+i] != searchFor[i] {
			return false
		}
	}
	return true
}
