package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func openFile(path string, create bool, fallback *os.File) (*os.File, error) {
	var file *os.File
	var err error
	if path == "-" {
		return fallback, nil
	} else {
		perms := os.O_RDWR
		if create {
			perms = perms | os.O_CREATE
		}
		file, err = os.OpenFile(path, perms, 0644)
		if err != nil {
			return nil, fmt.Errorf("could not open file %s: %w", path, err)
		}
	}

	return file, nil
}

func main() {
	inputPath := ""
	outputPath := ""
	flag.StringVar(&inputPath, "in", "-", "Input file, stdin if not defined")
	flag.StringVar(&outputPath, "out", "-", "Output file, stdout if not defined")
	flag.Parse()

	in, err := openFile(inputPath, false, os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	out, err := openFile(outputPath, true, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Count the number of lines in the input file
	numLines := 0
	s := bufio.NewScanner(in)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		numLines += 1
	}
	maxLineGutterWidth := len(fmt.Sprintf("%d", numLines))

	// Start from the beginning and write to the output
	in.Seek(0, 0)
	s = bufio.NewScanner(in)
	s.Split(bufio.ScanLines)
	lineNum := 1
	for s.Scan() {
		lineStr := fmt.Sprintf("%d", lineNum)
		for i := len(lineStr); i < maxLineGutterWidth; i++ {
			lineStr = " " + lineStr
		}
		output := " " + lineStr + " | " + s.Text()
		_, err := out.WriteString(output + "\n")
		if err != nil {
			fmt.Printf("could not write to output: %s\n", err)
			os.Exit(1)
		}
		lineNum += 1
	}
}
