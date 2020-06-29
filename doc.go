// Package gimel provides functionality for in-memory execution on Linux.
// It is primarily a wrapper for the 'memfd_create' system call, although
// it could easily be extended to utilize other strategies for running
// executables from memory.
//
// This library is based on example code by Stuart "MagisterQuis":
// https://magisterquis.github.io/2018/03/31/in-memory-only-elf-execution.html
package gimel
