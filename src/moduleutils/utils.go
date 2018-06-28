package moduleutils

import (
	"os"
	"time"
)

func FileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func FileModTime(file string) time.Time {
	info, err := os.Stat(file)
	if err != nil {
		return time.Unix(0, 0)
	}
	return info.ModTime()
}
