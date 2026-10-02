package channelcallbacks

import (
	"bufio"
	"bytes"
	"io"
)

// FetchLines fetches lines from given reader
func FetchLines(r io.Reader, fn func(line []byte) error) error {
	buf := bufio.NewReader(r)
	for {
		line, err := buf.ReadSlice('\n')
		eof := err == io.EOF
		if eof || err == nil {
			err = fn(bytes.TrimSpace(line))
		}
		if eof || err != nil {
			return err
		}
	}
}

// FetchCSV fetches CSV-values from given reader (CSV-file)
func FetchCSV(r io.Reader, fn func(values [][]byte) error) error {
	return FetchLines(r, func(line []byte) error {
		return fn(csvValues(line))
	})
}
