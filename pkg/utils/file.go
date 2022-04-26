package utils

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
	"errors"
	"os"

	"github.com/twpayne/go-vfs"
)

func Touch(s string, perms os.FileMode, fs vfs.FS) error {
	_, err := fs.Stat(s)

	switch {
	case os.IsNotExist(err):
		f, err := fs.Create(s)
		if err != nil {
			return err
		}
		if err = f.Chmod(perms); err != nil {
			return err
		}
		if err = f.Close(); err != nil {
			return err
		}
		_, err = fs.Stat(s)
		return err
	case err == nil:
		return nil
	default:
		return errors.New("could not create file")
	}

}

func Exists(s string) bool {
	if _, err := os.Stat(s); err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
