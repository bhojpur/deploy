package entities

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

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const (
	UserKind    = "user"
	ShadowKind  = "shadow"
	GroupKind   = "group"
	GShadowKind = "gshadow"
)

type EntitiesParser interface {
	ReadEntity(entity string) (Entity, error)
}

type Signature struct {
	Kind string `yaml:"kind"`
}

type Parser struct{}

func (p Parser) ReadEntityFromBytes(yamlFile []byte) (Entity, error) {

	var signature Signature
	err := yaml.Unmarshal(yamlFile, &signature)
	if err != nil {
		return nil, errors.Wrap(err, "Failed while parsing entity file")
	}

	switch signature.Kind {
	case UserKind:
		var user UserPasswd

		err = yaml.Unmarshal(yamlFile, &user)
		if err != nil {
			return nil, errors.Wrap(err, "Failed while parsing entity file")
		}
		return user, nil
	case ShadowKind:
		var shad Shadow

		err = yaml.Unmarshal(yamlFile, &shad)
		if err != nil {
			return nil, errors.Wrap(err, "Failed while parsing entity file")
		}
		return shad, nil
	case GroupKind:
		var group Group

		err = yaml.Unmarshal(yamlFile, &group)
		if err != nil {
			return nil, errors.Wrap(err, "Failed while parsing entity file")
		}
		return group, nil

	case GShadowKind:
		var group GShadow

		err = yaml.Unmarshal(yamlFile, &group)
		if err != nil {
			return nil, errors.Wrap(err, "Failed while parsing entity file")
		}
		return group, nil
	}

	return nil, errors.New("Unsupported format")
}
func (p Parser) ReadEntity(entity string) (Entity, error) {
	yamlFile, err := ioutil.ReadFile(entity)
	if err != nil {
		return nil, errors.Wrap(err, "Failed while reading entity file")
	}
	return p.ReadEntityFromBytes(yamlFile)

}
