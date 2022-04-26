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
	"io/ioutil"
	"os"

	. "github.com/bhojpur/deploy/pkg/plugins"
	"github.com/bhojpur/deploy/pkg/schema"
	consoletests "github.com/bhojpur/deploy/tests/console"
	"github.com/sirupsen/logrus"
	"github.com/twpayne/go-vfs/vfst"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Environment", func() {
	Context("setting", func() {
		testConsole := consoletests.TestConsole{}
		l := logrus.New()
		It("configures a /etc/environment setting", func() {
			fs, cleanup, err := vfst.NewTestFS(map[string]interface{}{"/etc/environment": ""})
			Expect(err).Should(BeNil())
			defer cleanup()

			err = Environment(l, schema.Stage{
				Environment: map[string]string{"foo": "0"},
			}, fs, testConsole)
			Expect(err).ShouldNot(HaveOccurred())

			file, err := fs.Open("/etc/environment")
			Expect(err).ShouldNot(HaveOccurred())

			b, err := ioutil.ReadAll(file)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(b)).Should(Equal("foo=\"0\""))
		})
		It("configures a /run/cos/cos-layout.env file and creates missing directories", func() {
			fs, cleanup, err := vfst.NewTestFS(map[string]interface{}{"/run": &vfst.Dir{Perm: 0o755}})
			Expect(err).Should(BeNil())
			defer cleanup()

			_, err = fs.Stat("/run/cos")
			Expect(err).NotTo(BeNil())

			err = Environment(l, schema.Stage{
				Environment:     map[string]string{"foo": "0"},
				EnvironmentFile: "/run/cos/cos-layout.env",
			}, fs, testConsole)
			Expect(err).ShouldNot(HaveOccurred())

			inf, _ := fs.Stat("/run/cos")
			Expect(inf.Mode().Perm()).To(Equal(os.FileMode(int(0744))))

			file, err := fs.Open("/run/cos/cos-layout.env")
			Expect(err).ShouldNot(HaveOccurred())

			b, err := ioutil.ReadAll(file)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(b)).Should(Equal("foo=\"0\""))
		})
	})
})
