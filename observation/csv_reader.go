package observation

import (
	"bufio"
	"io"
)

// CSVReader deserialises observations from an io.Reader containing CSV encoded observations.
type CSVReader struct {
	scanner  *bufio.Scanner
	rowIndex int64
}

// NewCSVReader returns a new CSVReader instance for the given io.CSVReader
func NewCSVReader(ioreader io.Reader) *CSVReader {

	scanner := bufio.NewScanner(ioreader)

	// Discard the header row.
	scanner.Scan()

	rowIndex := int64(1) // have discarded the header row so start at 1.

	return &CSVReader{
		scanner:  scanner,
		rowIndex: rowIndex,
	}
}

// Read will take a line from the input batchReader and convert it into an Observation instance.
func (reader *CSVReader) Read() (*Observation, error) {

	text, err := reader.readLine()
	if err != nil {
		return nil, err
	}

	observation := &Observation{
		Row:      text,
		RowIndex: reader.rowIndex,
	}

	reader.rowIndex = reader.rowIndex + 1

	return observation, nil
}

// ReadLine will read a single line from the input batchReader, returning an error if the read fails.
func (reader *CSVReader) readLine() (string, error) {
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
