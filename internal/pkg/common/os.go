package common

import (
	"errors"
	"runtime"
	"strconv"
)

var (
	ErrUnsupportedOS   = errors.New("unsupported os. Please use darwin/linux/windows")
	ErrUnsupportedArch = errors.New("unsupported architecture. Please use 64 bit system")
)

func OsArchSupported() error {
	switch runtime.GOOS {
	case "darwin", "linux", "windows":
		return nil
	default:
		return ErrUnsupportedOS
	}
	if strconv.IntSize == 64 {
		return nil
	}
	return ErrUnsupportedArch
}

func IsUnix() bool {
	return runtime.GOOS == "linux" || runtime.GOOS == "darwin"
}
