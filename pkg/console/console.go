package console

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
	"os/exec"

	"github.com/bhojpur/deploy/pkg/logger"
	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
)

type StandardConsole struct {
	logger logger.Interface
}

type StandardConsoleOptions func(*StandardConsole) error

func WithLogger(i logger.Interface) StandardConsoleOptions {
	return func(sc *StandardConsole) error {
		sc.logger = i
		return nil
	}
}

func NewStandardConsole(opts ...StandardConsoleOptions) *StandardConsole {
	c := &StandardConsole{
		logger: logrus.New(),
	}
	for _, o := range opts {
		o(c)
	}
	return c

}

func (s StandardConsole) Run(cmd string, opts ...func(cmd *exec.Cmd)) (string, error) {
	s.logger.Debugf("running command `%s`", cmd)
	c := exec.Command("sh", "-c", cmd)
	for _, o := range opts {
		o(c)
	}
	out, err := c.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("failed to run %s: %v", cmd, err)
	}

	return string(out), err
}

func (s StandardConsole) Start(cmd *exec.Cmd, opts ...func(cmd *exec.Cmd)) error {
	s.logger.Debugf("running command `%s`", cmd)
	for _, o := range opts {
		o(cmd)
	}
	return cmd.Run()
}

func (s StandardConsole) RunTemplate(st []string, template string) error {
	var errs error

	for _, svc := range st {
		out, err := s.Run(fmt.Sprintf(template, svc))
		if err != nil {
			s.logger.Error(out)
			s.logger.Error(err.Error())
			errs = multierror.Append(errs, err)
			continue
		}
	}
	return errs
}
