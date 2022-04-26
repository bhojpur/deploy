package consoletests

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
	"io/ioutil"
	"os/exec"

	"github.com/apex/log"
	"github.com/hashicorp/go-multierror"
)

var Commands []string
var Stdin string

type TestConsole struct {
}

func (s TestConsole) Run(cmd string, opts ...func(*exec.Cmd)) (string, error) {
	c := &exec.Cmd{}
	for _, o := range opts {
		o(c)
	}
	Commands = append(Commands, cmd)
	Commands = append(Commands, c.Args...)
	if c.Stdin != nil {
		b, _ := ioutil.ReadAll(c.Stdin)
		Stdin = string(b)
	}

	return "", nil
}

func Reset() {
	Commands = []string{}
	Stdin = ""
}
func (s TestConsole) Start(cmd *exec.Cmd, opts ...func(*exec.Cmd)) error {
	for _, o := range opts {
		o(cmd)
	}
	Commands = append(Commands, cmd.Args...)
	if cmd.Stdin != nil {
		b, _ := ioutil.ReadAll(cmd.Stdin)
		Stdin = string(b)
	}
	return nil
}

func (s TestConsole) RunTemplate(st []string, template string) error {
	var errs error

	for _, svc := range st {
		out, err := s.Run(fmt.Sprintf(template, svc))
		if err != nil {
			log.Error(out)
			log.Error(err.Error())
			errs = multierror.Append(errs, err)
			continue
		}
	}
	return errs
}
