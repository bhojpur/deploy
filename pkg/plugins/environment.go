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
	"path/filepath"

	"github.com/bhojpur/deploy/pkg/logger"
	"github.com/bhojpur/deploy/pkg/schema"
	"github.com/bhojpur/deploy/pkg/utils"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/twpayne/go-vfs"
)

const environmentFile = "/etc/environment"
const envFilePerm uint32 = 0644

func Environment(l logger.Interface, s schema.Stage, fs vfs.FS, console Console) error {
	if len(s.Environment) == 0 {
		return nil
	}
	environment := s.EnvironmentFile
	if environment == "" {
		environment = environmentFile
	}

	parentDir := filepath.Dir(environment)
	_, err := fs.Stat(parentDir)
	if err != nil {
		perm := envFilePerm
		if perm < 0700 {
			perm = perm + 0100
		}
		if err = EnsureDirectories(l, schema.Stage{
			Directories: []schema.Directory{
				{
					Path:        parentDir,
					Permissions: perm,
					Owner:       os.Getuid(),
					Group:       os.Getgid(),
				},
			},
		}, fs, console); err != nil {
			return err
		}
	}

	if err := utils.Touch(environment, os.ModePerm, fs); err != nil {
		return errors.Wrap(err, "failed touching environment file")
	}

	content, err := fs.ReadFile(environment)
	if err != nil {
		return err
	}

	env, _ := godotenv.Unmarshal(string(content))
	for key, val := range s.Environment {
		env[key] = templateSysData(l, val)
	}

	p, err := fs.RawPath(environment)
	if err != nil {
		return err
	}
	err = godotenv.Write(env, p)
	if err != nil {
		return err
	}

	return fs.Chmod(environment, os.FileMode(envFilePerm))
}
