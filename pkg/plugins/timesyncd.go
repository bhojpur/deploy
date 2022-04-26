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

	"github.com/bhojpur/deploy/pkg/logger"
	"github.com/bhojpur/deploy/pkg/schema"
	"github.com/twpayne/go-vfs"
	"gopkg.in/ini.v1"
)

const timeSyncd = "/etc/systemd/timesyncd.conf"

func Timesyncd(l logger.Interface, s schema.Stage, fs vfs.FS, console Console) error {
	if len(s.TimeSyncd) == 0 {
		return nil
	}
	var errs error

	path, err := fs.RawPath(timeSyncd)
	if err != nil {
		return err
	}

	if _, err := fs.Stat(timeSyncd); os.IsNotExist(err) {
		f, _ := fs.Create(timeSyncd)
		f.Close()
	}

	cfg, err := ini.Load(path)
	if err != nil {
		return err
	}

	for k, v := range s.TimeSyncd {
		cfg.Section("Time").Key(k).SetValue(v)
	}

	cfg.SaveTo(path)

	return errs
}
