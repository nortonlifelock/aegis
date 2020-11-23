package qualys

import "strconv"

// Convert a integer array to a string array
func intArrayToStringArray(intIn []int) (stringOut []string) {
	stringOut = make([]string, 0)

	if intIn != nil && len(intIn) > 0 {

		for _, value := range intIn {
			stringOut = append(stringOut, strconv.Itoa(value))
		}
	}

	return stringOut
}
