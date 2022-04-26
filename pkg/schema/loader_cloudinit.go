package schema

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
	"strconv"

	cloudconfig "github.com/rancher-sandbox/cloud-init/config"
	"github.com/twpayne/go-vfs"
)

type cloudInit struct{}

// Load transpiles a cloud-init style file to a Bhojpur Deploy schema.
// As Bhojpur Deploy supports multi-stages, it is encoded in the supplied one.
// fs is used to parse the user data required from /etc/passwd.
func (cloudInit) Load(s []byte, fs vfs.FS) (*BhojpurConfig, error) {
	cc, err := cloudconfig.NewCloudConfig(string(s))
	if err != nil {
		return nil, err
	}

	// Decode users and SSH Keys
	sshKeys := make(map[string][]string)
	users := make(map[string]User)
	userstoKey := []string{}

	for _, u := range cc.Users {
		userstoKey = append(userstoKey, u.Name)
		users[u.Name] = User{
			Name:         u.Name,
			PasswordHash: u.PasswordHash,
			GECOS:        u.GECOS,
			Homedir:      u.Homedir,
			NoCreateHome: u.NoCreateHome,
			PrimaryGroup: u.PrimaryGroup,
			Groups:       u.Groups,
			NoUserGroup:  u.NoUserGroup,
			System:       u.System,
			NoLogInit:    u.NoLogInit,
			Shell:        u.Shell,
			UID:          u.UID,
			LockPasswd:   u.LockPasswd,
		}
		sshKeys[u.Name] = u.SSHAuthorizedKeys
	}

	for _, uu := range userstoKey {
		_, exists := sshKeys[uu]
		if !exists {
			sshKeys[uu] = cc.SSHAuthorizedKeys
		} else {
			sshKeys[uu] = append(sshKeys[uu], cc.SSHAuthorizedKeys...)
		}
	}

	// If no users are defined, then assume global ssh_authorized_keys is assigned to root
	if len(userstoKey) == 0 && len(cc.SSHAuthorizedKeys) > 0 {
		sshKeys["root"] = cc.SSHAuthorizedKeys
	}

	// Decode writeFiles
	var f []File
	for _, ff := range append(cc.WriteFiles, cc.MilpaFiles...) {
		newFile := File{
			Path:        ff.Path,
			OwnerString: ff.Owner,
			Content:     ff.Content,
			Encoding:    ff.Encoding,
		}
		newFile.Permissions, err = parseOctal(ff.RawFilePermissions)
		if err != nil {
			return nil, fmt.Errorf("converting permission %s for %s: %w", ff.RawFilePermissions, ff.Path, err)
		}
		f = append(f, newFile)
	}

	stages := []Stage{{
		Commands: cc.RunCmd,
		Files:    f,
		Users:    users,
		SSHKeys:  sshKeys,
	}}

	for _, d := range cc.Partitioning.Devices {
		layout := &Layout{}
		layout.Expand = &Expand{Size: 0}
		layout.Device = &Device{Path: d}
		stages = append(stages, Stage{Layout: *layout})
	}

	result := &BhojpurConfig{
		Name: "Cloud init",
		Stages: map[string][]Stage{
			"boot": stages,
			"initramfs": {{
				Hostname: cc.Hostname,
			}},
		},
	}

	// optimistically load data as Bhojpur Deploy yaml
	bhojpurConfig, err := bhojpurYAML{}.Load(s, fs)
	if err == nil {
		for k, v := range bhojpurConfig.Stages {
			result.Stages[k] = append(result.Stages[k], v...)
		}
	}

	return result, nil
}

func parseOctal(srv string) (uint32, error) {
	if srv == "" {
		return 0, nil
	}
	i, err := strconv.ParseUint(srv, 8, 32)
	if err != nil {
		return 0, err
	}
	return uint32(i), nil
}
