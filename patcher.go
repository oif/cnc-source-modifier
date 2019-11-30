package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type PatchFunc func(line string, fileWriter io.Writer) (hint bool, err error)

func InsertG43HXBetweenTXM6(line string, fileWriter io.Writer) (bool, error) {
	blocks, ok := SplitLinePositionBlock(line)
	if !ok || len(blocks) != 2 || blocks[1] != "M6" || len(blocks[0]) != 2 || blocks[0][0] != 'T' {
		return false, nil
	}
	if !unicode.IsDigit(rune(blocks[0][1])) {
		return false, nil
	}
	TXDigitStr := string(blocks[0][1])
	TXDigit, err := strconv.Atoi(TXDigitStr)
	if err != nil {
		return false, nil
	}
	// Hint
	newBlocks := make([]string, 4)
	newBlocks[0] = blocks[0]
	newBlocks[1] = "G43"
	newBlocks[2] = fmt.Sprintf("H%d", TXDigit)
	newBlocks[3] = blocks[1]
	_, err = fileWriter.Write([]byte(strings.Join(newBlocks, " ") + "\n"))
	return true, err
}

type Patcher struct {
	sourceFilePath string
	sourceFile     *os.File
	patchChain     []PatchFunc
}

func NewPatcher(sourceFilePath string, patchFuncs ...PatchFunc) (*Patcher, error) {
	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return nil, err
	}
	return &Patcher{
		sourceFilePath: sourceFilePath,
		sourceFile:     sourceFile,
		patchChain:     patchFuncs,
	}, nil
}

func (p *Patcher) PatchAndWrite(w io.Writer) error {
	lineScanner := bufio.NewScanner(p.sourceFile)
	var lineNeverHint bool
	for lineScanner.Scan() {
		line := lineScanner.Text()
		lineNeverHint = true
		for _, patchFunc := range p.patchChain {
			hint, err := patchFunc(line, w)
			if err != nil {
				return err
			}
			if hint {
				// Next line
				lineNeverHint = false
				break
			}
		}
		if lineNeverHint {
			// Manual write to new file
			if line != "\n" && line != "\r\n" {
				line += "\n"
			}
			_, err := w.Write([]byte(line))
			if err != nil {
				return err
			}
		}
	}
	if err := lineScanner.Err(); err != nil {
		return err
	}
	return nil
}

func (p *Patcher) Close() {
	p.sourceFile.Close()
}
