# gimel

[![GoDoc][godoc-badge]][godoc]

[godoc-badge]: https://godoc.org/github.com/stephen-fox/gimel?status.svg
[godoc]: https://godoc.org/github.com/stephen-fox/gimel

Package gimel provides functionality for in-memory execution on Linux.
It can also load arbitrary files into memory, which (in certain cases)
can be shared with other Linux processes via file path or by passing the
file descriptor via exec.

This library is based on example code by @magisterquis. Please give their
[excellent blog post](https://magisterquis.github.io/2018/03/31/in-memory-only-elf-execution.html)
on the subject a read.

## APIs

#### `memfd_create(2)` wrappers
The [memfd_create(2)](https://man7.org/linux/man-pages/man2/memfd_create.2.html)
system call copies data into memory. The library offers several functions for
executing the system call:

- `MemfdCreate()` - Executes `memfd_create(2)`, returning a file descriptor
representing the in-memory file
- `MemfdCreateOSFile()` - Executes `memfd_create(2)`, returning a *os.File
representing the in-memory file
- `MemfdCreateFromExe()` - Executes `memfd_create(2)` and copies the specified
executable into the in-memory file, returning a *os.File and *exec.Cmd
- `MemfdCreateFromFile()` - Executes `memfd_create(2)` and copies the specified
file into the in-memory file, returning a *os.File
- `MemfdCreateFromReader()` - Executes `memfd_create(2)` and copies the contents
of the specified io.Reader into the in-memory file, returning a *os.File

#### Helper functions
Several helper functions are also available:

- `InMemoryFileToCmd()` - Converts an existing *os.File into a *exec.Cmd
- `CopyDataIntoMemFile()` - Copies the contents of an io.ReadCloser into the
specified *os.File and closes the io.ReadCloser

## Examples
The following examples can be found in the [examples/ directory](examples/):

- [examples/runonce/](examples/runonce/main.go) - Run an application from
memory only one time
- [examples/multirun/](examples/multirun/main.go) - Run an application from
memory several times in a row, reusing the existing in-memory file
- [examples/sharefile/](examples/sharefile/main.go) - Load an arbitrary file
into memory, print its "/proc" file path so that other programs can access it
using the file path
