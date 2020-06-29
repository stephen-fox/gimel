package gimel

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"
)

const (
	memfdCreateID = 319
)

// Various 'memfd_create(2)' constants.
//
// Refer to the man page for more information:
// https://man7.org/linux/man-pages/man2/memfd_create.2.html
const (
	// MfdCloExec specifies the "close-on-exec" flag.
	//
	// From the man page:
	//	Set the close-on-exec (FD_CLOEXEC) flag on the new file
	//	descriptor.  See the description of the O_CLOEXEC flag in
	//	open(2) for reasons why this may be useful.
	MfdCloExec = 1
)

// MemfdCreateFromExe creates a RAM-backed file using the 'memfd_create(2)'
// system call, and copies the specified executable file into it, returning
// a *exec.Cmd and a *os.File representing the in-memory file. Callers should
// close the *os.File only after they are finished running the executable.
//
// Note that a *exec.Cmd can only be run once. However, the in-memory file
// can be reused to create a new *exec.Cmd.
//
// Refer to MemfdCreate() for more information.
func MemfdCreateFromExe(optionalDisplayName string, exeFilePath string, args ...string) (*exec.Cmd, *os.File, error) {
	inMemory, err := MemfdCreateFromFile(optionalDisplayName, MfdCloExec, exeFilePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to call 'memfd_create' - %w", err)
	}

	return InMemoryFileToCmd(inMemory, args...), inMemory, nil
}

// MemfdCreateFromFile creates a RAM-backed file using the 'memfd_create(2)'
// system call, and copies the specified file into it, returning a *os.File
// representing the in-memory file. Callers should close the *os.File after
// they are finished using it.
//
// Refer to MemfdCreate() for more information.
func MemfdCreateFromFile(optionalDisplayName string, flags uint, sourcePath string) (*os.File, error) {
	source, err := os.Open(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file for reading - %w", err)
	}

	inMemory, err := MemfdCreateFromReader(optionalDisplayName, flags, source)
	source.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to create in-memory file from existing file on disk - %w", err)
	}

	return inMemory, nil
}

// MemfdCreateFromReader creates a RAM-backed file using the 'memfd_create(2)'
// system call, and copies the io.Reader into it, returning a *os.File
// representing the in-memory file. Callers should close the *os.File after
// they are finished using it.
//
// Refer to MemfdCreate() for more information.
func MemfdCreateFromReader(optionalDisplayName string, flags uint, source io.Reader) (*os.File, error) {
	inMemory, err := MemfdCreateOSFile(optionalDisplayName, flags)
	if err != nil {
		return nil, fmt.Errorf("failed to call 'memfd_create' - %w", err)
	}

	_, err = io.Copy(inMemory, source)
	if err != nil {
		return nil, fmt.Errorf("failed to copy source to in-memory file - %w", err)
	}

	return inMemory, nil
}

// MemfdCreateOSFile wraps the MemfdCreate() function, returning a *os.File
// with a properly populated name rather than a raw file descriptor.
// The *os.File returned by this function represents a RAM-backed file.
// Callers should close the *os.File only after all dependent resources are
// finished with it.
//
// Refer to MemfdCreate() for more information.
func MemfdCreateOSFile(optionalDisplayName string, flags uint) (*os.File, error) {
	fd, err := MemfdCreate(optionalDisplayName, flags)
	if err != nil {
		return nil, err
	}

	memFile := os.NewFile(fd, filepath.Join(
		"/proc",
		fmt.Sprintf("%d", os.Getpid()),
		"fd",
		fmt.Sprintf("%d", fd)))
	if memFile == nil {
		return nil, fmt.Errorf("os.NewFile returned nil when given mem fd %d", fd)
	}

	return memFile, nil
}

// MemfdCreate executes the 'memfd_create(2)' syscall.
//
// From the man page:
//	memfd_create() creates an anonymous file and returns a file
//	descriptor that refers to it. The file behaves like a regular file,
//	and so can be modified, truncated, memory-mapped, and so on.
//	However, unlike a regular file, it lives in RAM and has a volatile
//	backing storage.  Once all references to the file are dropped, it is
//	automatically released.
//
// The optionalDisplayName specifies what the RAM-backed file should be named.
// Per the 'memfd_create' documentation, the display name will always be
// prefixed by the string 'memfd:', even if a display name is specified.
//
// Refer to the man page for more information:
// https://man7.org/linux/man-pages/man2/memfd_create.2.html
func MemfdCreate(optionalDisplayName string, flags uint) (uintptr, error) {
	fdRaw, _, err := syscall.RawSyscall(
		memfdCreateID,
		uintptr(unsafe.Pointer(&optionalDisplayName)),
		uintptr(flags),
		0)
	if int(fdRaw) < 0 {
		return 0, fmt.Errorf(err.Error())
	}

	return fdRaw, nil
}

// CopyDataIntoMemFile is a helper function that simply copies data into the
// specified *os.File, closing data automatically on the caller's behalf.
func CopyDataIntoMemFile(data io.ReadCloser, inMemoryFile *os.File) error {
	_, err := io.Copy(inMemoryFile, data)
	data.Close()
	if err != nil {
		return fmt.Errorf("failed to copy data into in-memory file - %w", err)
	}

	return nil
}

// InMemoryFileToCmd is a helper functions that returns a *exec.Cmd
// representing an executable already loaded in memory.
//
// Refer to MemfdCreateOSFile() for more information.
func InMemoryFileToCmd(inMemory *os.File, args ...string) *exec.Cmd {
	return exec.Command(inMemory.Name(), args...)
}
