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
	"container/list"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/apex/log"
	"github.com/hashicorp/go-multierror"
	. "github.com/onsi/gomega"
)

type CmdMock struct {
	Cmd       string
	Output    string
	UseRegexp bool
}

type TestConsoleMock struct {
	Cmds *list.List
}

func New() *TestConsoleMock {
	return &TestConsoleMock{Cmds: list.New()}
}

func (s TestConsoleMock) AddCmd(cmd CmdMock) {
	s.Cmds.PushBack(&cmd)
}

func (s TestConsoleMock) AddCmds(cmds []CmdMock) {
	for _, cmd := range cmds {
		s.AddCmd(cmd)
	}
}

func (s TestConsoleMock) PopCmd() *CmdMock {
	e := s.Cmds.Front()
	if e == nil {
		return nil
	}
	s.Cmds.Remove(e)
	cmdMock := e.Value.(*CmdMock)
	return cmdMock
}

func (s TestConsoleMock) Run(cmd string, opts ...func(*exec.Cmd)) (string, error) {
	cmdMock := s.PopCmd()
	Expect(cmdMock).NotTo(BeNil())
	if cmdMock.UseRegexp {
		if matched, _ := regexp.MatchString(cmdMock.Cmd, cmd); matched {
			return cmdMock.Output, nil
		}
	} else {
		if cmdMock.Cmd == cmd {
			return cmdMock.Output, nil
		}
	}

	Expect(cmd).To(Equal(cmdMock.Cmd))
	return "", errors.New("Unexpected command")
}

func (s TestConsoleMock) Start(cmd *exec.Cmd, opts ...func(*exec.Cmd)) error {
	cmdMock := s.PopCmd()
	Expect(cmdMock).NotTo(BeNil())
	cmdStr := strings.Join(cmd.Args[:], " ")
	if cmdMock.Cmd == cmdStr {
		return nil
	} else {
		Expect(cmdStr).To(Equal(cmdMock.Cmd))
		return errors.New("Unexpected command")
	}
}

func (s TestConsoleMock) RunTemplate(st []string, template string) error {
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
