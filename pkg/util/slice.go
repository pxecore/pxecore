package util

// RemoveStringFromSlice removes one string from the slice.
func RemoveStringFromSlice(slice []string, elem string) []string {
	for i, e := range slice {
		if e == elem {
			slice[i] = slice[len(slice)-1]
			slice[len(slice)-1] = ""
			slice = slice[:len(slice)-1]
			break
		}
	}
	return slice
}

// AddUniqueStringToSlice add one string from the slice.
func AddUniqueStringToSlice(slice []string, elem string) []string {
	for _, e := range slice {
		if e == elem {
			return slice
		}
	}
	return append(slice, elem)
}
