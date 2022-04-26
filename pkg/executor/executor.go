package executor

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
	"github.com/bhojpur/deploy/pkg/plugins"
	"github.com/sirupsen/logrus"
	"github.com/twpayne/go-vfs"

	"github.com/bhojpur/deploy/pkg/schema"
)

// Executor an executor applies a Bhojpur Deploy config
type Executor interface {
	Apply(string, schema.BhojpurConfig, vfs.FS, plugins.Console) error
	Run(string, vfs.FS, plugins.Console, ...string) error
	Plugins([]Plugin)
	Conditionals([]Plugin)
	Modifier(m schema.Modifier)
}

type Plugin func(logger.Interface, schema.Stage, vfs.FS, plugins.Console) error

type Options func(d *DefaultExecutor) error

// WithLogger sets the logger for the cloudrunner
func WithLogger(i logger.Interface) Options {
	return func(d *DefaultExecutor) error {
		d.logger = i
		return nil
	}
}

// WithPlugins sets the plugins for the cloudrunner
func WithPlugins(p ...Plugin) Options {
	return func(d *DefaultExecutor) error {
		d.plugins = p
		return nil
	}
}

// WithConditionals sets the conditionals for the cloudrunner
func WithConditionals(p ...Plugin) Options {
	return func(d *DefaultExecutor) error {
		d.conditionals = p
		return nil
	}
}

// NewExecutor returns an executor from the stringified version of it.
func NewExecutor(opts ...Options) Executor {
	d := &DefaultExecutor{
		logger: logrus.New(),
		conditionals: []Plugin{
			plugins.NodeConditional,
			plugins.IfConditional,
		},
		plugins: []Plugin{
			plugins.DNS,
			plugins.Download,
			plugins.Git,
			plugins.Entities,
			plugins.EnsureDirectories,
			plugins.EnsureFiles,
			plugins.Commands,
			plugins.DeleteEntities,
			plugins.Hostname,
			plugins.Sysctl,
			plugins.User,
			plugins.SSH,
			plugins.LoadModules,
			plugins.Timesyncd,
			plugins.Systemctl,
			plugins.Environment,
			plugins.SystemdFirstboot,
			plugins.DataSources,
			plugins.Layout,
		},
	}

	for _, o := range opts {
		o(d)
	}
	return d
}
