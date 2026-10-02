# Disuse chan

Channels are useful for coordinating concurrent work, but they come with a cost. When a caller simply needs to read and process a sequence of values, a callback or a reader interface can make the API simpler and faster.

This repository compares two ways to process the same input in Go: a goroutine that sends each item through a channel, and a synchronous function that calls a callback for each item.

## Choose an API that fits the work

For sequential processing:

### ❌ Bad: a channel for every item

A channel-based API moves data between a producer goroutine and its consumer:

```go
func DataChannel() (<-chan Data, <-chan error)
```

The caller needs to receive values, handle errors, and detect when the producer finishes. Each item also requires communication between goroutines.

### ✅ Good: a direct callback

A callback-based API keeps reading and processing in the same goroutine:

```go
func FetchData(fn func(Data) error) error
```

The function calls the callback as values become available and returns when it finishes or encounters an error. The callback can stop processing by returning an error. For sequential work, this avoids channel synchronization and keeps error handling straightforward. A reader or iterator interface is another natural fit when the caller should control when to read the next value.

Use channels when concurrency is part of the task. For a sequential data stream, start with a callback, reader, or iterator and measure before adding goroutines.

## Benchmarks

The benchmarks compare line reading and simple comma-separated value splitting. Each operation processes the same generated input: a header and 100 rows. The callbacks and channel consumers discard the output, so the results measure reading and delivery overhead rather than application processing.

Measured on an Apple M2 with Go 1.27.0, macOS (`darwin/arm64`), and `GOMAXPROCS=8`. Times below are the median of three runs; memory figures are per operation.

| Workload | ❌ Chan (ns/op) | ✅ Fetch (ns/op) | Fetch speedup |
| --- | ---: | ---: | ---: |
| Lines | 24,781 | **1,766** | **14×** |
| CSV | 66,461 | **3,696** | **18×** |

Fetch also uses less memory and fewer allocations:

| Workload | Chan bytes/op | Fetch bytes/op | Chan allocs/op | Fetch allocs/op |
| --- | ---: | ---: | ---: | ---: |
| Lines | 4,432 | **4,144** | 7 | **2** |
| CSV | 7,146 | **6,568** | 113 | **103** |

Speedups are rounded and describe the implementations and input in this repository; hardware, buffering, batch size, and work performed by the consumer can change the result.

## Run the benchmarks

```sh
git clone https://github.com/denisskin/disusechan.git
cd disusechan
go test -run '^$' -bench . -benchmem -count=3
```

The benchmark data is generated in memory, so no external data files are required. Each reported operation processes the entire input, not a single row.
