package plugins

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bhojpur/deploy/pkg/logger"
	"github.com/bhojpur/deploy/pkg/schema"
	"github.com/hashicorp/go-multierror"
	"github.com/twpayne/go-vfs"
)

func EnsureDirectories(l logger.Interface, s schema.Stage, fs vfs.FS, console Console) error {
	var errs error
	for _, dir := range s.Directories {
		if err := writePath(l, dir, fs, true); err != nil {
			l.Error(err.Error())
			errs = multierror.Append(errs, err)
			continue
		}
	}
	return errs
}

func writeDirectory(l logger.Interface, dir schema.Directory, fs vfs.FS) error {
	l.Debug("Creating directory ", dir.Path)
	err := fs.Mkdir(dir.Path, os.FileMode(dir.Permissions))
	if err != nil {
		return err
	}

	return fs.Chown(dir.Path, dir.Owner, dir.Group)
}

func writePath(l logger.Interface, dir schema.Directory, fs vfs.FS, topLevel bool) error {
	inf, err := fs.Stat(dir.Path)
	if err == nil && inf.IsDir() && topLevel {
		// The path already exists, apply permissions and ownership only
		err = fs.Chmod(dir.Path, os.FileMode(dir.Permissions))
		if err != nil {
			return err
		}
		return fs.Chown(dir.Path, dir.Owner, dir.Group)
	} else if err == nil && !inf.IsDir() {
		return fmt.Errorf("Error, '%s' already exists and it is not a directory", dir.Path)
	} else if err == nil {
		return nil
	} else {
		parentDir := filepath.Dir(dir.Path)
		_, err = fs.Stat(parentDir)
		if parentDir == "/" || parentDir == "." || err == nil {
			//There is no parent dir or it already exists
			return writeDirectory(l, dir, fs)
		} else {
			//Parent dir needs to be created
			pDir := schema.Directory{parentDir, dir.Permissions, dir.Owner, dir.Group}
			err = writePath(l, pDir, fs, false)
			if err != nil {
				return err
			}
			return writeDirectory(l, dir, fs)
		}
	}
}
