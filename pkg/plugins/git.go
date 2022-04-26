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
	"path/filepath"

	"github.com/bhojpur/deploy/pkg/logger"
	"github.com/bhojpur/deploy/pkg/schema"
	"github.com/bhojpur/deploy/pkg/utils"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	gith "github.com/go-git/go-git/v5/plumbing/transport/http"
	ssh2 "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/pkg/errors"
	"github.com/twpayne/go-vfs"
	"golang.org/x/crypto/ssh"
)

func Git(l logger.Interface, s schema.Stage, fs vfs.FS, console Console) error {
	if s.Git.URL == "" {
		return nil
	}

	branch := "master"
	if s.Git.Branch != "" {
		branch = s.Git.Branch
	}

	gitconfig := s.Git
	path, err := fs.RawPath(s.Git.Path)
	if err != nil {
		return err
	}
	l.Infof("Cloning git repository '%s'", s.Git.URL)

	if utils.Exists(filepath.Join(path, ".git")) {
		l.Info("Repository already exists, updating it")
		// is a git repo, update it
		// We instantiate a new repository targeting the given path (the .git folder)
		r, err := git.PlainOpen(path)
		if err != nil {
			return err
		}

		w, err := r.Worktree()
		if err != nil {
			return err
		}

		err = w.Pull(&git.PullOptions{
			Auth:            authMethod(s),
			SingleBranch:    s.Git.BranchOnly,
			Force:           true,
			InsecureSkipTLS: s.Git.Auth.Insecure,
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return err
		}

		err = w.Reset(&git.ResetOptions{
			Commit: plumbing.NewHash(branch),
			Mode:   git.HardReset,
		})

		if err != nil {
			return err
		}
		return nil

	}

	opts := &git.CloneOptions{
		URL:          gitconfig.URL,
		SingleBranch: s.Git.BranchOnly,
	}

	applyOptions(s, opts)

	_, err = git.PlainClone(path, false, opts)
	if err != nil {
		return errors.Wrap(err, "failed cloning repo")
	}
	return nil
}

func authMethod(s schema.Stage) transport.AuthMethod {
	var t transport.AuthMethod

	if s.Git.Auth.Username != "" {
		t = &gith.BasicAuth{Username: s.Git.Auth.Username, Password: s.Git.Auth.Password}
	}

	if s.Git.Auth.PrivateKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(s.Git.Auth.PrivateKey))
		if err != nil {
			return t
		}

		userName := "git"
		if s.Git.Auth.Username != "" {
			userName = s.Git.Auth.Username
		}
		sshAuth := &ssh2.PublicKeys{
			User:   userName,
			Signer: signer,
		}
		if s.Git.Auth.Insecure {
			sshAuth.HostKeyCallbackHelper = ssh2.HostKeyCallbackHelper{
				HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			}
		}
		if s.Git.Auth.PublicKey != "" {
			key, err := ssh.ParsePublicKey([]byte(s.Git.Auth.PublicKey))
			if err != nil {
				return t
			}
			sshAuth.HostKeyCallbackHelper = ssh2.HostKeyCallbackHelper{
				HostKeyCallback: ssh.FixedHostKey(key),
			}
		}

		t = sshAuth
	}
	return t
}

func applyOptions(s schema.Stage, g *git.CloneOptions) {

	g.Auth = authMethod(s)

	if s.Git.Branch != "" {
		g.ReferenceName = plumbing.NewBranchReferenceName(s.Git.Branch)
	}
	if s.Git.BranchOnly {
		g.SingleBranch = true
	}
}
