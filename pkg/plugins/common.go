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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/bhojpur/deploy/pkg/logger"
	"github.com/bhojpur/deploy/pkg/utils"
	"github.com/pkg/errors"
	"github.com/zcalusic/sysinfo"
)

var system sysinfo.SysInfo

func init() {
	system.GetSysInfo()
}

type Console interface {
	Run(string, ...func(*exec.Cmd)) (string, error)
	Start(*exec.Cmd, ...func(*exec.Cmd)) error
	RunTemplate([]string, string) error
}

func templateSysData(l logger.Interface, s string) string {
	interpolateOpts := map[string]interface{}{}

	data, err := json.Marshal(&system)
	if err != nil {
		l.Warn(fmt.Sprintf("Failed marshalling '%s': %s", s, err.Error()))
		return s
	}
	l.Debug(string(data))

	err = json.Unmarshal(data, &interpolateOpts)
	if err != nil {
		l.Warn(fmt.Sprintf("Failed marshalling '%s': %s", s, err.Error()))
		return s
	}

	rendered, err := utils.TemplatedString(s, map[string]interface{}{"Values": interpolateOpts})
	if err != nil {
		l.Warn(fmt.Sprintf("Failed rendering '%s': %s", s, err.Error()))
		return s
	}
	return rendered
}

func download(url string) (string, error) {
	var resp *http.Response
	var err error
	for i := 0; i < 10; i++ {
		resp, err = http.Get(url)
		if err == nil || strings.Contains(err.Error(), "unsupported protocol scheme") {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		return "", errors.Wrap(err, "failed while getting file")
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode/100 > 2 {
		return "", fmt.Errorf("%s %s", resp.Proto, resp.Status)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	return string(bytes), err
}
