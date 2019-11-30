package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {

	var (
		sourceFilePath string
		printHelp      bool
	)
	flag.StringVar(&sourceFilePath, "source", "",
		"The CNC source file path, e.g /Users/john/a.cnc")
	flag.BoolVar(&printHelp, "help", false, "Print help info")
	flag.Parse()
	if printHelp {
		fmt.Println(`use ./cnc-source-modifier --source C:\source.nc to modify the CNC file,
and new file will store at same folder with timestamp suffix, e.g C:\source-10086.nc`)
		os.Exit(0)
	}
	startAt := time.Now()
	patcher, err := NewPatcher(sourceFilePath, InsertG43HXBetweenTXM6)
	if err != nil {
		fmt.Printf("Unexpected error while initializing patcher: %s", err)
		os.Exit(1)
	}
	defer patcher.Close()
	// Open modified new file
	modifiedFilePath := ConstructCopyFilePath(sourceFilePath)
	modifiedFile, err := os.OpenFile(modifiedFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open new file: %s", err)
		os.Exit(1)
	}
	modifiedFileWriter := bufio.NewWriter(modifiedFile)
	defer modifiedFile.Close()
	err = patcher.PatchAndWrite(modifiedFileWriter)
	if err != nil {
		fmt.Printf("Unexpected error while patch and write new file: %s", err)
		os.Exit(1)
	}
	err = modifiedFileWriter.Flush()
	if err != nil {
		fmt.Printf("New file flush failed: %s", err)
	}
	fmt.Printf("Patch work done, new file path: %s\nTime used: %s\n", modifiedFilePath, time.Since(startAt))
}
