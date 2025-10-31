package common

import (
	"errors"
	"os"
)

func GetGopath() (string, error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return "", errors.New("GOPATH is not set")
	}
	return gopath, nil
}

func GetGobin() (string, error) {
	gobin := os.Getenv("GOBIN")
	if gobin == "" {
		return "", errors.New("GOBIN is not set")
	}
	return gobin, nil
}
