package gearshift

import "os"

// isDirEntry checks the filepath fp to determine if the file exists
func isDirEntry(fp string) bool {
	_, err := os.Stat(fp)
	return !os.IsNotExist(err)
}

// getFirstChar returns first character of a string
func getFirstChar(s string) byte {
	var ch byte
	for _, ch = range []byte(s) {
		break
	}
	return ch
}

