package main

import (
	"bufio"
	"encoding/hex"
	"os"
)

func ReadHexLines(file string) ([][]byte, error) {
	return ReadLinesWithParser(file, func(line string, lineNum int) ([]byte, error) {
		return hex.DecodeString(line)
	})
}

func ReadLinesWithParser(file string, parse func(line string, lineNum int) ([]byte, error)) ([][]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	lines := [][]byte{}
	scanner := bufio.NewScanner(f)
	lineCtr := 0
	for scanner.Scan() {
		lineCtr += 1
		line := scanner.Text()

		parsedLine, err := parse(line, lineCtr)
		if err != nil {
			return nil, err
		}

		if parsedLine == nil {
			continue
		}

		lines = append(lines, parsedLine)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
