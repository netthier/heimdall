// Copyright 2023 Dimitrij Drus <dadrus@gmx.de>
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

package decision

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.uber.org/fx"

	"github.com/dadrus/heimdall/internal/cache"
	"github.com/dadrus/heimdall/internal/config"
	accesslogmiddleware "github.com/dadrus/heimdall/internal/fiber/middleware/accesslog"
	cachemiddleware "github.com/dadrus/heimdall/internal/fiber/middleware/cache"
	errormiddleware "github.com/dadrus/heimdall/internal/fiber/middleware/errorhandler"
	loggermiddlerware "github.com/dadrus/heimdall/internal/fiber/middleware/logger"
	tracingmiddleware "github.com/dadrus/heimdall/internal/fiber/middleware/opentelemetry"
	prometheusmiddleware "github.com/dadrus/heimdall/internal/fiber/middleware/prometheus"
	"github.com/dadrus/heimdall/internal/x"
)

type appArgs struct {
	fx.In

	Config     *config.Configuration
	Registerer prometheus.Registerer
	Cache      cache.Cache
	Logger     zerolog.Logger
}

func newApp(args appArgs) *fiber.App {
	service := args.Config.Serve.Decision

	app := fiber.New(fiber.Config{
		AppName:                 "Heimdall Decision Service",
		ReadTimeout:             service.Timeout.Read,
		WriteTimeout:            service.Timeout.Write,
		IdleTimeout:             service.Timeout.Idle,
		DisableStartupMessage:   true,
		EnableTrustedProxyCheck: true,
		TrustedProxies: x.IfThenElseExec(service.TrustedProxies != nil,
			func() []string { return *service.TrustedProxies },
			func() []string { return []string{} }),
		JSONDecoder: json.Unmarshal,
		JSONEncoder: json.Marshal,
	})

	app.Use(
		recover.New(recover.Config{EnableStackTrace: true}),
		tracingmiddleware.New(
			tracingmiddleware.WithTracer(otel.GetTracerProvider().Tracer("github.com/dadrus/heimdall/decision")),
		),
	)

	if args.Config.Metrics.Enabled {
		app.Use(prometheusmiddleware.New(
			prometheusmiddleware.WithServiceName("decision"),
			prometheusmiddleware.WithRegisterer(args.Registerer),
		))
	}

	app.Use(
		accesslogmiddleware.New(args.Logger),
		loggermiddlerware.New(args.Logger),
		errormiddleware.New(
			errormiddleware.WithVerboseErrors(service.Respond.Verbose),
			errormiddleware.WithPreconditionErrorCode(service.Respond.With.ArgumentError.Code),
			errormiddleware.WithAuthenticationErrorCode(service.Respond.With.AuthenticationError.Code),
			errormiddleware.WithAuthorizationErrorCode(service.Respond.With.AuthorizationError.Code),
			errormiddleware.WithCommunicationErrorCode(service.Respond.With.CommunicationError.Code),
			errormiddleware.WithMethodErrorCode(service.Respond.With.BadMethodError.Code),
			errormiddleware.WithNoRuleErrorCode(service.Respond.With.NoRuleError.Code),
			errormiddleware.WithInternalServerErrorCode(service.Respond.With.InternalError.Code),
		),
		cachemiddleware.New(args.Cache),
	)

	return app
}
