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
	"os"
	"path/filepath"

	"github.com/bhojpur/deploy/pkg/logger"
	"github.com/bhojpur/deploy/pkg/schema"
	"github.com/bhojpur/deploy/pkg/utils"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/twpayne/go-vfs"
)

func EnsureFiles(l logger.Interface, s schema.Stage, fs vfs.FS, console Console) error {
	var errs error
	for _, file := range s.Files {
		if err := writeFile(l, file, fs, console); err != nil {
			l.Error(err.Error())
			errs = multierror.Append(errs, err)
			continue
		}
	}
	return errs
}

func writeFile(l logger.Interface, file schema.File, fs vfs.FS, console Console) error {
	l.Debug("Creating file ", file.Path)
	parentDir := filepath.Dir(file.Path)
	_, err := fs.Stat(parentDir)
	if err != nil {
		l.Debug("Creating parent directories")
		perm := file.Permissions
		if perm < 0700 {
			l.Debug("Adding execution bit to parent directory")
			perm = perm + 0100
		}
		if err = EnsureDirectories(l, schema.Stage{
			Directories: []schema.Directory{
				{
					Path:        parentDir,
					Permissions: perm,
					Owner:       file.Owner,
					Group:       file.Group,
				},
			},
		}, fs, console); err != nil {
			l.Infof("Failed to write %s: %s", parentDir, err)
			return err
		}
	}
	fsfile, err := fs.Create(file.Path)
	if err != nil {
		return err
	}
	defer fsfile.Close()

	d := newDecoder(file.Encoding)
	c, err := d.Decode(file.Content)
	if err != nil {
		return errors.Wrapf(err, "failed decoding content with encoding %s", file.Encoding)
	}

	_, err = fsfile.WriteString(templateSysData(l, string(c)))
	if err != nil {
		return err

	}
	err = fs.Chmod(file.Path, os.FileMode(file.Permissions))
	if err != nil {
		return err

	}

	if file.OwnerString != "" {
		// FIXUP: Doesn't support fs. It reads real /etc/passwd files
		uid, gid, err := utils.GetUserDataFromString(file.OwnerString)
		if err != nil {
			return errors.Wrap(err, "Failed getting gid")
		}
		return fs.Chown(file.Path, uid, gid)
	}

	return fs.Chown(file.Path, file.Owner, file.Group)
}
