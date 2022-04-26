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
	"bufio"
	"strings"

	"github.com/bhojpur/deploy/pkg/logger"
	"github.com/bhojpur/deploy/pkg/schema"
	"github.com/hashicorp/go-multierror"
	"github.com/twpayne/go-vfs"
	"pault.ag/go/modprobe"
)

const (
	modules = "/proc/modules"
)

func loadedModules(l logger.Interface, fs vfs.FS) map[string]interface{} {
	loaded := map[string]interface{}{}
	f, err := fs.Open(modules)
	if err != nil {
		l.Warnf("Cannot open %s: %s", modules, err.Error())
		return loaded
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		mod := strings.SplitN(sc.Text(), " ", 2)
		if len(mod) == 0 {
			continue
		}
		loaded[mod[0]] = nil
	}
	return loaded
}

func LoadModules(l logger.Interface, s schema.Stage, fs vfs.FS, console Console) error {
	var errs error

	if len(s.Modules) == 0 {
		return nil
	}

	loaded := loadedModules(l, fs)

	for _, m := range s.Modules {
		if _, ok := loaded[m]; ok {
			continue
		}
		params := strings.SplitN(m, " ", -1)
		l.Debugf("loading module %s with parameters [%s]", m, params)
		if err := modprobe.Load(params[0], strings.Join(params[1:], " ")); err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
		l.Debugf("module %s loaded", m)
	}
	return errs
}
