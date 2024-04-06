//go:build wireinject

//go:generate go run github.com/google/wire/cmd/wire

package main

import (
	"context"
	"github.com/google/wire"
	"github.com/samgozman/go-bloggy/internal/captcha"
	"github.com/samgozman/go-bloggy/internal/config"
	"github.com/samgozman/go-bloggy/internal/db"
	"github.com/samgozman/go-bloggy/internal/github"
	"github.com/samgozman/go-bloggy/internal/jwt"
	"github.com/samgozman/go-bloggy/internal/mailer"
)

func initApp(ctx context.Context, cfg *config.Config) (*tempApp, error) {
	wire.Build(
		db.ProviderSet,
		github.ProviderSet,
		jwt.ProviderSet,
		captcha.ProviderSet,
		mailer.ProviderSet,

		newTempApp,
	)

	return &tempApp{}, nil
}
