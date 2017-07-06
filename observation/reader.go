package observation

import (
	"bufio"
	"github.com/ONSdigital/dp-observation-extractor/model"
	"io"
)

// Reader provides an interface to read Observations
type Reader interface {
	Read() (*model.Observation, error)
}

// ensure reader satisfies the observation.Reader interface
var _ Reader = (*reader)(nil)

// NewReader returns a new Reader instance for the given io.Reader
func NewReader(ioreader io.Reader) Reader {

	scanner := bufio.NewScanner(ioreader)

	return &reader{
		scanner: scanner,
	}
}

type reader struct {
	scanner *bufio.Scanner
}

// Read will take a line from the input batchReader and convert it into an Observation instance.
func (reader *reader) Read() (*model.Observation, error) {

	text, err := reader.readLine()
	if err != nil {
		return nil, err
	}

	return &model.Observation{
		Row: text,
	}, nil
}

// ReadLine will read a single line from the input batchReader, returning an error if the read fails.
func (reader *reader) readLine() (string, error) {
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
