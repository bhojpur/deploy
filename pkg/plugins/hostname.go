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
	"fmt"
	"math/rand"
	"strings"
	"syscall"
	"time"

	"github.com/bhojpur/deploy/pkg/logger"
	"github.com/bhojpur/deploy/pkg/schema"
	"github.com/bhojpur/deploy/pkg/utils"
	"github.com/denisbrodbeck/machineid"
	"github.com/hashicorp/go-multierror"
	uuid "github.com/satori/go.uuid"
	"github.com/twpayne/go-vfs"
)

const localHost = "127.0.0.1"

func Hostname(l logger.Interface, s schema.Stage, fs vfs.FS, console Console) error {
	var errs error
	hostname := s.Hostname
	if hostname == "" {
		return nil
	}

	// Template the input string with random generated strings and UUID.
	// Those can be used to e.g. generate random node names based on patterns "foo-{{.UUID}}"
	rand.Seed(time.Now().UnixNano())

	id, _ := machineid.ID()
	myuuid := uuid.NewV4()
	tmpl, err := utils.TemplatedString(hostname,
		struct {
			UUID      string
			Random    string
			MachineID string
		}{
			UUID:      myuuid.String(),
			MachineID: id,
			Random:    utils.RandomString(32),
		},
	)
	if err != nil {
		return err
	}

	if err := syscall.Sethostname([]byte(tmpl)); err != nil {
		errs = multierror.Append(errs, err)
	}
	if err := SystemHostname(tmpl, fs); err != nil {
		errs = multierror.Append(errs, err)
	}
	if err := UpdateHostsFile(tmpl, fs); err != nil {
		errs = multierror.Append(errs, err)
	}
	return errs
}

func UpdateHostsFile(hostname string, fs vfs.FS) error {
	hosts, err := fs.Open("/etc/hosts")
	if err != nil {
		return err
	}
	defer hosts.Close()

	lines := bufio.NewScanner(hosts)
	content := ""
	for lines.Scan() {
		line := strings.TrimSpace(lines.Text())
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == localHost {
			content += fmt.Sprintf("%s localhost %s\n", localHost, hostname)
			continue
		}
		content += line + "\n"
	}
	return fs.WriteFile("/etc/hosts", []byte(content), 0600)
}

func SystemHostname(hostname string, fs vfs.FS) error {
	return fs.WriteFile("/etc/hostname", []byte(hostname+"\n"), 0644)
}
