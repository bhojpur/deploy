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

	. "github.com/bhojpur/deploy/pkg/plugins"
	"github.com/bhojpur/deploy/pkg/schema"
	consoletests "github.com/bhojpur/deploy/tests/console"
	"github.com/sirupsen/logrus"
	"github.com/twpayne/go-vfs/vfst"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SSH", func() {
	Context("setting", func() {
		testConsole := consoletests.TestConsole{}
		l := logrus.New()

		It("configures a user authorized_key", func() {
			fs, cleanup, err := vfst.NewTestFS(map[string]interface{}{
				"/etc/passwd":     `foo:x:1000:100:foo:/home/foo:/bin/zsh`,
				"/home/foo/.keep": "",
			})
			Expect(err).Should(BeNil())
			defer cleanup()

			err = SSH(l, schema.Stage{
				SSHKeys: map[string][]string{"foo": {"github:bhojpur", "efafeeafea,t,t,pgl3,pbar"}},
			}, fs, testConsole)
			//Expect(err).ShouldNot(HaveOccurred())

			file, err := fs.Open("/home/foo/.ssh/authorized_keys")
			Expect(err).ShouldNot(HaveOccurred())

			b, err := ioutil.ReadAll(file)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(b)).Should(Equal("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDR9zjXvyzg1HFMC7RT4LgtR+YGstxWDPPRoAcNrAWjtQcJVrcVo4WLFnT0BMU5mtMxWSrulpC6yrwnt2TE3Ul86yMxO2hbSyGP/xOdYm/nQzufY49rd3tKeJl1+6DkczuPa+XYh1GBcW5E2laNM5ZK+RjABppMpDgmnrM3AsGNE6G8RSuUvc/6Rwt61ma+jak3F5YMj4kwr5PhY2MTPo2YshsL3ouRXP/uPsbaBM6AdQakjWGJR8tPbrnHenzF65813d9zuY4y78TG0AHfomx9btmha7Mc0YF+BpELnvSQLlYrlRY/ziGhP65aQc8lFMc+XBnHeaXF4NHnzq6dIH2D\nssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDjWfZUB5W9HU70yOD1QW/7DSYZsisg8pPHnrxzS5WFnUvhnd7x3r9i+L8mRfk0tXk9p599e5uTryqaHW74bQK360+TnVens0JRF5vGeABe2L2GGrIkTIF8aTlPVq2BTDhu0R0rU28Cw3HwywX7cNjZdpFN2MtF74QbwqB0Ue7Nj6XxJjgV7GcecKEWc23Vjie6KEHlkFcgS0objZsiSt+hY3v3wJ94t+WZ8d1vEwvp7PX2J20W8Zq0bGcJiGMGuhDPRAZ4ju6HxIm60fUo9WzMNrZKVyEbMSYo6frLcmcMN0cDpDXE9WWnCwKDKnZEB0WqQcwOh1TQLYvRYEgMJair\n\nefafeeafea,t,t,pgl3,pbar\n"))
		})
	})
})
