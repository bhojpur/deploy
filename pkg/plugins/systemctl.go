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
	"github.com/bhojpur/deploy/pkg/logger"
	"github.com/bhojpur/deploy/pkg/schema"
	"github.com/hashicorp/go-multierror"
	"github.com/twpayne/go-vfs"
)

func Systemctl(l logger.Interface, s schema.Stage, fs vfs.FS, console Console) error {
	var errs error

	if err := console.RunTemplate(s.Systemctl.Enable, "systemctl enable %s"); err != nil {
		errs = multierror.Append(errs, err)
	}
	if err := console.RunTemplate(s.Systemctl.Disable, "systemctl disable %s"); err != nil {
		errs = multierror.Append(errs, err)
	}
	if err := console.RunTemplate(s.Systemctl.Mask, "systemctl mask %s"); err != nil {
		errs = multierror.Append(errs, err)
	}
	if err := console.RunTemplate(s.Systemctl.Start, "systemctl start %s"); err != nil {
		errs = multierror.Append(errs, err)
	}
	return errs
}
