package observation

import (
	"io"
	"bufio"
)

// CSVReader deserialises observations from an io.Reader containing CSV encoded observations.
type CSVReader struct {
	scanner *bufio.Scanner
}

// NewCSVReader returns a new CSVReader instance for the given io.CSVReader
func NewCSVReader(ioreader io.Reader, discardHeaderRow bool) *CSVReader {

	scanner := bufio.NewScanner(ioreader)

	if discardHeaderRow {
		scanner.Scan()
	}

	return &CSVReader{
		scanner: scanner,
	}
}

// Read will take a line from the input batchReader and convert it into an Observation instance.
func (reader CSVReader) Read() (*Observation, error) {

	text, err := reader.readLine()
	if err != nil {
		return nil, err
	}

	return &Observation{
		Row: text,
	}, nil
}

// ReadLine will read a single line from the input batchReader, returning an error if the read fails.
func (reader CSVReader) readLine() (string, error) {
	scanner := reader.scanner
	scanSuccessful := scanner.Scan()

	if !scanSuccessful {

		if scanner.Err() != nil {
			return "", scanner.Err()
		}

		// Scanner can return a nil error for EOF
		return "", io.EOF
	}

	return scanner.Text(), nil
}
