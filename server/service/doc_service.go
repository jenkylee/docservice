package service

import (
	"errors"
	"strings"
)

type DocService struct {}

func (DocService) Import (s string) (string, error){
	if s == "" {
		return "", ErrEmpty
	}

	return strings.ToUpper(s), nil
}

func (DocService) Export(s string) int {
	return len(s)
}

var ErrEmpty = errors.New("empty strin")