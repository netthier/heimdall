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

package unifiers

import (
	"github.com/rs/zerolog"

	"github.com/dadrus/heimdall/internal/heimdall"
	"github.com/dadrus/heimdall/internal/rules/mechanisms/subject"
	"github.com/dadrus/heimdall/internal/rules/mechanisms/template"
	"github.com/dadrus/heimdall/internal/x/errorchain"
)

// by intention. Used only during application bootstrap
// nolint
func init() {
	registerUnifierTypeFactory(
		func(id string, typ string, conf map[string]any) (bool, Unifier, error) {
			if typ != UnifierCookie {
				return false, nil, nil
			}

			unifier, err := newCookieUnifier(id, conf)

			return true, unifier, err
		})
}

type cookieUnifier struct {
	id      string
	cookies map[string]template.Template
}

func newCookieUnifier(id string, rawConfig map[string]any) (*cookieUnifier, error) {
	type Config struct {
		Cookies map[string]template.Template `mapstructure:"cookies"`
	}

	var conf Config
	if err := decodeConfig(rawConfig, &conf); err != nil {
		return nil, errorchain.
			NewWithMessage(heimdall.ErrConfiguration, "failed to unmarshal cookie unifier config").
			CausedBy(err)
	}

	if len(conf.Cookies) == 0 {
		return nil, errorchain.
			NewWithMessage(heimdall.ErrConfiguration, "no cookie definitions provided")
	}

	return &cookieUnifier{
		id:      id,
		cookies: conf.Cookies,
	}, nil
}

func (u *cookieUnifier) Execute(ctx heimdall.Context, sub *subject.Subject) error {
	logger := zerolog.Ctx(ctx.AppContext())
	logger.Debug().Str("_id", u.id).Msg("Unifying using cookie unifier")

	if sub == nil {
		return errorchain.
			NewWithMessage(heimdall.ErrInternal, "failed to execute cookie unifier due to 'nil' subject").
			WithErrorContext(u)
	}

	for name, tmpl := range u.cookies {
		value, err := tmpl.Render(map[string]any{
			"Request": ctx.Request(),
			"Subject": sub,
		})
		if err != nil {
			return errorchain.
				NewWithMessagef(heimdall.ErrInternal, "failed to render value for '%s' cookie", name).
				WithErrorContext(u).
				CausedBy(err)
		}

		ctx.AddCookieForUpstream(name, value)
	}

	return nil
}

func (u *cookieUnifier) WithConfig(config map[string]any) (Unifier, error) {
	if len(config) == 0 {
		return u, nil
	}

	return newCookieUnifier(u.id, config)
}

func (u *cookieUnifier) HandlerID() string { return u.id }

func (u *cookieUnifier) ContinueOnError() bool { return false }
