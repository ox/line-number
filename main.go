package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
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

func formatLine(lineNum int, maxLineGutterWidth int, line string) string {
	lineStr := fmt.Sprintf("%d", lineNum)
	for i := len(lineStr); i < maxLineGutterWidth; i++ {
		lineStr = " " + lineStr
	}
	return " " + lineStr + " | " + line
}

func writeFromPipe(f *os.File, w io.Writer) error {
	// Buffer the input so we can count the number of lines
	var buf bytes.Buffer
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	numLines := 0
	for scanner.Scan() {
		buf.WriteString(scanner.Text() + "\n")
		numLines++
	}

	maxLineGutterWidth := len(fmt.Sprintf("%d", numLines))

	for lineNum := 1; ; lineNum++ {
		line, err := buf.ReadString('\n')

		if len(line) > 0 {
			output := formatLine(lineNum, maxLineGutterWidth, line)
			w.Write([]byte(output))
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
	}
}

func writeFromFile(f *os.File, w io.Writer) error {
	// Count the number of lines in the input file
	numLines := 0
	s := bufio.NewScanner(f)
	for s.Scan() {
		numLines += 1
	}
	maxLineGutterWidth := len(fmt.Sprintf("%d", numLines))

	// Start from the beginning and write to the output
	f.Seek(0, 0)
	s = bufio.NewScanner(f)

	lineNum := 1
	for s.Scan() {
		output := formatLine(lineNum, maxLineGutterWidth, s.Text()) + "\n"
		_, err := w.Write([]byte(output))
		if err != nil {
			return fmt.Errorf("could not write to output: %w", err)
		}
		lineNum += 1
	}
	return nil
}

func main() {
	inputPath := ""
	outputPath := ""
	flag.StringVar(&inputPath, "in", "-", "Input file, stdin if not defined")
	flag.StringVar(&outputPath, "out", "-", "Output file, stdout if not defined")
	flag.Parse()

	outFile, err := openFile(outputPath, true, os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Check if the input is a pipe
	if inputPath == "-" {
		err = writeFromPipe(os.Stdin, outFile)
	} else {
		inFile, openErr := openFile(inputPath, false, os.Stdin)
		if openErr != nil {
			fmt.Println(openErr)
			os.Exit(1)
		}

		err = writeFromFile(inFile, outFile)
	}

	if err != nil {
		fmt.Printf("error writing to output: %s", err)
		os.Exit(1)
	}
}
