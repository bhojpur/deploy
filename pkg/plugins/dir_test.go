package plugins_test

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

	. "github.com/bhojpur/deploy/pkg/plugins"
	"github.com/bhojpur/deploy/pkg/schema"
	consoletests "github.com/bhojpur/deploy/tests/console"
	"github.com/sirupsen/logrus"
	"github.com/twpayne/go-vfs/vfst"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Files", func() {
	Context("creating", func() {
		testConsole := consoletests.TestConsole{}
		l := logrus.New()
		It("Creates a /tmp/dir directory", func() {
			fs, cleanup, err := vfst.NewTestFS(map[string]interface{}{"/tmp": &vfst.Dir{Perm: 0o755}})
			Expect(err).Should(BeNil())
			defer cleanup()

			err = EnsureDirectories(l, schema.Stage{
				Directories: []schema.Directory{{Path: "/tmp/dir", Permissions: 0740, Owner: os.Getuid(), Group: os.Getgid()}},
			}, fs, testConsole)
			Expect(err).ShouldNot(HaveOccurred())
			inf, _ := fs.Stat("/tmp/dir")
			Expect(inf.Mode().Perm()).To(Equal(os.FileMode(int(0740))))
		})

		It("Changes permissions of existing directory /tmp/dir directory", func() {
			fs, cleanup, err := vfst.NewTestFS(map[string]interface{}{"/tmp/dir": &vfst.Dir{Perm: 0o755}})
			Expect(err).Should(BeNil())
			defer cleanup()
			inf, _ := fs.Stat("/tmp/dir")
			Expect(inf.Mode().Perm()).To(Equal(os.FileMode(int(0755))))
			err = EnsureDirectories(l, schema.Stage{
				Directories: []schema.Directory{{Path: "/tmp/dir", Permissions: 0740, Owner: os.Getuid(), Group: os.Getgid()}},
			}, fs, testConsole)
			Expect(err).ShouldNot(HaveOccurred())
			inf, _ = fs.Stat("/tmp/dir")
			Expect(inf.Mode().Perm()).To(Equal(os.FileMode(int(0740))))
		})

		It("Creates /tmp/dir/subdir1/subdir2 directory and its missing parent dirs", func() {
			fs, cleanup, err := vfst.NewTestFS(map[string]interface{}{"/tmp": &vfst.Dir{Perm: 0o755}})
			Expect(err).Should(BeNil())
			defer cleanup()
			err = EnsureDirectories(l, schema.Stage{
				Directories: []schema.Directory{{Path: "/tmp/dir/subdir1/subdir2", Permissions: 0740, Owner: os.Getuid(), Group: os.Getgid()}},
			}, fs, testConsole)
			Expect(err).ShouldNot(HaveOccurred())
			inf, _ := fs.Stat("/tmp")
			Expect(inf.Mode().Perm()).To(Equal(os.FileMode(int(0755))))
			inf, _ = fs.Stat("/tmp/dir/subdir1/subdir2")
			Expect(inf.Mode().Perm()).To(Equal(os.FileMode(int(0740))))
		})
	})
})
