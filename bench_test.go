package disusechan

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"testing"
)

var testData = func() []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteString("x,y\n")
	for x := range 100 {
		_, _ = fmt.Fprintf(buf, "%d,%.6f\n", x, math.Sin(float64(x)))
	}
	return buf.Bytes()
}()

func newTestReader() io.Reader {
	return bytes.NewBuffer(testData)
}

func BenchmarkLinesChannel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func() {
			lines, errs := LinesChannel(newTestReader())
			for {
				select {
				case _, ok := <-lines:
					if !ok {
						return
					}
				case err := <-errs:
					if err != nil {
						panic(err)
					}
				}
			}
		}()
	}
}

func BenchmarkCSVChannel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func() {
			data, errs := CSVChannel(newTestReader())
			for {
				select {
				case _, ok := <-data:
					if !ok {
						return
					}
				case err := <-errs:
					if err != nil {
						panic(err)
					}
				}
			}
		}()
	}
}

func BenchmarkFetchLines(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FetchLines(newTestReader(), func(line []byte) error {
			return nil
		})
	}
}

func BenchmarkFetchCSV(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FetchCSV(newTestReader(), func(line [][]byte) error {
			return nil
		})
	}
}
