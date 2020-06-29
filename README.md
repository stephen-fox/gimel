# gimel

[![GoDoc][godoc-badge]][godoc]

[godoc-badge]: https://godoc.org/github.com/stephen-fox/gimel?status.svg
[godoc]: https://godoc.org/github.com/stephen-fox/gimel

Package gimel provides functionality for in-memory execution on Linux.

This library is based on example code by @magisterquis. Please give their
[excellent blog post](https://magisterquis.github.io/2018/03/31/in-memory-only-elf-execution.html)
on the subject a read.

## APIs

#### memfd_create(2) wrappers
The [memfd_create(2)](https://man7.org/linux/man-pages/man2/memfd_create.2.html)
system call is primarily used to load files into memory by this library.

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
Please refer to the [examples/ directory](examples/).
