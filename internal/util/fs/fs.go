/*
 * Copyright (C) 2023 Stefan Kühnel
 *
 * SPDX-License-Identifier: EUPL-1.2
 */

package fs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
)

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
func CopyDir(sourcePath, destinationPath string) error {
	// Copyright (c) 2014 The Go Authors. All rights reserved.
	//
	// Redistribution and use in source and binary forms, with or without
	// modification, are permitted provided that the following conditions are
	// met:
	//
	//    * Redistributions of source code must retain the above copyright
	// notice, this list of conditions and the following disclaimer.
	//    * Redistributions in binary form must reproduce the above
	// copyright notice, this list of conditions and the following disclaimer
	// in the documentation and/or other materials provided with the
	// distribution.
	//    * Neither the name of Google Inc. nor the names of its
	// contributors may be used to endorse or promote products derived from
	// this software without specific prior written permission.
	//
	// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
	// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
	// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
	// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
	// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
	// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
	// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
	// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
	// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
	// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
	// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
	//
	// See: https://github.com/golang/dep/blob/master/internal/fs/fs.go#L356
	sourcePath = filepath.Clean(sourcePath)
	destinationPath = filepath.Clean(destinationPath)

	// os.Lstat() is used here to ensure that a loop is not encountered
	// where a symlink actually links to one of its parent directories.
	fileInfo, err := os.Lstat(sourcePath)
	if err != nil {
		return err
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("source path '%s' is not a directory", sourcePath)
	}

	_, err = os.Stat(destinationPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err = os.MkdirAll(destinationPath, fileInfo.Mode()); err != nil {
		return fmt.Errorf("cannot create directory '%s'", destinationPath)
	}

	directoryEntries, err := os.ReadDir(sourcePath)
	if err != nil {
		return fmt.Errorf("cannot read directory '%s'", sourcePath)
	}

	for _, directoryEntry := range directoryEntries {
		directoryEntrySourcePath := filepath.Join(sourcePath, directoryEntry.Name())
		directoryEntryDestinationPath := filepath.Join(destinationPath, directoryEntry.Name())

		if directoryEntry.IsDir() {
			if err = CopyDir(directoryEntrySourcePath, directoryEntryDestinationPath); err != nil {
				return fmt.Errorf("copying directory from '%s' to '%s' failed", directoryEntrySourcePath, directoryEntryDestinationPath)
			}
		} else {
			// This includes symlinks, which is desired when copying things.
			if err = CopyFile(directoryEntrySourcePath, directoryEntryDestinationPath); err != nil {
				return fmt.Errorf("copying file '%s' to '%s' failed", directoryEntrySourcePath, directoryEntryDestinationPath)
			}
		}
	}

	return nil
}

// CopyFile copies a file from the source path to the destination path.
func CopyFile(sourcePath, destinationPath string) error {
	// Copyright (c) 2014 The Go Authors. All rights reserved.
	//
	// Redistribution and use in source and binary forms, with or without
	// modification, are permitted provided that the following conditions are
	// met:
	//
	//    * Redistributions of source code must retain the above copyright
	// notice, this list of conditions and the following disclaimer.
	//    * Redistributions in binary form must reproduce the above
	// copyright notice, this list of conditions and the following disclaimer
	// in the documentation and/or other materials provided with the
	// distribution.
	//    * Neither the name of Google Inc. nor the names of its
	// contributors may be used to endorse or promote products derived from
	// this software without specific prior written permission.
	//
	// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
	// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
	// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
	// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
	// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
	// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
	// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
	// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
	// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
	// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
	// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
	//
	// See: https://github.com/golang/dep/blob/master/internal/fs/fs.go#L411
	if isSymlink, err := IsSymlink(sourcePath); err != nil {
		return fmt.Errorf("check for symlink failed for path '%s'", sourcePath)
	} else if isSymlink {
		if err := CloneSymlink(sourcePath, destinationPath); err != nil {
			if runtime.GOOS == "windows" {
				// If cloning the symlink fails on Windows because the user
				// does not have the required privileges, ignore the error and
				// fall back to copying the file contents.
				//
				// ERROR_PRIVILEGE_NOT_HELD is 1314 (0x522):
				// https://msdn.microsoft.com/en-us/library/windows/desktop/ms681385(v=vs.85).aspx
				if linkError, ok := err.(*os.LinkError); ok && linkError.Err != syscall.Errno(1314) {
					return err
				}
			} else {
				return err
			}
		} else {
			return nil
		}
	}

	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}

	defer inputFile.Close()

	outputFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}

	if _, err = io.Copy(outputFile, inputFile); err != nil {
		err := outputFile.Close()
		if err != nil {
			return err
		}
		return err
	}

	if err = outputFile.Close(); err != nil {
		return err
	}

	statInfo, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	// "os.Chmod" on Windows doesn't handle long paths (>=248).
	// Go > 1.9 implements a fix by using the "os.fixLongPath" function.
	//
	// See: https://golang.org/issue/20829
	// See: https://github.com/golang/dep/issues/774#issuecomment-311560825
	// See: https://msdn.microsoft.com/en-us/library/windows/desktop/aa365247(v=vs.85).aspx#maxpath
	err = os.Chmod(destinationPath, statInfo.Mode())
	if err != nil {
		return err
	}

	return nil
}

// MakeDir creates a directory named path, along with any
// necessary parents, or else returns an error. If path is
// already a directory, MakeDir does nothing.
func MakeDir(path string) error {
	// when using "os.MkdirAll()" an error is thrown only
	// if the path exists but is not a directory.
	err := os.MkdirAll(path, os.ModePerm)

	return err
}

// RemoveDir removes path and any children it contains.
// It removes everything it can but returns the first error
// it encounters. If there is an error, it will be of type *PathError.
func RemoveDir(path string) error {
	err := os.RemoveAll(path)

	return err
}

// ExistDir returns a boolean indicating if a directory named path exists.
// If there is an error, it will be of type *PathError.
func ExistDir(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}

// ExistFile returns a boolean indicating if a filename exists.
// If there is an error, it will be of type *PathError.
func ExistFile(filename string) bool {
	fileInfo, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !fileInfo.IsDir()
}

// StdinAll reads from os.Stdin until an error or EOF and returns the data it read.
func StdinAll() ([]byte, error) {
	return io.ReadAll(os.Stdin)
}

// IsSymlink returns a boolean indicating if a symlink exists at the path.
// If there is an error, it will be of type *PathError.
func IsSymlink(path string) (bool, error) {
	// Copyright (c) 2014 The Go Authors. All rights reserved.
	//
	// Redistribution and use in source and binary forms, with or without
	// modification, are permitted provided that the following conditions are
	// met:
	//
	//    * Redistributions of source code must retain the above copyright
	// notice, this list of conditions and the following disclaimer.
	//    * Redistributions in binary form must reproduce the above
	// copyright notice, this list of conditions and the following disclaimer
	// in the documentation and/or other materials provided with the
	// distribution.
	//    * Neither the name of Google Inc. nor the names of its
	// contributors may be used to endorse or promote products derived from
	// this software without specific prior written permission.
	//
	// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
	// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
	// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
	// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
	// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
	// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
	// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
	// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
	// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
	// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
	// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
	//
	// See: https://github.com/golang/dep/blob/master/internal/fs/fs.go#L557
	lstat, err := os.Lstat(path)
	if err != nil {
		return false, err
	}

	return lstat.Mode()&os.ModeSymlink == os.ModeSymlink, nil
}

// CloneSymlink clones a symlink from symlinkPath to destinationPath.
func CloneSymlink(symlinkPath, destinationPath string) error {
	// Copyright (c) 2014 The Go Authors. All rights reserved.
	//
	// Redistribution and use in source and binary forms, with or without
	// modification, are permitted provided that the following conditions are
	// met:
	//
	//    * Redistributions of source code must retain the above copyright
	// notice, this list of conditions and the following disclaimer.
	//    * Redistributions in binary form must reproduce the above
	// copyright notice, this list of conditions and the following disclaimer
	// in the documentation and/or other materials provided with the
	// distribution.
	//    * Neither the name of Google Inc. nor the names of its
	// contributors may be used to endorse or promote products derived from
	// this software without specific prior written permission.
	//
	// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
	// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
	// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
	// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
	// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
	// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
	// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
	// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
	// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
	// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
	// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
	//
	// See: https://github.com/golang/dep/blob/master/internal/fs/fs.go#L474
	resolvedSymlinkPath, err := os.Readlink(symlinkPath)

	if err != nil {
		return err
	}

	return os.Symlink(resolvedSymlinkPath, destinationPath)
}
