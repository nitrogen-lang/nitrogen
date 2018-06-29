package moduleutils

import (
	"os"
	"time"
)

func FileExists(file string) bool {
	stat, err := os.Stat(file)
	return !os.IsNotExist(err) && !stat.IsDir()
}

func IsDir(file string) bool {
	stat, err := os.Stat(file)
	return !os.IsNotExist(err) && stat.IsDir()
}

func FileModTime(file string) time.Time {
	info, err := os.Stat(file)
	if err != nil {
		return time.Unix(0, 0)
	}
	return info.ModTime()
}
