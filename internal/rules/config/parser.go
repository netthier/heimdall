// Copyright 2022 Dimitrij Drus <dadrus@gmx.de>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"errors"
	"io"

	"gopkg.in/yaml.v3"

	"github.com/dadrus/heimdall/internal/heimdall"
	"github.com/dadrus/heimdall/internal/x/errorchain"
)

var ErrEmptyRuleSet = errors.New("empty rule set")

func ParseRules(contentType string, reader io.Reader) (*RuleSet, error) {
	switch contentType {
	case "application/json":
		fallthrough
	case "application/yaml":
		return parseYAML(reader)
	default:
		// check if the contents are empty. in that case nothing needs to be decoded anyway
		b := make([]byte, 1)
		if _, err := reader.Read(b); err != nil && errors.Is(err, io.EOF) {
			return nil, ErrEmptyRuleSet
		}

		// otherwise
		return nil, errorchain.NewWithMessagef(heimdall.ErrInternal,
			"unsupported '%s' content type", contentType)
	}
}

func parseYAML(reader io.Reader) (*RuleSet, error) {
	var (
		rawConfig map[string]any
		ruleSet   RuleSet
	)

	dec := yaml.NewDecoder(reader)
	if err := dec.Decode(&rawConfig); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, ErrEmptyRuleSet
		}

		return nil, err
	}

	err := DecodeConfig(rawConfig, &ruleSet)

	return &ruleSet, err
}
