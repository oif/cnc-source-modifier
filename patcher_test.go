package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

type writer struct {
	content []byte
}

func (w *writer) Write(p []byte) (n int, err error) {
	w.content = append(w.content, p...)
	return len(p), nil
}

func TestPatcher(t *testing.T) {
	patcher, err := NewPatcher("./case/origin.nc",
		InsertG43HXBetweenTXM6)
	if err != nil {
		t.Fatalf("Failed to initialize patcher: %s", err)
	}
	content := new(writer)
	modifiedFileWriter := bufio.NewWriter(content)
	defer func() {
		patcher.Close()
	}()
	err = patcher.PatchAndWrite(modifiedFileWriter)
	if err != nil {
		t.Fatalf("Failed to patch and write new file: %s", err)
	}
	err = modifiedFileWriter.Flush()
	if err != nil {
		fmt.Printf("New file flush failed: %s", err)
	}
	// Compare result
	sampleFile, err := ioutil.ReadFile("./case/modified.nc")
	if err != nil {
		t.Fatalf("Failed to read sample file: %s", err)
	}
	if !bytes.Equal(content.content, sampleFile) {
		t.Fatal("Test result file not equal with sample file as expect")
	}
}
