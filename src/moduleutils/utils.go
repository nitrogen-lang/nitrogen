package moduleutils

import "os"

func FileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}
