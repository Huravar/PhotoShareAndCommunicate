package gadgets

import "strconv"

func StringToUint(PString string) uint {
	IntParam, _ := strconv.Atoi(PString)
	return uint(IntParam)
}
