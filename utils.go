package main

import (
	"fmt"
	"strings"
	"time"
	"unicode"
)

func ConstructCopyFilePath(originPath string) string {
	timestamp := time.Now().Unix()
	pathWithoutSuffix := originPath
	if strings.HasSuffix(originPath, ".nc") {
		pathWithoutSuffix = originPath[:len(originPath)-3]
	}
	return fmt.Sprintf("%s-%d.nc", pathWithoutSuffix, timestamp)
}

func SplitLinePositionBlock(line string) ([]string, bool) {
	if len(line) == 0 {
		return nil, false
	}
	// Check first character is A-Z
	if unicode.IsLetter(rune(line[0])) {
		return strings.Split(line, " "), true
	}
	return nil, false
}
