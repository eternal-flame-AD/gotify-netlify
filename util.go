package main

import (
	"crypto/sha256"
	"fmt"
)

func sha256Str(b []byte) string {
	s := sha256.New()
	s.Write(b)
	return fmt.Sprintf("%x", s.Sum(nil))
}
