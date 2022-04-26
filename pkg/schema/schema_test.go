package schema_test

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
	. "github.com/bhojpur/deploy/pkg/schema"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/twpayne/go-vfs/vfst"
)

func loadstdBhojpur(s string) *BhojpurConfig {
	fs, cleanup, err := vfst.NewTestFS(map[string]interface{}{"/deploy.yaml": s, "/etc/passwd": ""})
	Expect(err).Should(BeNil())
	defer cleanup()

	bhojpurConfig, err := Load("/deploy.yaml", fs, FromFile, nil)
	Expect(err).ToNot(HaveOccurred())
	return bhojpurConfig
}

func loadBhojpur(s string) *BhojpurConfig {
	fs, cleanup, err := vfst.NewTestFS(map[string]interface{}{"/deploy.yaml": s})
	Expect(err).Should(BeNil())
	defer cleanup()

	bhojpurConfig, err := Load("/deploy.yaml", fs, FromFile, DotNotationModifier)
	Expect(err).ToNot(HaveOccurred())
	return bhojpurConfig
}

var _ = Describe("Schema", func() {
	Context("Loading from dot notation", func() {
		oneConfigwithGarbageS := "stages.foo[0].name=bar boo.baz"
		twoConfigsS := "stages.foo[0].name=bar   stages.foo[0].commands[0]=baz"
		threeConfigInvalid := `ip=dhcp test="echo ping_test_host=127.0.0.1  > /tmp/jojo"`
		fourConfigHalfInvalid := `stages.foo[0].name=bar ip=dhcp test="echo ping_test_host=127.0.0.1  > /tmp/dio"`

		It("Reads deploy file correctly", func() {
			bhojpurConfig := loadBhojpur(oneConfigwithGarbageS)
			Expect(bhojpurConfig.Stages["foo"][0].Name).To(Equal("bar"))
		})
		It("Reads deploy file correctly", func() {
			bhojpurConfig := loadBhojpur(twoConfigsS)
			Expect(bhojpurConfig.Stages["foo"][0].Name).To(Equal("bar"))
			Expect(bhojpurConfig.Stages["foo"][0].Commands[0]).To(Equal("baz"))
		})

		It("Reads deploy file correctly", func() {
			bhojpurConfig, err := Load(twoConfigsS, nil, nil, DotNotationModifier)
			Expect(err).ToNot(HaveOccurred())
			Expect(bhojpurConfig.Stages["foo"][0].Name).To(Equal("bar"))
			Expect(bhojpurConfig.Stages["foo"][0].Commands[0]).To(Equal("baz"))
		})

		It("Reads bhojpur file correctly", func() {
			bhojpurConfig, err := Load(threeConfigInvalid, nil, nil, DotNotationModifier)
			Expect(err).ToNot(HaveOccurred())
			// should look like an empty Bhojpur Deploy Config as its an invalid config, so nothing should be loaded
			Expect(bhojpurConfig.Stages).To(Equal(BhojpurConfig{}.Stages))
			Expect(bhojpurConfig.Name).To(Equal(BhojpurConfig{}.Name))
		})

		It("Reads deploy file correctly", func() {
			bhojpurConfig, err := Load(fourConfigHalfInvalid, nil, nil, DotNotationModifier)
			Expect(err).ToNot(HaveOccurred())
			Expect(bhojpurConfig.Name).To(Equal(BhojpurConfig{}.Name))
			// Even if broken config, it should load the valid parts of the config
			Expect(bhojpurConfig.Stages["foo"][0].Name).To(Equal("bar"))
		})
	})

	Context("Loading CloudConfig", func() {
		It("Reads cloudconfig to boot stage", func() {
			bhojpurConfig := loadstdBhojpur(`#cloud-config
growpart:
 devices: ['/']
stages:
  test:
  - environment:
      foo: bar
users:
- name: "bar"
  passwd: "foo"
  uid: "1002"
  lock_passwd: true
  groups: "users"
  ssh_authorized_keys:
  - faaapploo
ssh_authorized_keys:
  - asdd
runcmd:
- foo
hostname: "bar"
write_files:
- encoding: b64
  content: CiMgVGhpcyBmaWxlIGNvbnRyb2xzIHRoZSBzdGF0ZSBvZiBTRUxpbnV4
  path: /foo/bar
  permissions: "0644"
  owner: "bar"
`)
			Expect(len(bhojpurConfig.Stages)).To(Equal(3))
			Expect(bhojpurConfig.Stages["boot"][0].Users["bar"].UID).To(Equal("1002"))
			Expect(bhojpurConfig.Stages["boot"][0].Users["bar"].PasswordHash).To(Equal("foo"))
			Expect(bhojpurConfig.Stages["boot"][0].SSHKeys).To(Equal(map[string][]string{"bar": {"faaapploo", "asdd"}}))
			Expect(bhojpurConfig.Stages["boot"][0].Files[0].Path).To(Equal("/foo/bar"))
			Expect(bhojpurConfig.Stages["boot"][0].Files[0].Permissions).To(Equal(uint32(0644)))
			Expect(bhojpurConfig.Stages["boot"][0].Hostname).To(Equal(""))
			Expect(bhojpurConfig.Stages["initramfs"][0].Hostname).To(Equal("bar"))
			Expect(bhojpurConfig.Stages["boot"][0].Commands).To(Equal([]string{"foo"}))
			Expect(bhojpurConfig.Stages["test"][0].Environment["foo"]).To(Equal("bar"))
			Expect(bhojpurConfig.Stages["boot"][0].Users["bar"].LockPasswd).To(Equal(true))
			Expect(bhojpurConfig.Stages["boot"][1].Layout.Expand.Size).To(Equal(uint(0)))
			Expect(bhojpurConfig.Stages["boot"][1].Layout.Device.Path).To(Equal("/"))
		})
	})
})
