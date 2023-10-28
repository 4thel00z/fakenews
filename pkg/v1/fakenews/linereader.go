package fakenews

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"math/rand"
)

type LineReaderSeeker struct {
	rs io.ReadSeeker
}

func NewLineReaderSeeker(rs io.ReadSeeker) *LineReaderSeeker {
	return &LineReaderSeeker{rs: rs}
}

func (lrs *LineReaderSeeker) Read(p []byte) (n int, err error) {
	return lrs.rs.Read(p)
}

func (lrs *LineReaderSeeker) Seek(offset int64, whence int) (int64, error) {
	return lrs.rs.Seek(offset, whence)
}

func (lrs *LineReaderSeeker) SeekLine(lineNumber int, whence int) error {
	switch whence {
	case io.SeekStart:
		if lineNumber < 0 {
			return lrs.seekFromEnd(lineNumber)
		}
		return lrs.seekFromStart(lineNumber)

	case io.SeekEnd:
		return lrs.seekFromEnd(lineNumber)

	default:
		return errors.New("unsupported seek method")
	}
}

func (lrs *LineReaderSeeker) seekFromEnd(lineNumber int) error {
	_, err := lrs.rs.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	if lineNumber == 0 {
		return nil
	}

	if lineNumber < 0 {
		lineNumber = -lineNumber
	}

	scanner := bufio.NewScanner(lrs.rs)
	scanner.Split(scanLinesReverse)

	lineCount := 0
	for scanner.Scan() {
		lineCount++
		if lineCount == lineNumber {
			return nil
		}
	}

	if lineCount > 0 {
		return lrs.seekFromEnd(-(lineNumber % lineCount))
	}

	return errors.New("no lines found in the file")
}

func scanLinesReverse(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.LastIndexByte(data, '\n'); i >= 0 {
		return len(data) - i, data[i+1:], nil
	}

	if atEOF {
		return len(data), data, nil
	}

	return 0, nil, nil
}

func (lrs *LineReaderSeeker) RandomLine() (string, error) {
	_, err := lrs.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	totalLines, err := lrs.Lines()
	if err != nil {
		return "", err
	}

	if totalLines == 0 {
		return "", errors.New("no lines found in the file")
	}

	randomLineNumber := rand.Intn(totalLines) + 1

	_, err = lrs.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	err = lrs.seekFromStart(randomLineNumber)
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(lrs.rs)
	if scanner.Scan() {
		return scanner.Text(), nil
	}

	return "", scanner.Err()
}

func (lrs *LineReaderSeeker) seekFromStart(lineNumber int) error {
	if lineNumber <= 0 {
		return errors.New("line number must be positive")
	}

	_, err := lrs.rs.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(lrs.rs)
	for lineNumber > 1 {
		if !scanner.Scan() {
			return errors.New("line number out of range")
		}
		lineNumber--
	}

	if scanner.Scan() {
		lineLength := len(scanner.Bytes())
		_, err := lrs.rs.Seek(int64(-lineLength), io.SeekCurrent)
		return err
	}

	return errors.New("no lines found in the file")
}

// ... Rest of the code ...

// Updated Lines function
func (lrs *LineReaderSeeker) Lines() (int, error) {
	currentPos, err := lrs.rs.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	lineCount, err := countLines(lrs.rs)
	if err != nil {
		return 0, err
	}

	_, err = lrs.rs.Seek(currentPos, io.SeekStart)
	return lineCount, err
}

// Helper function to count lines
func countLines(rs io.ReadSeeker) (int, error) {
	scanner := bufio.NewScanner(rs)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}
	if err := scanner.Err(); err != nil {
		if err == io.EOF {
			return lineCount, nil
		}
		return 0, err
	}
	return lineCount, nil
}
