package channelcallbacks

import (
	"bufio"
	"bytes"
	"io"
)

// LinesChannel returns channel with read lines and error-channel
func LinesChannel(f io.Reader) (lines chan []byte, errors chan error) {
	lines, errors = make(chan []byte), make(chan error)
	go func() {
		defer close(lines)
		defer close(errors)
		r := bufio.NewReader(f)
		for {
			if line, err := r.ReadSlice('\n'); err == io.EOF || err == nil {
				lines <- bytes.TrimSpace(line)
				if err == io.EOF {
					break
				}
			} else {
				errors <- err
				break
			}
		}
	}()
	return lines, errors
}

// CSVChannel returns channel with CSV-values of read CSV-line and error-channel
func CSVChannel(f io.Reader) (chValues chan [][]byte, chErrors chan error) {
	chValues, chErrors = make(chan [][]byte), make(chan error)
	go func() {
		defer close(chValues)
		defer close(chErrors)

		lines, errs := LinesChannel(f)
		for {
			select {
			case line, ok := <-lines:
				if !ok {
					return
				}
				chValues <- csvValues(line)

			case err := <-errs:
				if err != nil {
					chErrors <- err
					return
				}
			}
		}
	}()
	return chValues, chErrors
}

func csvValues(line []byte) (vals [][]byte) {
	for {
		if i := bytes.IndexByte(line, ','); i >= 0 {
			vals, line = append(vals, line[:i]), line[i+1:]
		} else {
			vals = append(vals)
			break
		}
	}
	return
}
