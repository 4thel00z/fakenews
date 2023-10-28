package fakenews

import (
	"bytes"
	"io"
	"testing"
)

func TestLineReaderSeeker(t *testing.T) {
	content := "Line1\nLine2\nLine3\nLine4\nLine5"
	reader := bytes.NewReader([]byte(content))
	lrs := NewLineReaderSeeker(reader)

	tests := []struct {
		lineNumber int
		whence     int
		expected   int
		name       string
	}{
		{3, io.SeekStart, 3, "SeekStart_Positive"},
		{-2, io.SeekEnd, 2, "SeekEnd_Negative"},
		{10, io.SeekStart, 5, "SeekStart_WrapAround"},
		{-10, io.SeekEnd, 5, "SeekEnd_WrapAround"},
		{0, io.SeekEnd, 5, "SeekEnd_Zero"},
		{0, io.SeekStart, 5, "SeekStart_Zero"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := lrs.SeekLine(test.lineNumber, test.whence)
			if err != nil {
				t.Fatalf("SeekLine failed: %v", err)
			}
			count, _ := lrs.Lines()
			if count != test.expected {
				t.Errorf("Expected %d lines, got %d", test.expected, count)
			}
		})
	}
}

func TestLineReaderSeekerEmptyFile(t *testing.T) {
	reader := bytes.NewReader([]byte(""))
	lrs := NewLineReaderSeeker(reader)

	err := lrs.SeekLine(1, io.SeekStart)
	if err == nil {
		t.Error("Expected error when seeking in an empty file, but got nil")
	}
}

func TestRandomLine(t *testing.T) {
	content := "Line1\nLine2\nLine3\nLine4\nLine5"
	reader := bytes.NewReader([]byte(content))
	lrs := NewLineReaderSeeker(reader)

	line, err := lrs.RandomLine()
	if err != nil {
		t.Fatalf("RandomLine failed: %v", err)
	}
	if line == "" {
		t.Error("Expected a non-empty line, but got an empty string")
	}
}

func TestRandomLineEmptyFile(t *testing.T) {
	reader := bytes.NewReader([]byte(""))
	lrs := NewLineReaderSeeker(reader)

	_, err := lrs.RandomLine()
	if err == nil {
		t.Error("Expected error when getting a random line from an empty file, but got nil")
	}
}
